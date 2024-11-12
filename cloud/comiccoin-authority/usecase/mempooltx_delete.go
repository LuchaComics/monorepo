package usecase

import (
	"context"
	"log/slog"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/config"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/domain"
)

type MempoolTransactionDeleteByChainIDUseCase struct {
	config *config.Configuration
	logger *slog.Logger
	repo   domain.MempoolTransactionRepository
}

func NewMempoolTransactionDeleteByChainIDUseCase(config *config.Configuration, logger *slog.Logger, repo domain.MempoolTransactionRepository) *MempoolTransactionDeleteByChainIDUseCase {
	return &MempoolTransactionDeleteByChainIDUseCase{config, logger, repo}
}

func (uc *MempoolTransactionDeleteByChainIDUseCase) Execute(ctx context.Context, chainID uint16) error {
	err := uc.repo.DeleteByChainID(ctx, chainID)
	if err != nil {
		uc.logger.Error("Failed deleting all mempool transactions",
			slog.Any("error", err))
		return err
	}
	return nil
}
