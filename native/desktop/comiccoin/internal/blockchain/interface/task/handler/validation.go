package handler

import (
	"context"
	"log/slog"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/service"
)

type ValidationTaskHandler struct {
	config  *config.Config
	logger  *slog.Logger
	service *service.ValidationService
}

func NewValidationTaskHandler(
	cfg *config.Config,
	logger *slog.Logger,
	s *service.ValidationService,
) *ValidationTaskHandler {
	return &ValidationTaskHandler{cfg, logger, s}
}

func (h *ValidationTaskHandler) Execute(ctx context.Context) error {
	h.logger.Info("Validation is running...")
	if serviceExecErr := h.service.Execute(ctx); serviceExecErr != nil {
		return serviceExecErr
	}
	return nil
}
