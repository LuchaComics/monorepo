package handler

import (
	"context"
	"log/slog"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/backend/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/backend/service"
)

type ProofOfAuthorityValidationTaskHandler struct {
	config  *config.Config
	logger  *slog.Logger
	service *service.ProofOfAuthorityValidationService
}

func NewProofOfAuthorityValidationTaskHandler(
	cfg *config.Config,
	logger *slog.Logger,
	s *service.ProofOfAuthorityValidationService,
) *ProofOfAuthorityValidationTaskHandler {
	return &ProofOfAuthorityValidationTaskHandler{cfg, logger, s}
}

func (h *ProofOfAuthorityValidationTaskHandler) Execute(ctx context.Context) error {
	if serviceExecErr := h.service.Execute(ctx); serviceExecErr != nil {
		return serviceExecErr
	}
	return nil
}
