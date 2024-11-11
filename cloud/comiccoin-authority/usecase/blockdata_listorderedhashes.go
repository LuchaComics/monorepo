package usecase

import (
	"context"
	"log/slog"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/config"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/domain"
)

type ListAllBlockNumberByHashArrayUseCase struct {
	config *config.Configuration
	logger *slog.Logger
	repo   domain.BlockDataRepository
}

func NewListAllBlockNumberByHashArrayUseCase(config *config.Configuration, logger *slog.Logger, repo domain.BlockDataRepository) *ListAllBlockNumberByHashArrayUseCase {
	return &ListAllBlockNumberByHashArrayUseCase{config, logger, repo}
}

func (uc *ListAllBlockNumberByHashArrayUseCase) Execute(ctx context.Context, chainID uint16) ([]domain.BlockNumberByHash, error) {
	data, err := uc.repo.ListBlockNumberByHashArrayForChainID(ctx, chainID)
	if err != nil {
		uc.logger.Error("failed listing all block data", slog.Any("error", err))
		return nil, err
	}
	return data, nil
}
