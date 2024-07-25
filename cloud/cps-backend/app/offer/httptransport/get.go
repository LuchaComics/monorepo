package httptransport

import (
	"encoding/json"
	"net/http"
	"strconv"

	"go.mongodb.org/mongo-driver/bson/primitive"

	sub_s "github.com/LuchaComics/monorepo/cloud/cps-backend/app/offer/datastore"
	"github.com/LuchaComics/monorepo/cloud/cps-backend/utils/httperror"
)

func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request, id string) {
	ctx := r.Context()

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		httperror.ResponseError(w, err)
		return
	}

	m, err := h.Controller.GetByID(ctx, objectID)
	if err != nil {
		httperror.ResponseError(w, err)
		return
	}

	MarshalDetailResponse(m, w)
}

func MarshalDetailResponse(res *sub_s.Offer, w http.ResponseWriter) {
	if err := json.NewEncoder(w).Encode(&res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) GetByServiceType(w http.ResponseWriter, r *http.Request, serviceTypeStr string) {
	ctx := r.Context()
	serviceType, _ := strconv.ParseInt(serviceTypeStr, 10, 64)
	m, err := h.Controller.GetByServiceType(ctx, int8(serviceType))
	if err != nil {
		httperror.ResponseError(w, err)
		return
	}
	MarshalDetailResponse(m, w)
}
