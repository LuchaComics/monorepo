package handler

import (
	"context"
	"log/slog"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/service"
)

type ConsensusTaskHandler struct {
	config  *config.Config
	logger  *slog.Logger
	service *service.ConsensusService
}

func NewConsensusTaskHandler(
	cfg *config.Config,
	logger *slog.Logger,
	s *service.ConsensusService,
) *ConsensusTaskHandler {
	return &ConsensusTaskHandler{cfg, logger, s}
}

func (h *ConsensusTaskHandler) Execute(ctx context.Context) error {
	if serviceExecErr := h.service.Execute(ctx); serviceExecErr != nil {
		return serviceExecErr
	}
	return nil
}
