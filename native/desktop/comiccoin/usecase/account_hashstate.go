package usecase

import (
	"log/slog"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/domain"
)

type GetAccountsHashStateUseCase struct {
	config *config.Config
	logger *slog.Logger
	repo   domain.AccountRepository
}

func NewGetAccountsHashStateUseCase(config *config.Config, logger *slog.Logger, repo domain.AccountRepository) *GetAccountsHashStateUseCase {
	return &GetAccountsHashStateUseCase{config, logger, repo}
}

func (uc *GetAccountsHashStateUseCase) Execute() (string, error) {
	return uc.repo.HashState()
}
