package usecase

import (
	"context"
	"log/slog"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/domain"
)

type BlockDataDTOReceiveP2PResponsetUseCase struct {
	config *config.Config
	logger *slog.Logger
	repo   domain.BlockDataDTORepository
}

func NewBlockDataDTOReceiveP2PResponsetUseCase(config *config.Config, logger *slog.Logger, repo domain.BlockDataDTORepository) *BlockDataDTOReceiveP2PResponsetUseCase {
	return &BlockDataDTOReceiveP2PResponsetUseCase{config, logger, repo}
}

func (uc *BlockDataDTOReceiveP2PResponsetUseCase) Execute(ctx context.Context) (*domain.BlockDataDTO, error) {
	// ctx, cancel := context.WithTimeout(ctx, 45*time.Second) // add a 45-second timeout
	// defer cancel()
	return uc.repo.ReceiveResponseFromNetwork(ctx)
}
