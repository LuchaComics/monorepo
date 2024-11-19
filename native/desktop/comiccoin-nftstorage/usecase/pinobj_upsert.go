package usecase

import (
	"context"
	"log/slog"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin-nftstorage/common/httperror"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin-nftstorage/domain"
)

type UpsertPinObjectUseCase struct {
	logger *slog.Logger
	repo   domain.PinObjectRepository
}

func NewUpsertPinObjectUseCase(logger *slog.Logger, r1 domain.PinObjectRepository) *UpsertPinObjectUseCase {
	return &UpsertPinObjectUseCase{logger, r1}
}

func (uc *UpsertPinObjectUseCase) Execute(ctx context.Context, pinobj *domain.PinObject) error {
	//
	// STEP 1:
	// Validation.
	//

	e := make(map[string]string)

	if pinobj.CID == "" {
		e["cid"] = "missing value"
	}
	if pinobj.RequestID == 0 {
		e["requestid"] = "missing value"
	}
	if pinobj.Meta == nil {
		e["meta"] = "missing value"
	} else {
		if pinobj.Meta["filename"] == "" {
			e["meta"] = "missing `filename` value"
		}
		if pinobj.Meta["content_type"] == "" {
			e["meta"] = "missing `content_type` value"
		}
	}
	if len(e) != 0 {
		uc.logger.Warn("Validation failed", slog.Any("e", e))
		return httperror.NewForBadRequest(&e)
	}

	// Save to database.
	if err := uc.repo.Upsert(pinobj); err != nil {
		uc.logger.Error("database upsert error", slog.Any("error", err))
		return err
	}
	return nil
}
