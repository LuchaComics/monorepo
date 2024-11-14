package usecase

import (
	"log/slog"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/domain"
)

type GetNonFungibleTokenUseCase struct {
	config *config.Config
	logger *slog.Logger
	repo   domain.NonFungibleTokenRepository
}

func NewGetNonFungibleTokenUseCase(config *config.Config, logger *slog.Logger, repo domain.NonFungibleTokenRepository) *GetNonFungibleTokenUseCase {
	return &GetNonFungibleTokenUseCase{config, logger, repo}
}

func (uc *GetNonFungibleTokenUseCase) Execute(tokenID uint64) (*domain.NonFungibleToken, error) {
	return uc.repo.GetByTokenID(tokenID)
}
