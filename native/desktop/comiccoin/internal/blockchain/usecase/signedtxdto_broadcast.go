package usecase

import (
	"context"
	"log/slog"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/domain"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/pkg/httperror"
)

type BroadcastSignedTransactionDTOUseCase struct {
	config *config.Config
	logger *slog.Logger
	repo   domain.SignedTransactionDTORepository
}

func NewBroadcastSignedTransactionDTOUseCase(config *config.Config, logger *slog.Logger, repo domain.SignedTransactionDTORepository) *BroadcastSignedTransactionDTOUseCase {
	return &BroadcastSignedTransactionDTOUseCase{config, logger, repo}
}

func (uc *BroadcastSignedTransactionDTOUseCase) Execute(ctx context.Context, stx *domain.SignedTransaction) error {
	//
	// STEP 1: Validation.
	//

	e := make(map[string]string)
	if stx == nil {
		e["signed_transaction"] = "missing value"
	}
	if len(e) != 0 {
		uc.logger.Warn("Validation failed",
			slog.Any("error", e))
		return httperror.NewForBadRequest(&e)
	}

	//
	// STEP 2: Create our strucutre.
	//

	dto := &domain.SignedTransactionDTO{
		Transaction: stx.Transaction,
		V:           stx.V,
		R:           stx.R,
		S:           stx.S,
	}

	//
	// STEP 3: Insert into database.
	//

	return uc.repo.BroadcastToP2PNetwork(ctx, dto)
}
