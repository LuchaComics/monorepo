package usecase

import (
	"context"
	"log/slog"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/config"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/domain"
)

type ListBlockDataFilteredBetweenBlockNumbersUseCase struct {
	config *config.Configuration
	logger *slog.Logger
	repo   domain.BlockDataRepository
}

func NewListBlockDataFilteredBetweenBlockNumbersUseCase(config *config.Configuration, logger *slog.Logger, repo domain.BlockDataRepository) *ListBlockDataFilteredBetweenBlockNumbersUseCase {
	return &ListBlockDataFilteredBetweenBlockNumbersUseCase{config, logger, repo}
}

func (uc *ListBlockDataFilteredBetweenBlockNumbersUseCase) Execute(ctx context.Context, startBlockNumber uint64, finishBlockNumber uint64, chainID uint16) ([]*domain.BlockData, error) {
	data, err := uc.repo.ListInBetweenBlockNumbersForChainID(ctx, startBlockNumber, finishBlockNumber, chainID)
	if err != nil {
		uc.logger.Error("failed listing all block data", slog.Any("error", err))
		return nil, err
	}
	return data, nil
}
