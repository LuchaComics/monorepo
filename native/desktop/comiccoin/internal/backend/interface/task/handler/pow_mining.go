package handler

import (
	"context"
	"log/slog"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/backend/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/backend/service"
)

type ProofOfWorkMiningTaskHandler struct {
	config  *config.Config
	logger  *slog.Logger
	service *service.ProofOfWorkMiningService
}

func NewProofOfWorkMiningTaskHandler(
	cfg *config.Config,
	logger *slog.Logger,
	s *service.ProofOfWorkMiningService,
) *ProofOfWorkMiningTaskHandler {
	return &ProofOfWorkMiningTaskHandler{cfg, logger, s}
}

func (h *ProofOfWorkMiningTaskHandler) Execute(ctx context.Context) error {
	if serviceExecErr := h.service.Execute(ctx); serviceExecErr != nil {
		return serviceExecErr
	}
	return nil
}
