package httptransport

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	_ "time/tzdata"

	sub_s "github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/app/nft/datastore"
	"github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/utils/httperror"
)

func UnmarshalCreateRequest(ctx context.Context, r *http.Request) (*sub_s.NFT, error) {
	// Initialize our array which will tenant all the results from the remote server.
	var requestData sub_s.NFT

	defer r.Body.Close()

	// Read the JSON string and convert it into our golang stuct else we need
	// to send a `400 Bad Request` errror message back to the client,
	err := json.NewDecoder(r.Body).Decode(&requestData) // [1]
	if err != nil {
		log.Println(err)
		return nil, httperror.NewForSingleField(http.StatusBadRequest, "non_field_error", "payload structure is wrong")
	}

	// Perform our validation and return validation error on any issues detected.
	if err := ValidateCreateRequest(&requestData); err != nil {
		return nil, err
	}
	return &requestData, nil
}

func ValidateCreateRequest(dirtyData *sub_s.NFT) error {
	e := make(map[string]string)

	if dirtyData.TenantID.IsZero() {
		e["tenant_id"] = "missing value"
	}
	// if dirtyData.Name == "" {
	// 	e["name"] = "missing value"
	// }

	if len(e) != 0 {
		return httperror.NewForBadRequest(&e)
	}
	return nil
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	payload, err := UnmarshalCreateRequest(ctx, r)
	if err != nil {
		httperror.ResponseError(w, err)
		return
	}

	res, err := h.Controller.Create(ctx, payload)
	if err != nil {
		httperror.ResponseError(w, err)
		return
	}

	MarshalCreateResponse(res, w)
}

func MarshalCreateResponse(res *sub_s.NFT, w http.ResponseWriter) {
	if err := json.NewEncoder(w).Encode(&res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
