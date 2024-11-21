package usecase

import (
	"context"
	"log/slog"

	"github.com/ethereum/go-ethereum/common"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/domain"
)

type ListTokensByOwnerUseCase struct {
	logger *slog.Logger
	repo   domain.TokenRepository
}

func NewListTokensByOwnerUseCase(logger *slog.Logger, repo domain.TokenRepository) *ListTokensByOwnerUseCase {
	return &ListTokensByOwnerUseCase{logger, repo}
}

func (uc *ListTokensByOwnerUseCase) Execute(ctx context.Context, owner *common.Address) ([]*domain.Token, error) {
	return uc.repo.ListByOwner(ctx, owner)
}
