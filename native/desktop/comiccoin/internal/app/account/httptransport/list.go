package httptransport

import (
	"encoding/json"
	"net/http"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/utils/httperror"
)

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	res, err := h.Controller.List(ctx)
	if err != nil {
		httperror.ResponseError(w, err)
		return
	}

	if err := json.NewEncoder(w).Encode(&res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
