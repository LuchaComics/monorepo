package usecase

import (
	"log/slog"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/domain"
)

type ListAllSignedTransactionUseCase struct {
	config *config.Config
	logger *slog.Logger
	repo   domain.SignedTransactionRepository
}

func NewListAllSignedTransactionUseCase(config *config.Config, logger *slog.Logger, repo domain.SignedTransactionRepository) *ListAllSignedTransactionUseCase {
	return &ListAllSignedTransactionUseCase{config, logger, repo}
}

func (uc *ListAllSignedTransactionUseCase) Execute() ([]*domain.SignedTransaction, error) {
	return uc.repo.ListAll()
}
