package handler

import (
	"context"
	"log/slog"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/service"
)

type SyncBlockDataDTOTaskHandler struct {
	config  *config.Config
	logger  *slog.Logger
	service *service.SyncBlockDataDTOService
}

func NewSyncBlockDataDTOTaskHandler(
	cfg *config.Config,
	logger *slog.Logger,
	s *service.SyncBlockDataDTOService,
) *SyncBlockDataDTOTaskHandler {
	return &SyncBlockDataDTOTaskHandler{cfg, logger, s}
}

func (h *SyncBlockDataDTOTaskHandler) Execute(ctx context.Context) error {
	if serviceExecErr := h.service.Execute(ctx); serviceExecErr != nil {
		return serviceExecErr
	}
	return nil
}
