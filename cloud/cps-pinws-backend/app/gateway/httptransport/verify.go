package httptransport

import (
	"context"
	"encoding/json"
	"net/http"

	gateway_s "github.com/LuchaComics/monorepo/cloud/cps-pinws-backend/app/gateway/datastore"
	"github.com/LuchaComics/monorepo/cloud/cps-pinws-backend/utils/httperror"
)

type VerifyRequestIDO struct {
	Code string `json:"code"`
}

func UnmarshalVerifyRequest(ctx context.Context, r *http.Request) (*VerifyRequestIDO, error) {
	// Initialize our array which will store all the results from the remote server.
	var requestData VerifyRequestIDO

	defer r.Body.Close()

	// Read the JSON string and convert it into our golang stuct else we need
	// to send a `400 Bad Request` errror message back to the client,
	err := json.NewDecoder(r.Body).Decode(&requestData) // [1]
	if err != nil {
		return nil, httperror.NewForSingleField(http.StatusBadRequest, "non_field_error", "payload structure is wrong")
	}

	// Perform our validation and return validation error on any issues detected.
	if err = ValidateVerifyRequest(&requestData); err != nil {
		return nil, err
	}
	return &requestData, nil
}

func ValidateVerifyRequest(dirtyData *VerifyRequestIDO) error {
	e := make(map[string]string)

	if dirtyData.Code == "" {
		e["code"] = "missing value"
	}
	return nil
}

func (h *Handler) Verify(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	data, err := UnmarshalVerifyRequest(ctx, r)
	if err != nil {
		httperror.ResponseError(w, err)
		return
	}
	res, err := h.Controller.Verify(ctx, data.Code)
	if err != nil {
		httperror.ResponseError(w, err)
		return
	}

	MarshalVerifyResponse(res, w)
}

func MarshalVerifyResponse(responseData *gateway_s.VerifyResponseIDO, w http.ResponseWriter) {
	if err := json.NewEncoder(w).Encode(&responseData); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
