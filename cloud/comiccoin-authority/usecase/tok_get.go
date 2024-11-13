package usecase

import (
	"context"
	"log/slog"
	"math/big"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/config"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/domain"
)

type GetTokenUseCase struct {
	config *config.Configuration
	logger *slog.Logger
	repo   domain.TokenRepository
}

func NewGetTokenUseCase(config *config.Configuration, logger *slog.Logger, repo domain.TokenRepository) *GetTokenUseCase {
	return &GetTokenUseCase{config, logger, repo}
}

func (uc *GetTokenUseCase) Execute(ctx context.Context, tokenID *big.Int) (*domain.Token, error) {
	return uc.repo.GetByID(ctx, tokenID)
}
