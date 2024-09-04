package httptransport

import (
	"encoding/json"
	"net/http"
	"strconv"

	pinobject_s "github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/app/pinobject/datastore"
	"github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/utils/httperror"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	f := &pinobject_s.PinObjectPaginationListFilter{
		Cursor:    "",
		ProjectID: primitive.NilObjectID,
		PageSize:  25,
		SortField: "created",
		SortOrder: 1, // 1=ascending | -1=descending
	}

	// Here is where you extract url parameters.
	query := r.URL.Query()

	cursor := query.Get("cursor")
	if cursor != "" {
		f.Cursor = cursor
	}

	projectID := query.Get("project_id")
	if projectID != "" {
		projectID, err := primitive.ObjectIDFromHex(projectID)
		if err != nil {
			httperror.ResponseError(w, err)
			return
		}
		f.ProjectID = projectID
	}

	pageSize := query.Get("page_size")
	if pageSize != "" {
		pageSize, _ := strconv.ParseInt(pageSize, 10, 64)
		if pageSize == 0 || pageSize > 250 {
			pageSize = 250
		}
		f.PageSize = pageSize
	}

	m, err := h.Controller.ListByFilter(ctx, f)
	if err != nil {
		httperror.ResponseError(w, err)
		return
	}

	MarshalListResponse(m, w)
}

func MarshalListResponse(res *pinobject_s.PinObjectPaginationListResult, w http.ResponseWriter) {
	if err := json.NewEncoder(w).Encode(&res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) ListAsSelectOptionByFilter(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	f := &pinobject_s.PinObjectPaginationListFilter{
		Cursor:    "",
		PageSize:  1_000_000,
		SortField: "created",
		SortOrder: pinobject_s.SortOrderDescending,
	}

	m, err := h.Controller.ListAsSelectOptionByFilter(ctx, f)
	if err != nil {
		httperror.ResponseError(w, err)
		return
	}

	MarshalListAsSelectOptionResponse(m, w)
}

func MarshalListAsSelectOptionResponse(res []*pinobject_s.PinObjectAsSelectOption, w http.ResponseWriter) {
	if err := json.NewEncoder(w).Encode(&res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
