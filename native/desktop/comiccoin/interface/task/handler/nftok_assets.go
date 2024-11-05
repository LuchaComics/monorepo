package handler

import (
	"context"
	"log/slog"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/service"
)

type NonFungibleTokenAssetsServiceTaskHandler struct {
	config  *config.Config
	logger  *slog.Logger
	service *service.NonFungibleTokenAssetsService
}

func NewNonFungibleTokenAssetsServiceTaskHandler(
	cfg *config.Config,
	logger *slog.Logger,
	s *service.NonFungibleTokenAssetsService,
) *NonFungibleTokenAssetsServiceTaskHandler {
	return &NonFungibleTokenAssetsServiceTaskHandler{cfg, logger, s}
}

func (h *NonFungibleTokenAssetsServiceTaskHandler) Execute(ctx context.Context) error {
	if serviceExecErr := h.service.Execute(); serviceExecErr != nil {
		return serviceExecErr
	}
	return nil
}
