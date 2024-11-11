package usecase

import (
	"context"
	"log/slog"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/config"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/domain"
)

type ListBlockDataUnorderedHashArrayUseCase struct {
	config *config.Configuration
	logger *slog.Logger
	repo   domain.BlockDataRepository
}

func NewListBlockDataUnorderedHashArrayUseCase(config *config.Configuration, logger *slog.Logger, repo domain.BlockDataRepository) *ListBlockDataUnorderedHashArrayUseCase {
	return &ListBlockDataUnorderedHashArrayUseCase{config, logger, repo}
}

func (uc *ListBlockDataUnorderedHashArrayUseCase) Execute(ctx context.Context, chainID uint16) ([]string, error) {
	data, err := uc.repo.ListUnorderedHashArrayForChainID(ctx, chainID)
	if err != nil {
		uc.logger.Error("failed listing all block data", slog.Any("error", err))
		return nil, err
	}
	return data, nil
}
