package httptransport

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"strings"
	_ "time/tzdata"

	gateway_s "github.com/LuchaComics/monorepo/cloud/cps-backend/app/gateway/datastore"
	"github.com/LuchaComics/monorepo/cloud/cps-backend/utils/httperror"
)

func (h *Handler) unmarshalRegisterCustomerRequest(ctx context.Context, r *http.Request) (*gateway_s.RegisterCustomerRequestIDO, error) {
	// Initialize our array which will store all the results from the remote server.
	var requestData gateway_s.RegisterCustomerRequestIDO

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

	return &requestData, nil
}

func (h *Handler) RegisterCustomer(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	data, err := h.unmarshalRegisterCustomerRequest(ctx, r)
	if err != nil {
		httperror.ResponseError(w, err)
		return
	}
	if err := h.Controller.RegisterCustomer(ctx, data); err != nil {
		httperror.ResponseError(w, err)
		return
	}
	w.WriteHeader(http.StatusCreated)
}
