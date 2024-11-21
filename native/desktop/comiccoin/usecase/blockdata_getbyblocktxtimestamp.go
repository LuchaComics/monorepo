package usecase

import (
	"context"
	"log/slog"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/common/httperror"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/domain"
)

type GetByBlockTransactionTimestampUseCase struct {
	logger *slog.Logger
	repo   domain.BlockDataRepository
}

func NewGetByBlockTransactionTimestampUseCase(logger *slog.Logger, repo domain.BlockDataRepository) *GetByBlockTransactionTimestampUseCase {
	return &GetByBlockTransactionTimestampUseCase{logger, repo}
}

func (uc *GetByBlockTransactionTimestampUseCase) Execute(ctx context.Context, timestamp uint64) (*domain.BlockData, error) {
	//
	// STEP 1: Validation.
	//

	e := make(map[string]string)
	if timestamp == 0 {
		e["timestamp"] = "missing value"
	}
	if len(e) != 0 {
		uc.logger.Warn("Failed validating",
			slog.Any("error", e))
		return nil, httperror.NewForBadRequest(&e)
	}

	//
	// STEP 2: Insert into database.
	//

	return uc.repo.GetByBlockTransactionTimestamp(ctx, timestamp)
}
