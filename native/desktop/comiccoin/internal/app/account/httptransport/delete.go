package httptransport

import (
	"net/http"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/utils/httperror"
)

func (h *Handler) DeleteByName(w http.ResponseWriter, r *http.Request, name string) {
	ctx := r.Context()
	if err := h.Controller.DeleteByName(ctx, name); err != nil {
		httperror.ResponseError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
