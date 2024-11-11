package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/common/httperror"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/service"
)

type GetGenesisBlockDataHTTPHandler struct {
	logger  *slog.Logger
	service *service.GetGenesisBlockDataService
}

func NewGetGenesisBlockDataHTTPHandler(
	logger *slog.Logger,
	s1 *service.GetGenesisBlockDataService,
) *GetGenesisBlockDataHTTPHandler {
	return &GetGenesisBlockDataHTTPHandler{logger, s1}
}

type GenesisBlockDataResponseIDO struct {
	GenesisBlockData string `json:"GenesisBlockData"`
}

func (h *GetGenesisBlockDataHTTPHandler) Execute(w http.ResponseWriter, r *http.Request, chainIDstr string) {
	ctx := r.Context()
	h.logger.Debug("GenesisBlockData requested")

	genesis, err := h.service.Execute(ctx, 1)
	if err != nil {
		httperror.ResponseError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(&genesis); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
