package httptransport

import (
	"encoding/json"
	"net/http"
	"strconv"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/utils/httperror"
)

func (h *Handler) OperationGetTokenURI(w http.ResponseWriter, r *http.Request) {
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

	tokenIDStr := query.Get("token_id")
	if tokenIDStr == "" {
		httperror.ResponseError(w, httperror.NewForBadRequestWithSingleField("token_id", "missing url parameter"))
		return
	}

	tokenID, _ := strconv.ParseUint(tokenIDStr, 10, 64)

	res, err := h.Controller.OperationGetTokenURI(ctx, collectionID, tokenID)
	if err != nil {
		httperror.ResponseError(w, err)
		return
	}

	if err := json.NewEncoder(w).Encode(&res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
