package usecase

import (
	"context"
	"log/slog"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/domain"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/pkg/httperror"
)

type BroadcastMempoolTransactionDTOUseCase struct {
	config *config.Config
	logger *slog.Logger
	repo   domain.MempoolTransactionDTORepository
}

func NewBroadcastMempoolTransactionDTOUseCase(config *config.Config, logger *slog.Logger, repo domain.MempoolTransactionDTORepository) *BroadcastMempoolTransactionDTOUseCase {
	return &BroadcastMempoolTransactionDTOUseCase{config, logger, repo}
}

func (uc *BroadcastMempoolTransactionDTOUseCase) Execute(ctx context.Context, stx *domain.MempoolTransaction) error {
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

	dto := &domain.MempoolTransactionDTO{
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
