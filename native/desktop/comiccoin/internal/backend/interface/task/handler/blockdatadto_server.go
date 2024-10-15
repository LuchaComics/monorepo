package handler

import (
	"context"
	"log/slog"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/backend/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/backend/service"
)

type BlockDataDTOServerTaskHandler struct {
	config  *config.Config
	logger  *slog.Logger
	service *service.BlockDataDTOServerService
}

func NewBlockDataDTOServerTaskHandler(
	cfg *config.Config,
	logger *slog.Logger,
	s *service.BlockDataDTOServerService,
) *BlockDataDTOServerTaskHandler {
	return &BlockDataDTOServerTaskHandler{cfg, logger, s}
}

func (h *BlockDataDTOServerTaskHandler) Execute(ctx context.Context) error {
	if serviceExecErr := h.service.Execute(ctx); serviceExecErr != nil {
		return serviceExecErr
	}
	return nil
}
