package usecase

import (
	"log/slog"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/domain"
)

type GetSignedIssuedTokenUseCase struct {
	config *config.Config
	logger *slog.Logger
	repo   domain.SignedIssuedTokenRepository
}

func NewGetSignedIssuedTokenUseCase(config *config.Config, logger *slog.Logger, repo domain.SignedIssuedTokenRepository) *GetSignedIssuedTokenUseCase {
	return &GetSignedIssuedTokenUseCase{config, logger, repo}
}

func (uc *GetSignedIssuedTokenUseCase) Execute(id uint64) (*domain.SignedIssuedToken, error) {
	return uc.repo.GetByID(id)
}
