package stripe

import (
	"encoding/json"
	"net/http"

	"github.com/LuchaComics/monorepo/cloud/cps-backend/utils/httperror"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CreateStripeCheckoutSessionURLForComicSubmissionIDResponseIDO struct {
	CheckoutSessionURL string
}

func (h *Handler) CreateStripeCheckoutSessionURLForComicSubmissionID(w http.ResponseWriter, r *http.Request, comicSubmissionIDString string) {
	ctx := r.Context()

	comicSubmissionID, err := primitive.ObjectIDFromHex(comicSubmissionIDString)
	if err != nil {
		httperror.ResponseError(w, err)
		return
	}

	checkoutSessionURL, err := h.Controller.CreateStripeCheckoutSessionURLForComicSubmissionID(ctx, comicSubmissionID)
	if err != nil {
		httperror.ResponseError(w, err)
		return
	}

	// Create temporary structure
	res := &CreateStripeCheckoutSessionURLForComicSubmissionIDResponseIDO{
		CheckoutSessionURL: checkoutSessionURL,
	}

	if err := json.NewEncoder(w).Encode(&res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
