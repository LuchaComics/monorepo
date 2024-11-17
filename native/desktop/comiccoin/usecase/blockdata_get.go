package usecase

import (
	"context"
	"log/slog"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/common/httperror"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/domain"
)

type GetBlockDataUseCase struct {
	logger *slog.Logger
	repo   domain.BlockDataRepository
}

func NewGetBlockDataUseCase(logger *slog.Logger, repo domain.BlockDataRepository) *GetBlockDataUseCase {
	return &GetBlockDataUseCase{logger, repo}
}

func (uc *GetBlockDataUseCase) Execute(ctx context.Context, hash string) (*domain.BlockData, error) {
	//
	// STEP 1: Validation.
	//

	e := make(map[string]string)
	if hash == "" {
		e["hash"] = "missing value"
	}
	if len(e) != 0 {
		uc.logger.Warn("Failed getting account",
			slog.Any("error", e))
		return nil, httperror.NewForBadRequest(&e)
	}

	//
	// STEP 2: Insert into database.
	//

	return uc.repo.GetByHash(ctx, hash)
}
