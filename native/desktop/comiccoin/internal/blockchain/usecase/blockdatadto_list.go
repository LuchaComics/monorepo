package usecase

import (
	"context"
	"log/slog"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/domain"
)

type ListLatestBlockDataAfterHashDTOUseCase struct {
	config *config.Config
	logger *slog.Logger
	repo   domain.BlockDataDTORepository
}

func NewListLatestBlockDataAfterHashDTOUseCase(config *config.Config, logger *slog.Logger, repo domain.BlockDataDTORepository) *ListLatestBlockDataAfterHashDTOUseCase {
	return &ListLatestBlockDataAfterHashDTOUseCase{config, logger, repo}
}

func (uc *ListLatestBlockDataAfterHashDTOUseCase) Execute(ctx context.Context, hash string) ([]*domain.BlockDataDTO, error) {
	return uc.repo.ListLatestAfterHash(ctx, hash)
}
