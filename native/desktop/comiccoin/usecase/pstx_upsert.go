package usecase

import (
	"context"
	"log/slog"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/common/httperror"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/domain"
)

type UpsertPendingSignedTransactionUseCase struct {
	logger *slog.Logger
	repo   domain.PendingSignedTransactionRepository
}

func NewUpsertPendingSignedTransactionUseCase(logger *slog.Logger, repo domain.PendingSignedTransactionRepository) *UpsertPendingSignedTransactionUseCase {
	return &UpsertPendingSignedTransactionUseCase{logger, repo}
}

func (uc *UpsertPendingSignedTransactionUseCase) Execute(ctx context.Context, pstx *domain.PendingSignedTransaction) error {
	//
	// STEP 1: Validation.
	//

	e := make(map[string]string)
	if pstx == nil {
		e["pending_signed_transaction"] = "missing value"
	}
	if len(e) != 0 {
		uc.logger.Warn("Failed validating",
			slog.Any("error", e))
		return httperror.NewForBadRequest(&e)
	}

	//
	// STEP 2: Insert into database.
	//

	return uc.repo.Upsert(pstx)
}
