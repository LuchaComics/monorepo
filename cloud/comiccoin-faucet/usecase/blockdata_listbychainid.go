package usecase

import (
	"context"
	"log/slog"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/config"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/domain"
)

//
// Copied from `github.com/LuchaComics/monorepo/cloud/comiccoin-authority/usecase`
//

type ListBlockDataByChainIDUseCase struct {
	config *config.Configuration
	logger *slog.Logger
	repo   domain.BlockDataRepository
}

func NewListBlockDataByChainIDUseCase(config *config.Configuration, logger *slog.Logger, repo domain.BlockDataRepository) *ListBlockDataByChainIDUseCase {
	return &ListBlockDataByChainIDUseCase{config, logger, repo}
}

func (uc *ListBlockDataByChainIDUseCase) Execute(ctx context.Context) ([]*domain.BlockData, error) {
	data, err := uc.repo.ListByChainID(ctx, uc.config.Blockchain.ChainID)
	if err != nil {
		uc.logger.Error("failed listing all block data", slog.Any("error", err))
		return nil, err
	}
	return data, nil
}
