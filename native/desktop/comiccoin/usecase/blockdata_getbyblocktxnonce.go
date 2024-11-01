package usecase

import (
	"log/slog"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/domain"
)

type GetByBlockTransactionNonceUseCase struct {
	config *config.Config
	logger *slog.Logger
	repo   domain.BlockDataRepository
}

func NewGetByBlockTransactionNonceUseCase(config *config.Config, logger *slog.Logger, repo domain.BlockDataRepository) *GetByBlockTransactionNonceUseCase {
	return &GetByBlockTransactionNonceUseCase{config, logger, repo}
}

func (uc *GetByBlockTransactionNonceUseCase) Execute(nonce uint64) (*domain.BlockData, error) {
	data, err := uc.repo.GetByBlockTransactionNonce(nonce)
	if err != nil {
		uc.logger.Error("failed getting block data by block transaction nonce",
			slog.Any("error", err))
		return nil, err
	}
	return data, nil
}
