package usecase

import (
	"context"
	"log/slog"
	"math/big"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/common/httperror"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/domain"
)

type DeletePendingSignedTransactionUseCase struct {
	logger *slog.Logger
	repo   domain.PendingSignedTransactionRepository
}

func NewDeletePendingSignedTransactionUseCase(logger *slog.Logger, repo domain.PendingSignedTransactionRepository) *DeletePendingSignedTransactionUseCase {
	return &DeletePendingSignedTransactionUseCase{logger, repo}
}

func (uc *DeletePendingSignedTransactionUseCase) Execute(ctx context.Context, nonce *big.Int) error {
	//
	// STEP 1: Validation.
	//

	e := make(map[string]string)
	if nonce == nil {
		e["nonce"] = "missing value"
	} else {
		// Special thanks to "Is there another way of testing if a big.Int is 0?" via https://stackoverflow.com/a/64257532
		if len(nonce.Bits()) == 0 {
			e["nonce"] = "missing value"
		}
	}
	if len(e) != 0 {
		uc.logger.Warn("Failed validating",
			slog.Any("error", e))
		return httperror.NewForBadRequest(&e)
	}

	//
	// STEP 2: Delete from database.
	//

	return uc.repo.DeleteByNonce(nonce)
}
