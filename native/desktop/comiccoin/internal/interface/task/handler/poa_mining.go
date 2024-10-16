package handler

import (
	"context"
	"log/slog"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/service"
)

type ProofOfAuthorityMiningTaskHandler struct {
	config  *config.Config
	logger  *slog.Logger
	service *service.ProofOfAuthorityMiningService
}

func NewProofOfAuthorityMiningTaskHandler(
	cfg *config.Config,
	logger *slog.Logger,
	s *service.ProofOfAuthorityMiningService,
) *ProofOfAuthorityMiningTaskHandler {
	return &ProofOfAuthorityMiningTaskHandler{cfg, logger, s}
}

func (h *ProofOfAuthorityMiningTaskHandler) Execute(ctx context.Context) error {
	if serviceExecErr := h.service.Execute(ctx); serviceExecErr != nil {
		return serviceExecErr
	}
	return nil
}
