package usecase

import (
	"context"
	"log/slog"
	"math/big"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/domain"
)

type GetNonFungibleTokenUseCase struct {
	logger *slog.Logger
	repo   domain.NonFungibleTokenRepository
}

func NewGetNonFungibleTokenUseCase(logger *slog.Logger, repo domain.NonFungibleTokenRepository) *GetNonFungibleTokenUseCase {
	return &GetNonFungibleTokenUseCase{logger, repo}
}

func (uc *GetNonFungibleTokenUseCase) Execute(ctx context.Context, tokenID *big.Int) (*domain.NonFungibleToken, error) {
	return uc.repo.GetByTokenID(tokenID)
}
