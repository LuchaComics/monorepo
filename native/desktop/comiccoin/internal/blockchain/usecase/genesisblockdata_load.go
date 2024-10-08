package usecase

import (
	"log/slog"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/domain"
)

type LoadGenesisBlockDataUseCase struct {
	config *config.Config
	logger *slog.Logger
	repo   domain.GenesisBlockDataRepository
}

func NewLoadGenesisBlockDataUseCase(config *config.Config, logger *slog.Logger, repo domain.GenesisBlockDataRepository) *LoadGenesisBlockDataUseCase {
	return &LoadGenesisBlockDataUseCase{config, logger, repo}
}

func (uc *LoadGenesisBlockDataUseCase) Execute() (*domain.GenesisBlockData, error) {
	return uc.repo.LoadGenesisData()
}
