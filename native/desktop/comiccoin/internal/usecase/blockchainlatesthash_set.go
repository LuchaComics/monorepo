package usecase

import (
	"log/slog"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/domain"
)

type SetBlockchainLastestHashUseCase struct {
	config *config.Config
	logger *slog.Logger
	repo   domain.BlockchainLastestHashRepository
}

func NewSetBlockchainLastestHashUseCase(config *config.Config, logger *slog.Logger, repo domain.BlockchainLastestHashRepository) *SetBlockchainLastestHashUseCase {
	return &SetBlockchainLastestHashUseCase{config, logger, repo}
}

func (uc *SetBlockchainLastestHashUseCase) Execute(hash string) error {
	return uc.repo.Set(hash)
}
