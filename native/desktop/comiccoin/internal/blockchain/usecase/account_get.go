package usecase

import (
	"log/slog"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/domain"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/pkg/httperror"
)

type GetAccountUseCase struct {
	config *config.Config
	logger *slog.Logger
	repo   domain.AccountRepository
}

func NewGetAccountUseCase(config *config.Config, logger *slog.Logger, repo domain.AccountRepository) *GetAccountUseCase {
	return &GetAccountUseCase{config, logger, repo}
}

func (uc *GetAccountUseCase) Execute(id string) (*domain.Account, error) {
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
