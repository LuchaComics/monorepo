package usecase

import (
	"context"
	"log/slog"

	"github.com/ethereum/go-ethereum/common"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/config"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/domain"
)

type CountTokensByOwnerUseCase struct {
	config *config.Configuration
	logger *slog.Logger
	repo   domain.TokenRepository
}

func NewCountTokensByOwnerUseCase(config *config.Configuration, logger *slog.Logger, repo domain.TokenRepository) *CountTokensByOwnerUseCase {
	return &CountTokensByOwnerUseCase{config, logger, repo}
}

func (uc *CountTokensByOwnerUseCase) Execute(ctx context.Context, owner *common.Address) (int, error) {
	toks, err := uc.repo.ListByOwner(ctx, owner)
	if err != nil {
		return 0, err
	}
	return len(toks), nil
}
