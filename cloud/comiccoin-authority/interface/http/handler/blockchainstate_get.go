package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/common/httperror"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/service"
)

type GetBlockchainStateHTTPHandler struct {
	logger  *slog.Logger
	service *service.GetBlockchainStateService
}

func NewGetBlockchainStateHTTPHandler(
	logger *slog.Logger,
	s1 *service.GetBlockchainStateService,
) *GetBlockchainStateHTTPHandler {
	return &GetBlockchainStateHTTPHandler{logger, s1}
}

func (h *GetBlockchainStateHTTPHandler) Execute(w http.ResponseWriter, r *http.Request, chainIDstr string) {
	ctx := r.Context()
	h.logger.Debug("Blockchain state requested")

	blockchainState, err := h.service.Execute(ctx, 1)
	if err != nil {
		httperror.ResponseError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(&blockchainState); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
