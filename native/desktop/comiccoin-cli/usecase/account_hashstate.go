package usecase

import (
	"context"
	"log/slog"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/domain"
)

type GetAccountsHashStateUseCase struct {
	logger *slog.Logger
	repo   domain.AccountRepository
}

func NewGetAccountsHashStateUseCase(logger *slog.Logger, repo domain.AccountRepository) *GetAccountsHashStateUseCase {
	return &GetAccountsHashStateUseCase{logger, repo}
}

func (uc *GetAccountsHashStateUseCase) Execute(ctx context.Context) (string, error) {
	return uc.repo.HashState(ctx)
}
