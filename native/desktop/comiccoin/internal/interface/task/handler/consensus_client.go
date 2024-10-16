package handler

import (
	"context"
	"log/slog"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/service"
)

type MajorityVoteConsensusClientTaskHandler struct {
	config  *config.Config
	logger  *slog.Logger
	service *service.MajorityVoteConsensusClientService
}

func NewMajorityVoteConsensusClientTaskHandler(
	cfg *config.Config,
	logger *slog.Logger,
	s *service.MajorityVoteConsensusClientService,
) *MajorityVoteConsensusClientTaskHandler {
	return &MajorityVoteConsensusClientTaskHandler{cfg, logger, s}
}

func (h *MajorityVoteConsensusClientTaskHandler) Execute(ctx context.Context) error {
	if serviceExecErr := h.service.Execute(ctx); serviceExecErr != nil {
		return serviceExecErr
	}
	return nil
}
