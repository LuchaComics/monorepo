package httptransport

import (
	"context"
	"encoding/json"
	"net/http"

	gateway_c "github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/app/gateway/controller"
	"github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/utils/httperror"
)

func UnmarshalProfileChangePasswordRequest(ctx context.Context, r *http.Request) (*gateway_c.ProfileChangePasswordRequestIDO, error) {
	// Initialize our array which will store all the results from the remote server.
	var requestData gateway_c.ProfileChangePasswordRequestIDO

	defer r.Body.Close()

	// Read the JSON string and convert it into our golang stuct else we need
	// to send a `400 Bad Request` errror message back to the client,
	err := json.NewDecoder(r.Body).Decode(&requestData) // [1]
	if err != nil {
		return nil, httperror.NewForSingleField(http.StatusBadRequest, "non_field_error", "payload structure is wrong")
	}

	// Return our result
	return &requestData, nil
}

func (h *Handler) ProfileChangePassword(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	data, err := UnmarshalProfileChangePasswordRequest(ctx, r)
	if err != nil {
		httperror.ResponseError(w, err)
		return
	}

	if err := h.Controller.ProfileChangePassword(ctx, data); err != nil {
		httperror.ResponseError(w, err)
		return
	}

	// Get the request
	h.Profile(w, r)
}
