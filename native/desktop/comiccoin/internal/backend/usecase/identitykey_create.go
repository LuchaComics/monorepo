package usecase

import (
	"log/slog"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/backend/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/backend/domain"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/pkg/httperror"
)

type CreateIdentityKeyUseCase struct {
	config *config.Config
	logger *slog.Logger
	repo   domain.IdentityKeyRepository
}

func NewCreateIdentityKeyUseCase(config *config.Config, logger *slog.Logger, repo domain.IdentityKeyRepository) *CreateIdentityKeyUseCase {
	return &CreateIdentityKeyUseCase{config, logger, repo}
}

func (uc *CreateIdentityKeyUseCase) Execute(id string) error {
	//
	// STEP 1: Validation.
	//

	e := make(map[string]string)
	if id == "" {
		e["id"] = "missing value"
	}
	if len(e) != 0 {
		uc.logger.Warn("Failed creating new identity key",
			slog.Any("error", e))
		return httperror.NewForBadRequest(&e)
	}

	//
	// STEP 2: Create our strucutre.
	//

	ik, err := domain.NewIdentityKey(id)
	if err != nil {
		uc.logger.Warn("Failed creating new identity key",
			slog.Any("error", e))
		return err
	}

	//
	// STEP 3: Insert into database.
	//

	return uc.repo.Upsert(ik)
}
