package httptransport

import (
	"encoding/json"
	"net/http"
	"strconv"

	sub_s "github.com/LuchaComics/monorepo/cloud/cps-backend/app/comicsub/datastore"
	"github.com/LuchaComics/monorepo/cloud/cps-backend/utils/httperror"
	"github.com/bartmika/timekit"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Initialize the list filter with base results and then override them with the URL parameters.
	f := &sub_s.ComicSubmissionPaginationListFilter{
		Cursor:          "",
		PageSize:        25,
		SortField:       "created_at",
		SortOrder:       -1, // 1=ascending | -1=descending
		ExcludeArchived: true,
	}

	// Here is where you extract url parameters.
	query := r.URL.Query()
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

	searchText := query.Get("search")
	if searchText != "" {
		f.SearchText = searchText
	}

	statusStr := query.Get("status")
	if statusStr != "" {
		status, _ := strconv.ParseInt(statusStr, 10, 64)
		f.Status = int8(status)
	}
	createdAtGTEStr := query.Get("created_at_gte")
	if createdAtGTEStr != "" {
		createdAtGTE, err := timekit.ParseJavaScriptTimeString(createdAtGTEStr)
		if err != nil {
			httperror.ResponseError(w, err)
			return
		}
		f.CreatedAtGTE = createdAtGTE
	}

	// Fet
	m, err := h.Controller.ListByFilter(ctx, f)
	if err != nil {
		httperror.ResponseError(w, err)
		return
	}

	MarshalListResponse(m, w)
}

func MarshalListResponse(res *sub_s.ComicSubmissionPaginationListResult, w http.ResponseWriter) {
	if err := json.NewEncoder(w).Encode(&res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
