package usecase

import (
	"log/slog"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/domain"
)

// Usecase (UC) represents always having the greatest token ID saved and
// never having any value lesser then the current token ID saved in the
// database, therefore keeping a consistent token ID sequence.
type SetBlockchainLastestTokenIDIfGTEUseCase struct {
	config *config.Config
	logger *slog.Logger
	repo   domain.BlockchainLastestTokenIDRepository
}

func NewSetBlockchainLastestTokenIDIfGTEUseCase(
	config *config.Config,
	logger *slog.Logger,
	repo domain.BlockchainLastestTokenIDRepository,
) *SetBlockchainLastestTokenIDIfGTEUseCase {
	return &SetBlockchainLastestTokenIDIfGTEUseCase{config, logger, repo}
}

func (uc *SetBlockchainLastestTokenIDIfGTEUseCase) Execute(tokenID uint64) error {
	// Developers Note:
	// The following code check the existence of the previous most recent token
	// ID value so we can check and enforce actual latest token ID values get
	// set in the database and nothing less then the current token ID is set.

	latestTokenID, err := uc.repo.Get()
	if err != nil {
		uc.logger.Warn("Failed getting latest token id, automatically setting returned value to zero",
			slog.Any("error", err))
		latestTokenID = 0
	}
	if tokenID >= latestTokenID {
		return uc.repo.Set(tokenID)
	}
	return nil
}
