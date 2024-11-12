package usecase

import (
	"log/slog"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/common/httperror"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/domain"
)

type UpsertNonFungibleTokenUseCase struct {
	config *config.Config
	logger *slog.Logger
	repo   domain.NonFungibleTokenRepository
}

func NewUpsertNonFungibleTokenUseCase(config *config.Config, logger *slog.Logger, repo domain.NonFungibleTokenRepository) *UpsertNonFungibleTokenUseCase {
	return &UpsertNonFungibleTokenUseCase{config, logger, repo}
}

func (uc *UpsertNonFungibleTokenUseCase) Execute(nftok *domain.NonFungibleToken) error {
	//
	// STEP 1: Validation.
	//

	e := make(map[string]string)
	if nftok == nil {
		e["nftok"] = "missing value"
	}
	if len(e) != 0 {
		uc.logger.Warn("Failed creating non-fungible token because validation failed",
			slog.Any("error", e))
		return httperror.NewForBadRequest(&e)
	}

	//
	// STEP 2: Insert into database.
	//

	return uc.repo.Upsert(nftok)
}