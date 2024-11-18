package usecase

import (
	"log/slog"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/domain"
)

type ListNonFungibleTokensWithFilterByTokenIDsyUseCase struct {
	logger *slog.Logger
	repo   domain.NonFungibleTokenRepository
}

func NewListNonFungibleTokensWithFilterByTokenIDsyUseCase(logger *slog.Logger, repo domain.NonFungibleTokenRepository) *ListNonFungibleTokensWithFilterByTokenIDsyUseCase {
	return &ListNonFungibleTokensWithFilterByTokenIDsyUseCase{logger, repo}
}

func (uc *ListNonFungibleTokensWithFilterByTokenIDsyUseCase) Execute(tokIDs []uint64) ([]*domain.NonFungibleToken, error) {
	return uc.repo.ListWithFilterByTokenIDs(tokIDs)
}
