package usecase

import (
	"log/slog"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/domain"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/pkg/httperror"
)

type CreateBlockDataUseCase struct {
	config *config.Config
	logger *slog.Logger
	repo   domain.BlockDataRepository
}

func NewCreateBlockDataUseCase(config *config.Config, logger *slog.Logger, repo domain.BlockDataRepository) *CreateBlockDataUseCase {
	return &CreateBlockDataUseCase{config, logger, repo}
}

func (uc *CreateBlockDataUseCase) Execute(hash string, header *domain.BlockHeader, trans []domain.BlockTransaction) error {
	//
	// STEP 1: Validation.
	//

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
	// STEP 2: Create our strucutre.
	//

	account := &domain.BlockData{
		Hash:   hash,
		Header: header,
		Trans:  trans,
	}

	//
	// STEP 3: Insert into database.
	//

	return uc.repo.Upsert(account)
}
