package usecase

import (
	"log/slog"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/domain"
)

type ListAllPendingBlockTransactionUseCase struct {
	config *config.Config
	logger *slog.Logger
	repo   domain.PendingBlockTransactionRepository
}

func NewListAllPendingBlockTransactionUseCase(config *config.Config, logger *slog.Logger, repo domain.PendingBlockTransactionRepository) *ListAllPendingBlockTransactionUseCase {
	return &ListAllPendingBlockTransactionUseCase{config, logger, repo}
}

func (uc *ListAllPendingBlockTransactionUseCase) Execute() ([]*domain.PendingBlockTransaction, error) {
	return uc.repo.ListAll()
}
