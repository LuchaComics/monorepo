package usecase

import (
	"log/slog"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/domain"
)

type ListAllBlockDataUseCase struct {
	config *config.Config
	logger *slog.Logger
	repo   domain.BlockDataRepository
}

func NewListAllBlockDataUseCase(config *config.Config, logger *slog.Logger, repo domain.BlockDataRepository) *ListAllBlockDataUseCase {
	return &ListAllBlockDataUseCase{config, logger, repo}
}

func (uc *ListAllBlockDataUseCase) Execute() ([]*domain.BlockData, error) {
	return nil, nil
}
