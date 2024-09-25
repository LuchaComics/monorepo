package httptransport

import (
	"context"
	"encoding/json"
	"net/http"

	pt_c "github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/app/blockchain/controller"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/utils/httperror"
)

func unmarshalTransferRequest(ctx context.Context, r *http.Request) (*pt_c.BlockchainTransferRequestIDO, error) {
	// Initialize our array which will store all the results from the remote server.
	var requestData *pt_c.BlockchainTransferRequestIDO

	defer r.Body.Close()

	// Read the JSON string and convert it into our golang stuct else we need
	// to send a `400 Bad Request` errror message back to the client,
	err := json.NewDecoder(r.Body).Decode(&requestData) // [1]
	if err != nil {
		return nil, httperror.NewForSingleField(http.StatusBadRequest, "non_field_error", "payload structure is wrong")
	}

	return requestData, nil
}

func (h *Handler) Transfer(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	data, err := unmarshalTransferRequest(ctx, r)
	if err != nil {
		httperror.ResponseError(w, err)
		return
	}

	res, err := h.Controller.Transfer(ctx, data)
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
