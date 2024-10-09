package handler

import (
	"context"
	"log/slog"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/service"
)

type ProofOfWorkValidationTaskHandler struct {
	config  *config.Config
	logger  *slog.Logger
	service *service.ProofOfWorkValidationService
}

func NewProofOfWorkValidationTaskHandler(
	cfg *config.Config,
	logger *slog.Logger,
	s *service.ProofOfWorkValidationService,
) *ProofOfWorkValidationTaskHandler {
	return &ProofOfWorkValidationTaskHandler{cfg, logger, s}
}

func (h *ProofOfWorkValidationTaskHandler) Execute(ctx context.Context) error {
	if serviceExecErr := h.service.Execute(ctx); serviceExecErr != nil {
		return serviceExecErr
	}
	return nil
}
