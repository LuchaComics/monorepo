package usecase

import (
	"log/slog"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/domain"
)

type GetLastBlockDataHashUseCase struct {
	config *config.Config
	logger *slog.Logger
	repo   domain.LastBlockDataHashRepository
}

func NewGetLastBlockDataHashUseCase(config *config.Config, logger *slog.Logger, repo domain.LastBlockDataHashRepository) *GetLastBlockDataHashUseCase {
	return &GetLastBlockDataHashUseCase{config, logger, repo}
}

func (uc *GetLastBlockDataHashUseCase) Execute() (domain.LastBlockDataHash, error) {
	return uc.repo.Get()
}
