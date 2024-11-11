package usecase

import (
	"log/slog"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/domain"
)

type ListNonFungibleTokensWithFilterByTokenIDsyUseCase struct {
	config *config.Config
	logger *slog.Logger
	repo   domain.NonFungibleTokenRepository
}

func NewListNonFungibleTokensWithFilterByTokenIDsyUseCase(config *config.Config, logger *slog.Logger, repo domain.NonFungibleTokenRepository) *ListNonFungibleTokensWithFilterByTokenIDsyUseCase {
	return &ListNonFungibleTokensWithFilterByTokenIDsyUseCase{config, logger, repo}
}

func (uc *ListNonFungibleTokensWithFilterByTokenIDsyUseCase) Execute(tokIDs []uint64) ([]*domain.NonFungibleToken, error) {
	return uc.repo.ListWithFilterByTokenIDs(tokIDs)
}
