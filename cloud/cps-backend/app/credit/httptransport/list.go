package httptransport

import (
	"encoding/json"
	"net/http"
	"strconv"

	sub_s "github.com/LuchaComics/monorepo/cloud/cps-backend/app/credit/datastore"
	"github.com/LuchaComics/monorepo/cloud/cps-backend/utils/httperror"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	f := &sub_s.CreditPaginationListFilter{
		Cursor:    "",
		PageSize:  25,
		SortField: "created_at",
		SortOrder: 1, // 1=ascending | -1=descending
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
	}

	userID := query.Get("user_id")
	if userID != "" {
		userID, err := primitive.ObjectIDFromHex(userID)
		if err != nil {
			httperror.ResponseError(w, err)
			return
		}
		f.UserID = userID
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

func MarshalListResponse(res *sub_s.CreditPaginationListResult, w http.ResponseWriter) {
	if err := json.NewEncoder(w).Encode(&res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
