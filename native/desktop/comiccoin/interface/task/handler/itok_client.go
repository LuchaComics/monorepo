package handler

import (
	"context"
	"log/slog"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/service"
)

type IssuedTokenClientServiceTaskHandler struct {
	config  *config.Config
	logger  *slog.Logger
	service *service.IssuedTokenClientService
}

func NewIssuedTokenClientServiceTaskHandler(
	cfg *config.Config,
	logger *slog.Logger,
	s *service.IssuedTokenClientService,
) *IssuedTokenClientServiceTaskHandler {
	return &IssuedTokenClientServiceTaskHandler{cfg, logger, s}
}

func (h *IssuedTokenClientServiceTaskHandler) Execute(ctx context.Context) error {
	if serviceExecErr := h.service.Execute(ctx); serviceExecErr != nil {
		return serviceExecErr
	}
	return nil
}
