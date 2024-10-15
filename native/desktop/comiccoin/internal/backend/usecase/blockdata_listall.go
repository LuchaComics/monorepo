package usecase

import (
	"log/slog"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/backend/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/backend/domain"
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
	data, err := uc.repo.ListAll()
	if err != nil {
		uc.logger.Error("failed listing all block data", slog.Any("error", err))
		return nil, err
	}
	return data, nil
}
