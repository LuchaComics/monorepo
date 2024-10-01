package usecase

import (
	"log/slog"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/domain"
)

type ListAllMempoolTransactionUseCase struct {
	config *config.Config
	logger *slog.Logger
	repo   domain.MempoolTransactionRepository
}

func NewListAllMempoolTransactionUseCase(config *config.Config, logger *slog.Logger, repo domain.MempoolTransactionRepository) *ListAllMempoolTransactionUseCase {
	return &ListAllMempoolTransactionUseCase{config, logger, repo}
}

func (uc *ListAllMempoolTransactionUseCase) Execute() ([]*domain.MempoolTransaction, error) {
	return uc.repo.ListAll()
}
