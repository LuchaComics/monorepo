package handler

import (
	"context"
	"log/slog"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/service"
)

type SignedIssuedTokenClientServiceTaskHandler struct {
	config  *config.Config
	logger  *slog.Logger
	service *service.SignedIssuedTokenClientService
}

func NewSignedIssuedTokenClientServiceTaskHandler(
	cfg *config.Config,
	logger *slog.Logger,
	s *service.SignedIssuedTokenClientService,
) *SignedIssuedTokenClientServiceTaskHandler {
	return &SignedIssuedTokenClientServiceTaskHandler{cfg, logger, s}
}

func (h *SignedIssuedTokenClientServiceTaskHandler) Execute(ctx context.Context) error {
	if serviceExecErr := h.service.Execute(ctx); serviceExecErr != nil {
		return serviceExecErr
	}
	return nil
}
