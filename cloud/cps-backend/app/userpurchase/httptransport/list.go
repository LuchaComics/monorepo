package httptransport

import (
	"encoding/json"
	"net/http"
	"strconv"

	sub_s "github.com/LuchaComics/monorepo/cloud/cps-backend/app/userpurchase/datastore"
	"github.com/LuchaComics/monorepo/cloud/cps-backend/utils/httperror"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	f := &sub_s.UserPurchasePaginationListFilter{
		Cursor:          "",
		PageSize:        25,
		SortField:       "created_at",
		SortOrder:       -1, // 1=ascending | -1=descending
		ExcludeArchived: true,
	}

	// Here is where you extract url parameters.
	query := r.URL.Query()

	cursor := query.Get("cursor")
	if cursor != "" {
		f.Cursor = cursor
	}

	pageSize := query.Get("page_size")
	if pageSize != "" {
		pageSize, _ := strconv.ParseInt(pageSize, 10, 64)
		if pageSize == 0 || pageSize > 250 {
			pageSize = 250
		}
		f.PageSize = pageSize
	}

	storeID := query.Get("store_id")
	if storeID != "" {
		storeID, err := primitive.ObjectIDFromHex(storeID)
		if err != nil {
			httperror.ResponseError(w, err)
			return
		}
		f.StoreID = storeID
	} else {
		err := httperror.NewForSingleField(http.StatusBadRequest, "store_id", "missing value")
		httperror.ResponseError(w, err)
		return
	}

	searchKeyword := query.Get("search")
	if searchKeyword != "" {
		f.SearchText = searchKeyword
	}

	m, err := h.Controller.ListByFilter(ctx, f)
	if err != nil {
		httperror.ResponseError(w, err)
		return
	}

	MarshalListResponse(m, w)
}

func MarshalListResponse(res *sub_s.UserPurchasePaginationListResult, w http.ResponseWriter) {
	if err := json.NewEncoder(w).Encode(&res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
