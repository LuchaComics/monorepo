package usecase

import (
	"context"
	"log/slog"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/config"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/domain"
)

type MempoolTransactionListByChainIDUseCase struct {
	config *config.Configuration
	logger *slog.Logger
	repo   domain.MempoolTransactionRepository
}

func NewMempoolTransactionListByChainIDUseCase(config *config.Configuration, logger *slog.Logger, repo domain.MempoolTransactionRepository) *MempoolTransactionListByChainIDUseCase {
	return &MempoolTransactionListByChainIDUseCase{config, logger, repo}
}

func (uc *MempoolTransactionListByChainIDUseCase) Execute(ctx context.Context, chainID uint16) ([]*domain.MempoolTransaction, error) {
	return uc.repo.ListByChainID(ctx, chainID)
}
