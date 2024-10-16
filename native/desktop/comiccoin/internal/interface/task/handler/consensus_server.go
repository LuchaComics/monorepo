package handler

import (
	"context"
	"log/slog"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/service"
)

type MajorityVoteConsensusServerTaskHandler struct {
	config  *config.Config
	logger  *slog.Logger
	service *service.MajorityVoteConsensusServerService
}

func NewMajorityVoteConsensusServerTaskHandler(
	cfg *config.Config,
	logger *slog.Logger,
	s *service.MajorityVoteConsensusServerService,
) *MajorityVoteConsensusServerTaskHandler {
	return &MajorityVoteConsensusServerTaskHandler{cfg, logger, s}
}

func (h *MajorityVoteConsensusServerTaskHandler) Execute(ctx context.Context) error {
	if serviceExecErr := h.service.Execute(ctx); serviceExecErr != nil {
		return serviceExecErr
	}
	return nil
}
