package handler

import (
	"context"
	"log/slog"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/service"
)

type BlockchainSyncServerTaskHandler struct {
	config  *config.Config
	logger  *slog.Logger
	service *service.BlockchainSyncServerService
}

func NewBlockchainSyncServerTaskHandler(
	cfg *config.Config,
	logger *slog.Logger,
	s *service.BlockchainSyncServerService,
) *BlockchainSyncServerTaskHandler {
	return &BlockchainSyncServerTaskHandler{cfg, logger, s}
}

func (h *BlockchainSyncServerTaskHandler) Execute(ctx context.Context) error {
	if serviceExecErr := h.service.Execute(ctx); serviceExecErr != nil {
		return serviceExecErr
	}
	return nil
}
