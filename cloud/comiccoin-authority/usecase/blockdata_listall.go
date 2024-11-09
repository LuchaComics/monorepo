package usecase

import (
	"context"
	"log/slog"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/config"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/domain"
)

type ListAllBlockDataUseCase struct {
	config *config.Configuration
	logger *slog.Logger
	repo   domain.BlockDataRepository
}

func NewListAllBlockDataUseCase(config *config.Configuration, logger *slog.Logger, repo domain.BlockDataRepository) *ListAllBlockDataUseCase {
	return &ListAllBlockDataUseCase{config, logger, repo}
}

func (uc *ListAllBlockDataUseCase) Execute(ctx context.Context) ([]*domain.BlockData, error) {
	data, err := uc.repo.ListAll(ctx)
	if err != nil {
		uc.logger.Error("failed listing all block data", slog.Any("error", err))
		return nil, err
	}
	return data, nil
}
