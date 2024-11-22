package usecase

import (
	"context"
	"log/slog"
	"math/big"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/domain"
)

type ListNonFungibleTokensWithFilterByTokenIDsyUseCase struct {
	logger *slog.Logger
	repo   domain.NonFungibleTokenRepository
}

func NewListNonFungibleTokensWithFilterByTokenIDsyUseCase(logger *slog.Logger, repo domain.NonFungibleTokenRepository) *ListNonFungibleTokensWithFilterByTokenIDsyUseCase {
	return &ListNonFungibleTokensWithFilterByTokenIDsyUseCase{logger, repo}
}

func (uc *ListNonFungibleTokensWithFilterByTokenIDsyUseCase) Execute(ctx context.Context, tokIDs []*big.Int) ([]*domain.NonFungibleToken, error) {
	return uc.repo.ListWithFilterByTokenIDs(tokIDs)
}
