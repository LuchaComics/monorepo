package usecase

import (
	"context"
	"log/slog"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/domain"
)

type ListPendingSignedTransactionUseCase struct {
	logger *slog.Logger
	repo   domain.PendingSignedTransactionRepository
}

func NewListPendingSignedTransactionUseCase(logger *slog.Logger, repo domain.PendingSignedTransactionRepository) *ListPendingSignedTransactionUseCase {
	return &ListPendingSignedTransactionUseCase{logger, repo}
}

func (uc *ListPendingSignedTransactionUseCase) Execute(ctx context.Context) ([]*domain.PendingSignedTransaction, error) {
	return uc.repo.ListAll()
}
