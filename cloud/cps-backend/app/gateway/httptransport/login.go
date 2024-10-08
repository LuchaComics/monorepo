package httptransport

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"strings"

	gateway_s "github.com/LuchaComics/monorepo/cloud/cps-backend/app/gateway/datastore"
	"github.com/LuchaComics/monorepo/cloud/cps-backend/utils/httperror"
)

type LoginRequestIDO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *Handler) unmarshalLoginRequest(ctx context.Context, r *http.Request) (*LoginRequestIDO, error) {
	// Initialize our array which will store all the results from the remote server.
	var requestData LoginRequestIDO

	defer r.Body.Close()

	var rawJSON bytes.Buffer
	teeReader := io.TeeReader(r.Body, &rawJSON) // TeeReader allows you to read the JSON and capture it

	// Read the JSON string and convert it into our golang stuct else we need
	// to send a `400 Bad Request` errror message back to the client,
	err := json.NewDecoder(teeReader).Decode(&requestData) // [1]
	if err != nil {
		h.Logger.Error("decoding error",
			slog.Any("err", err),
			slog.String("json", rawJSON.String()),
		)
		return nil, httperror.NewForSingleField(http.StatusBadRequest, "non_field_error", "payload structure is wrong")
	}

	// Defensive Code: For security purposes we need to remove all whitespaces from the email and lower the characters.
	requestData.Email = strings.ToLower(requestData.Email)
	requestData.Email = strings.ReplaceAll(requestData.Email, " ", "")

	// Perform our validation and return validation error on any issues detected.
	if err = ValidateLoginRequest(&requestData); err != nil {
		return nil, err
	}
	return &requestData, nil
}

func ValidateLoginRequest(dirtyData *LoginRequestIDO) error {
	e := make(map[string]string)

	if dirtyData.Email == "" {
		e["email"] = "missing value"
	}
	if len(dirtyData.Email) > 255 {
		e["email"] = "too long"
	}
	if dirtyData.Password == "" {
		e["password"] = "missing value"
	}

	if len(e) != 0 {
		return httperror.NewForBadRequest(&e)
	}
	return nil
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	data, err := h.unmarshalLoginRequest(ctx, r)
	if err != nil {
		httperror.ResponseError(w, err)
		return
	}

	res, err := h.Controller.Login(ctx, data.Email, data.Password)
	if err != nil {
		httperror.ResponseError(w, err)
		return
	}
	MarshalLoginResponse(res, w)
}

func MarshalLoginResponse(responseData *gateway_s.LoginResponseIDO, w http.ResponseWriter) {
	if err := json.NewEncoder(w).Encode(&responseData); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
