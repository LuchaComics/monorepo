package usecase

import (
	"context"
	"log/slog"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/domain"
)

type DownloadFromNetworkBlockDataDTOUseCase struct {
	config *config.Config
	logger *slog.Logger
	repo   domain.BlockDataDTORepository
}

func NewDownloadFromNetworkBlockDataDTOUseCase(config *config.Config, logger *slog.Logger, repo domain.BlockDataDTORepository) *DownloadFromNetworkBlockDataDTOUseCase {
	return &DownloadFromNetworkBlockDataDTOUseCase{config, logger, repo}
}

func (uc *DownloadFromNetworkBlockDataDTOUseCase) Execute(ctx context.Context, blockDataHash string) (*domain.BlockDataDTO, error) {
	return uc.repo.DownloadFromNetwork(ctx, blockDataHash)
}
