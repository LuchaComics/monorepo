package httptransport

import (
	"encoding/json"
	"net/http"

	sub_s "github.com/LuchaComics/monorepo/cloud/cps-backend/app/offer/datastore"
	"github.com/LuchaComics/monorepo/cloud/cps-backend/utils/httperror"
)

func (h *Handler) ListAsSelectOptions(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	f := &sub_s.OfferPaginationListFilter{
		Cursor:    "",
		PageSize:  1_000_000,
		SortField: "name",
		SortOrder: 1, //1=Asc
		Status:    sub_s.StatusActive,
	}

	// Perform our database operation.
	m, err := h.Controller.ListAsSelectOptionByFilter(ctx, f)
	if err != nil {
		httperror.ResponseError(w, err)
		return
	}

	MarshalListAsSelectOptionResponse(m, w)
}

func MarshalListAsSelectOptionResponse(res []*sub_s.OfferAsSelectOption, w http.ResponseWriter) {
	if err := json.NewEncoder(w).Encode(&res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
