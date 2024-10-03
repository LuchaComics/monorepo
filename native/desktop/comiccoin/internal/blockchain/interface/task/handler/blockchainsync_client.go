package handler

import (
	"context"
	"log/slog"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/service"
)

type BlockchainSyncClientTaskHandler struct {
	config  *config.Config
	logger  *slog.Logger
	service *service.BlockchainSyncClientService
}

func NewBlockchainSyncClientTaskHandler(
	cfg *config.Config,
	logger *slog.Logger,
	s *service.BlockchainSyncClientService,
) *BlockchainSyncClientTaskHandler {
	return &BlockchainSyncClientTaskHandler{cfg, logger, s}
}

func (h *BlockchainSyncClientTaskHandler) Execute(ctx context.Context) error {
	if serviceExecErr := h.service.Execute(ctx); serviceExecErr != nil {
		return serviceExecErr
	}
	return nil
}
