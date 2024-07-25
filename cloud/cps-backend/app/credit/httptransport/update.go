package httptransport

import (
	"context"
	"encoding/json"
	"net/http"

	credit_d "github.com/LuchaComics/monorepo/cloud/cps-backend/app/credit/datastore"
	sub_s "github.com/LuchaComics/monorepo/cloud/cps-backend/app/credit/datastore"
	"github.com/LuchaComics/monorepo/cloud/cps-backend/utils/httperror"
)

func UnmarshalUpdateRequest(ctx context.Context, r *http.Request) (*sub_s.Credit, error) {
	// Initialize our array which will store all the results from the remote server.
	var requestData sub_s.Credit

	defer r.Body.Close()

	// Read the JSON string and convert it into our golang stuct else we need
	// to send a `400 Bad Request` errror message back to the client,
	err := json.NewDecoder(r.Body).Decode(&requestData) // [1]
	if err != nil {
		return nil, httperror.NewForSingleField(http.StatusBadRequest, "non_field_error", "payload structure is wrong")
	}

	// Perform our validation and return validation error on any issues detected.
	if err := ValidateUpdateRequest(&requestData); err != nil {
		return nil, err
	}

	return &requestData, nil
}

func ValidateUpdateRequest(dirtyData *sub_s.Credit) error {
	e := make(map[string]string)

	if dirtyData.BusinessFunction == 0 {
		e["business_function"] = "missing value"
	}
	if dirtyData.BusinessFunction == credit_d.BusinessFunctionGrantFreeSubmission {
		if dirtyData.OfferID.IsZero() {
			e["offer_id"] = "missing value"
		}
	}
	if dirtyData.UserID.IsZero() {
		e["user_id"] = "missing value"
	}

	if len(e) != 0 {
		return httperror.NewForBadRequest(&e)
	}
	return nil
}

func (h *Handler) UpdateByID(w http.ResponseWriter, r *http.Request, id string) {
	ctx := r.Context()

	data, err := UnmarshalUpdateRequest(ctx, r)
	if err != nil {
		httperror.ResponseError(w, err)
		return
	}

	org, err := h.Controller.UpdateByID(ctx, data)
	if err != nil {
		httperror.ResponseError(w, err)
		return
	}

	MarshalUpdateResponse(org, w)
}

func MarshalUpdateResponse(res *sub_s.Credit, w http.ResponseWriter) {
	if err := json.NewEncoder(w).Encode(&res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
