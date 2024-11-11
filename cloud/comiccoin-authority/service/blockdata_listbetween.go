package service

import (
	"context"
	"log/slog"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/config"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/domain"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/usecase"
)

type ListBlockDataFilteredBetweenBlockNumbersInChainIDService struct {
	config                                          *config.Configuration
	logger                                          *slog.Logger
	listBlockDataFilteredBetweenBlockNumbersUseCase *usecase.ListBlockDataFilteredBetweenBlockNumbersUseCase
}

func NewListBlockDataFilteredBetweenBlockNumbersInChainIDService(
	cfg *config.Configuration,
	logger *slog.Logger,
	uc *usecase.ListBlockDataFilteredBetweenBlockNumbersUseCase,
) *ListBlockDataFilteredBetweenBlockNumbersInChainIDService {
	return &ListBlockDataFilteredBetweenBlockNumbersInChainIDService{cfg, logger, uc}
}

func (s *ListBlockDataFilteredBetweenBlockNumbersInChainIDService) Execute(ctx context.Context, startBlockNumber uint64, finishBlockNumber uint64, chainID uint16) ([]*domain.BlockData, error) {
	data, err := s.listBlockDataFilteredBetweenBlockNumbersUseCase.Execute(ctx, startBlockNumber, finishBlockNumber, chainID)
	if err != nil {
		s.logger.Error("Failed listing", slog.Any("error", err))
		return nil, err
	}
	return data, nil
}
