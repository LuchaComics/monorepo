package usecase

import (
	"log/slog"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/domain"
)

type SetLastBlockDataHashUseCase struct {
	config *config.Config
	logger *slog.Logger
	repo   domain.LastBlockDataHashRepository
}

func NewSetLastBlockDataHashUseCase(config *config.Config, logger *slog.Logger, repo domain.LastBlockDataHashRepository) *SetLastBlockDataHashUseCase {
	return &SetLastBlockDataHashUseCase{config, logger, repo}
}

func (uc *SetLastBlockDataHashUseCase) Execute(hash string) error {
	return uc.repo.Set(hash)
}
