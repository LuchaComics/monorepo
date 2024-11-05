package usecase

import (
	"log/slog"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/domain"
)

type ListAllSignedIssuedTokenUseCase struct {
	config *config.Config
	logger *slog.Logger
	repo   domain.SignedIssuedTokenRepository
}

func NewListAllSignedIssuedTokenUseCase(config *config.Config, logger *slog.Logger, repo domain.SignedIssuedTokenRepository) *ListAllSignedIssuedTokenUseCase {
	return &ListAllSignedIssuedTokenUseCase{config, logger, repo}
}

func (uc *ListAllSignedIssuedTokenUseCase) Execute() ([]*domain.SignedIssuedToken, error) {
	return uc.repo.ListAll()
}
