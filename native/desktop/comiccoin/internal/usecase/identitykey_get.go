package usecase

import (
	"log/slog"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/domain"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/pkg/httperror"
)

type GetIdentityKeyUseCase struct {
	config *config.Config
	logger *slog.Logger
	repo   domain.IdentityKeyRepository
}

func NewGetIdentityKeyUseCase(config *config.Config, logger *slog.Logger, repo domain.IdentityKeyRepository) *GetIdentityKeyUseCase {
	return &GetIdentityKeyUseCase{config, logger, repo}
}

func (uc *GetIdentityKeyUseCase) Execute(id string) (*domain.IdentityKey, error) {
	//
	// STEP 1: Validation.
	//

	e := make(map[string]string)
	if id == "" {
		e["id"] = "missing value"
	}
	if len(e) != 0 {
		uc.logger.Warn("Failed getting account",
			slog.Any("error", e))
		return nil, httperror.NewForBadRequest(&e)
	}

	//
	// STEP 2: Insert into database.
	//

	return uc.repo.GetByID(id)
}
