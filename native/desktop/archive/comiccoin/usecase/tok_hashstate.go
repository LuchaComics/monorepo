package usecase

import (
	"log/slog"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/domain"
)

type GetTokensHashStateUseCase struct {
	config *config.Config
	logger *slog.Logger
	repo   domain.TokenRepository
}

func NewGetTokensHashStateUseCase(config *config.Config, logger *slog.Logger, repo domain.TokenRepository) *GetTokensHashStateUseCase {
	return &GetTokensHashStateUseCase{config, logger, repo}
}

func (uc *GetTokensHashStateUseCase) Execute() (string, error) {
	return uc.repo.HashState()
}
