package usecase

import (
	"context"
	"log/slog"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/domain"
)

type UploadToNetworkBlockDataDTOUseCase struct {
	config *config.Config
	logger *slog.Logger
	repo   domain.BlockDataDTORepository
}

func NewUploadToNetworkBlockDataDTOUseCase(config *config.Config, logger *slog.Logger, repo domain.BlockDataDTORepository) *UploadToNetworkBlockDataDTOUseCase {
	return &UploadToNetworkBlockDataDTOUseCase{config, logger, repo}
}

func (uc *UploadToNetworkBlockDataDTOUseCase) Execute(ctx context.Context, data *domain.BlockDataDTO) error {
	return uc.repo.UploadToNetwork(ctx, data)
}
