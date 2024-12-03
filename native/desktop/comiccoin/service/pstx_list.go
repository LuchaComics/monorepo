package service

import (
	"context"
	"log/slog"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/domain"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/usecase"
)

type PendingSignedTransactionListService struct {
	logger                              *slog.Logger
	listPendingSignedTransactionUseCase *usecase.ListPendingSignedTransactionUseCase
}

func NewPendingSignedTransactionListService(
	logger *slog.Logger,
	uc1 *usecase.ListPendingSignedTransactionUseCase,
) *PendingSignedTransactionListService {
	return &PendingSignedTransactionListService{logger, uc1}
}

func (s *PendingSignedTransactionListService) Execute(ctx context.Context) ([]*domain.PendingSignedTransaction, error) {
	return s.listPendingSignedTransactionUseCase.Execute(ctx)
}
