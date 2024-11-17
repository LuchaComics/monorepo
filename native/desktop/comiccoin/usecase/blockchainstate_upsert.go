package usecase

import (
	"context"
	"log/slog"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/common/httperror"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/domain"
)

type UpsertBlockchainStateUseCase struct {
	logger *slog.Logger
	repo   domain.BlockchainStateRepository
}

func NewUpsertBlockchainStateUseCase(logger *slog.Logger, repo domain.BlockchainStateRepository) *UpsertBlockchainStateUseCase {
	return &UpsertBlockchainStateUseCase{logger, repo}
}

func (uc *UpsertBlockchainStateUseCase) Execute(ctx context.Context, bcs *domain.BlockchainState) error {
	//
	// STEP 1: Validation.
	//

	e := make(map[string]string)
	if bcs == nil {
		e["blockchain_state"] = "missing value"
	}
	if len(e) != 0 {
		uc.logger.Warn("Failed upserting blockchain state",
			slog.Any("error", e))
		return httperror.NewForBadRequest(&e)
	}

	//
	// STEP 3: Insert into database.
	//

	return uc.repo.UpsertByChainID(ctx, bcs)
}
