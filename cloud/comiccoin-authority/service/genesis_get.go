package service

import (
	"context"
	"log/slog"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/config"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/domain"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/usecase"
)

type GetGenesisBlockDataService struct {
	config                     *config.Configuration
	logger                     *slog.Logger
	getGenesisBlockDataUseCase *usecase.GetGenesisBlockDataUseCase
}

func NewGetGenesisBlockDataService(
	cfg *config.Configuration,
	logger *slog.Logger,
	uc *usecase.GetGenesisBlockDataUseCase,
) *GetGenesisBlockDataService {
	return &GetGenesisBlockDataService{cfg, logger, uc}
}

func (s *GetGenesisBlockDataService) Execute(ctx context.Context, chainID uint16) (*domain.GenesisBlockData, error) {
	return s.getGenesisBlockDataUseCase.Execute(ctx, chainID)
}
