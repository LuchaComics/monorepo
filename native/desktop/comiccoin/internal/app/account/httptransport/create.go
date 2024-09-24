package httptransport

import (
	"context"
	"encoding/json"
	"net/http"

	acc_c "github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/app/account/controller"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/utils/httperror"
)

func unmarshalCreateRequest(ctx context.Context, r *http.Request) (*acc_c.AccountCreateRequestIDO, error) {
	// Initialize our array which will store all the results from the remote server.
	var requestData *acc_c.AccountCreateRequestIDO

	defer r.Body.Close()

	// Read the JSON string and convert it into our golang stuct else we need
	// to send a `400 Bad Request` errror message back to the client,
	err := json.NewDecoder(r.Body).Decode(&requestData) // [1]
	if err != nil {
		return nil, httperror.NewForSingleField(http.StatusBadRequest, "non_field_error", "payload structure is wrong")
	}

	return requestData, nil
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	data, err := unmarshalCreateRequest(ctx, r)
	if err != nil {
		httperror.ResponseError(w, err)
		return
	}

	res, err := h.Controller.Create(ctx, data)
	if err != nil {
		httperror.ResponseError(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(&res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
