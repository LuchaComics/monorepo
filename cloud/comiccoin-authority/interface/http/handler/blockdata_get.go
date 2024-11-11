package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/common/httperror"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/service"
)

type GetBlockDataHTTPHandler struct {
	logger  *slog.Logger
	service *service.GetBlockDataService
}

func NewGetBlockDataHTTPHandler(
	logger *slog.Logger,
	s1 *service.GetBlockDataService,
) *GetBlockDataHTTPHandler {
	return &GetBlockDataHTTPHandler{logger, s1}
}

func (h *GetBlockDataHTTPHandler) Execute(w http.ResponseWriter, r *http.Request, hash string) {
	ctx := r.Context()
	h.logger.Debug("BlockData requested")

	resp, err := h.service.Execute(ctx, hash)
	if err != nil {
		httperror.ResponseError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(&resp); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
