package usecase

import (
	"context"
	"log/slog"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/config"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/domain"
)

type ListBlockDataFilteredInHashesUseCase struct {
	config *config.Configuration
	logger *slog.Logger
	repo   domain.BlockDataRepository
}

func NewListBlockDataFilteredInHashesUseCase(config *config.Configuration, logger *slog.Logger, repo domain.BlockDataRepository) *ListBlockDataFilteredInHashesUseCase {
	return &ListBlockDataFilteredInHashesUseCase{config, logger, repo}
}

func (uc *ListBlockDataFilteredInHashesUseCase) Execute(ctx context.Context, hashes []string) ([]*domain.BlockData, error) {
	data, err := uc.repo.ListInHashes(ctx, hashes)
	if err != nil {
		uc.logger.Error("failed listing all block data", slog.Any("error", err))
		return nil, err
	}
	return data, nil
}
