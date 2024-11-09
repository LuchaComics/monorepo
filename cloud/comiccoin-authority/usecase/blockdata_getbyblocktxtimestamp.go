package usecase

import (
	"context"
	"log/slog"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/config"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/domain"
)

type GetBlockDataByBlockTransactionTimestampUseCase struct {
	config *config.Configuration
	logger *slog.Logger
	repo   domain.BlockDataRepository
}

func NewGetBlockDataByBlockTransactionTimestampUseCase(config *config.Configuration, logger *slog.Logger, repo domain.BlockDataRepository) *GetBlockDataByBlockTransactionTimestampUseCase {
	return &GetBlockDataByBlockTransactionTimestampUseCase{config, logger, repo}
}

func (uc *GetBlockDataByBlockTransactionTimestampUseCase) Execute(ctx context.Context, nonce uint64) (*domain.BlockData, error) {
	data, err := uc.repo.GetByBlockTransactionTimestamp(ctx, nonce)
	if err != nil {
		uc.logger.Error("failed getting block data by block transaction timestamp",
			slog.Any("error", err))
		return nil, err
	}
	return data, nil
}
