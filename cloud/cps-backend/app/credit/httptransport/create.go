package httptransport

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	credit_c "github.com/LuchaComics/monorepo/cloud/cps-backend/app/credit/controller"
	credit_d "github.com/LuchaComics/monorepo/cloud/cps-backend/app/credit/datastore"
	"github.com/LuchaComics/monorepo/cloud/cps-backend/utils/httperror"
)

func UnmarshalCreateRequest(ctx context.Context, r *http.Request) (*credit_c.CreditCreateRequest, error) {
	// Initialize our array which will store all the results from the remote server.
	var requestData credit_c.CreditCreateRequest

	defer r.Body.Close()

	// Read the JSON string and convert it into our golang stuct else we need
	// to send a `400 Bad Request` errror message back to the client,
	err := json.NewDecoder(r.Body).Decode(&requestData) // [1]
	if err != nil {
		log.Println(err)
		return nil, httperror.NewForSingleField(http.StatusBadRequest, "non_field_error", "payload structure is wrong")
	}

	// Perform our validation and return validation error on any issues detected.
	if err := ValidateCreateRequest(&requestData); err != nil {
		return nil, err
	}
	return &requestData, nil
}

func ValidateCreateRequest(dirtyData *credit_c.CreditCreateRequest) error {
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

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	data, err := UnmarshalCreateRequest(ctx, r)
	if err != nil {
		httperror.ResponseError(w, err)
		return
	}
	if err := h.Controller.Create(ctx, data); err != nil {
		httperror.ResponseError(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
