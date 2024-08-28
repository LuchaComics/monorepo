package httptransport

import (
	"encoding/json"
	"net/http"
	"strconv"

	sub_s "github.com/LuchaComics/monorepo/cloud/cps-ipfsstore-backend/app/user/datastore"
	"github.com/LuchaComics/monorepo/cloud/cps-ipfsstore-backend/utils/httperror"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (h *Handler) ListAsSelectOptions(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	f := &sub_s.UserPaginationListFilter{
		Cursor:    "",
		PageSize:  1_000_000,
		SortField: "name",
		SortOrder: 1, // 1=ascending | -1=descending
		Status:    0, // All
	}

	// Here is where you extract url parameters.
	query := r.URL.Query()

	// Apply search text if it exists in url parameter.
	searchKeyword := query.Get("search")
	if searchKeyword != "" {
		f.SearchText = searchKeyword
	}

	// Apply filters it exists in url parameter.
	firstName := query.Get("first_name")
	if firstName != "" {
		f.FirstName = firstName
	}
	lastName := query.Get("first_name")
	if lastName != "" {
		f.LastName = lastName
	}
	email := query.Get("email")
	if email != "" {
		f.Email = email
	}
	phone := query.Get("phone")
	if phone != "" {
		f.Phone = phone
	}
	statusStr := query.Get("status")
	if statusStr != "" {
		status, _ := strconv.ParseInt(statusStr, 10, 64)
		f.Status = int8(status)
	}
	storeID := query.Get("tenant_id")
	if storeID != "" {
		storeID, err := primitive.ObjectIDFromHex(storeID)
		if err != nil {
			httperror.ResponseError(w, err)
			return
		}
		f.TenantID = storeID
	}

	// Perform our database operation.
	m, err := h.Controller.ListAsSelectOptionByFilter(ctx, f)
	if err != nil {
		httperror.ResponseError(w, err)
		return
	}

	MarshalListAsSelectOptionResponse(m, w)
}

func MarshalListAsSelectOptionResponse(res []*sub_s.UserAsSelectOption, w http.ResponseWriter) {
	if err := json.NewEncoder(w).Encode(&res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
