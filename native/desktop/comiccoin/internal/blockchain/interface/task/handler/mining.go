package handler

import (
	"context"
	"log/slog"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/service"
)

type MiningTaskHandler struct {
	config  *config.Config
	logger  *slog.Logger
	service *service.MiningService
}

func NewMiningTaskHandler(
	cfg *config.Config,
	logger *slog.Logger,
	s *service.MiningService,
) *MiningTaskHandler {
	return &MiningTaskHandler{cfg, logger, s}
}

func (h *MiningTaskHandler) Execute(ctx context.Context) error {
	if serviceExecErr := h.service.Execute(ctx); serviceExecErr != nil {
		return serviceExecErr
	}
	return nil
}
