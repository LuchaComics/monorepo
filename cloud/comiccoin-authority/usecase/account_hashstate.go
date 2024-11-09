package usecase

import (
	"context"
	"log/slog"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/config"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/domain"
)

type GetAccountsHashStateUseCase struct {
	config *config.Configuration
	logger *slog.Logger
	repo   domain.AccountRepository
}

func NewGetAccountsHashStateUseCase(config *config.Configuration, logger *slog.Logger, repo domain.AccountRepository) *GetAccountsHashStateUseCase {
	return &GetAccountsHashStateUseCase{config, logger, repo}
}

func (uc *GetAccountsHashStateUseCase) Execute(ctx context.Context) (string, error) {
	return uc.repo.HashState(ctx)
}
