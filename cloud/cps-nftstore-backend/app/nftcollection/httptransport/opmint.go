package httptransport

import (
	"encoding/json"
	"net/http"
	_ "time/tzdata"

	sub_c "github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/app/nftcollection/controller"
	"github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/utils/httperror"
)

func (h *Handler) OperationMint(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Initialize our array which will tenant all the results from the remote server.
	var requestData sub_c.MintOperationRequestIDO

	defer r.Body.Close()

	// Read the JSON string and convert it into our golang stuct else we need
	// to send a `400 Bad Request` errror message back to the client,
	err := json.NewDecoder(r.Body).Decode(&requestData) // [1]
	if err != nil {
		httperror.ResponseError(w, httperror.NewForSingleField(http.StatusBadRequest, "non_field_error", "payload structure is wrong"))
		return
	}

	err := h.Controller.OperationMint(ctx, &requestData)
	if err != nil {
		httperror.ResponseError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
