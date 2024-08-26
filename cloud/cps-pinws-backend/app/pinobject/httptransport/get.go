package httptransport

import (
	"encoding/json"
	"net/http"

	"go.mongodb.org/mongo-driver/bson/primitive"

	sub_s "github.com/LuchaComics/monorepo/cloud/cps-pinws-backend/app/pinobject/datastore"
	"github.com/LuchaComics/monorepo/cloud/cps-pinws-backend/utils/httperror"
)

func (h *Handler) GetByRequestID(w http.ResponseWriter, r *http.Request, requestID string) {
	ctx := r.Context()

	objectRequestID, err := primitive.ObjectIDFromHex(requestID)
	if err != nil {
		httperror.ResponseError(w, err)
		return
	}

	m, err := h.Controller.GetByRequestID(ctx, objectRequestID)
	if err != nil {
		httperror.ResponseError(w, err)
		return
	}

	MarshalDetailResponse(m, w)
}

func MarshalDetailResponse(res *sub_s.PinObject, w http.ResponseWriter) {
	if err := json.NewEncoder(w).Encode(&res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
