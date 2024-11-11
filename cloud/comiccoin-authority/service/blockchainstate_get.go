package service

import (
	"context"
	"log/slog"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/config"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/domain"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/usecase"
)

type GetBlockchainStateService struct {
	config                    *config.Configuration
	logger                    *slog.Logger
	getBlockchainStateUseCase *usecase.GetBlockchainStateUseCase
}

func NewGetBlockchainStateService(
	cfg *config.Configuration,
	logger *slog.Logger,
	uc *usecase.GetBlockchainStateUseCase,
) *GetBlockchainStateService {
	return &GetBlockchainStateService{cfg, logger, uc}
}

func (s *GetBlockchainStateService) Execute(ctx context.Context, chainID uint16) (*domain.BlockchainState, error) {
	return s.getBlockchainStateUseCase.Execute(ctx, chainID)
}
