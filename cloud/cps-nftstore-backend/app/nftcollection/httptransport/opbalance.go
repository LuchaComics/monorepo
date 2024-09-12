package httptransport

import (
	"encoding/json"
	"net/http"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/utils/httperror"
)

func (h *Handler) OperationGetWalletBalanceByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Here is where you extract url parameters.
	query := r.URL.Query()

	collectionIDHex := query.Get("collection_id")
	if collectionIDHex == "" {
		httperror.ResponseError(w, httperror.NewForBadRequestWithSingleField("collection_id", "missing url parameter"))
		return
	}

	collectionID, err := primitive.ObjectIDFromHex(collectionIDHex)
	if err != nil {
		httperror.ResponseError(w, err)
		return
	}

	res, err := h.Controller.OperationGetWalletBalanceByID(ctx, collectionID)
	if err != nil {
		httperror.ResponseError(w, err)
		return
	}

	if err := json.NewEncoder(w).Encode(&res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
