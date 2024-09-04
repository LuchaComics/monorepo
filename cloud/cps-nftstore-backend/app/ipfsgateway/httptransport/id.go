package httptransport

import (
	"encoding/json"
	"net/http"

	"github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/utils/httperror"
)

// GetIpfsNodeInfo function returns the current running IPFS node server
// information. This is not according to IPFS spec but it's useful for
// our purposes.
func (h *Handler) GetIpfsNodeInfo(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	info, err := h.Controller.GetIpfsNodeInfo(ctx)
	if err != nil {
		httperror.ResponseError(w, err)
		return
	}
	if info == nil {
		httperror.ResponseError(w, httperror.NewForNotFoundWithSingleField("error", "does not exist"))
		return
	}

	if err := json.NewEncoder(w).Encode(&info); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
