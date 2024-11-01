package usecase

import (
	"log/slog"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/domain"
)

type GetBlockDataByBlockTransactionTimestampUseCase struct {
	config *config.Config
	logger *slog.Logger
	repo   domain.BlockDataRepository
}

func NewGetBlockDataByBlockTransactionTimestampUseCase(config *config.Config, logger *slog.Logger, repo domain.BlockDataRepository) *GetBlockDataByBlockTransactionTimestampUseCase {
	return &GetBlockDataByBlockTransactionTimestampUseCase{config, logger, repo}
}

func (uc *GetBlockDataByBlockTransactionTimestampUseCase) Execute(nonce uint64) (*domain.BlockData, error) {
	data, err := uc.repo.GetByBlockTransactionTimestamp(nonce)
	if err != nil {
		uc.logger.Error("failed getting block data by block transaction timestamp",
			slog.Any("error", err))
		return nil, err
	}
	return data, nil
}
