package usecase

import (
	"context"
	"log/slog"

	"github.com/ethereum/go-ethereum/common"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/domain"
)

type CountTokensByOwnerUseCase struct {
	logger *slog.Logger
	repo   domain.TokenRepository
}

func NewCountTokensByOwnerUseCase(logger *slog.Logger, repo domain.TokenRepository) *CountTokensByOwnerUseCase {
	return &CountTokensByOwnerUseCase{logger, repo}
}

func (uc *CountTokensByOwnerUseCase) Execute(ctx context.Context, owner *common.Address) (int64, error) {
	counter, err := uc.repo.CountByOwner(ctx, owner)
	if err != nil {
		return 0, err
	}
	return counter, nil
}
