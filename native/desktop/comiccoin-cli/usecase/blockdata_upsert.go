package usecase

import (
	"context"
	"log/slog"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/common/httperror"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/domain"
)

type UpsertBlockDataUseCase struct {
	logger *slog.Logger
	repo   domain.BlockDataRepository
}

func NewUpsertBlockDataUseCase(logger *slog.Logger, repo domain.BlockDataRepository) *UpsertBlockDataUseCase {
	return &UpsertBlockDataUseCase{logger, repo}
}

func (uc *UpsertBlockDataUseCase) Execute(
	ctx context.Context,
	hash string,
	header *domain.BlockHeader,
	headerSignature []byte,
	trans []domain.BlockTransaction,
	validator *domain.Validator,
) error {
	//
	// STEP 1: Validation.
	// Note: `headerSignature` is optional since PoW algorithm does not require it
	// the PoA algorithm requires it.

	e := make(map[string]string)
	if hash == "" {
		e["hash"] = "missing value"
	}
	if header == nil {
		e["header"] = "missing value"
	}
	if trans == nil {
		e["trans"] = "missing value"
	} else {
		if len(trans) <= 0 {
			e["trans"] = "missing value"
		}
	}
	if len(e) != 0 {
		uc.logger.Warn("Failed creating new block data",
			slog.Any("error", e))
		return httperror.NewForBadRequest(&e)
	}

	//
	// STEP 2: Upsert our strucutre.
	//

	blockData := &domain.BlockData{
		Hash:                 hash,
		Header:               header,
		HeaderSignatureBytes: headerSignature,
		Trans:                trans,
		Validator:            validator,
	}

	//
	// STEP 3: Insert into database.
	//

	return uc.repo.Upsert(ctx, blockData)
}
