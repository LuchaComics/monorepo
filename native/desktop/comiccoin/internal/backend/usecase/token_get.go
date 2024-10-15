package usecase

import (
	"log/slog"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/backend/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/backend/domain"
)

type GetTokenUseCase struct {
	config *config.Config
	logger *slog.Logger
	repo   domain.TokenRepository
}

func NewGetTokenUseCase(config *config.Config, logger *slog.Logger, repo domain.TokenRepository) *GetTokenUseCase {
	return &GetTokenUseCase{config, logger, repo}
}

func (uc *GetTokenUseCase) Execute(tokenID uint64) (*domain.Token, error) {
	return uc.repo.GetByID(tokenID)
}
