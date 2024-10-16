package usecase

import (
	"context"
	"log/slog"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/domain"
	"github.com/libp2p/go-libp2p/core/peer"
)

type BlockDataDTOReceiveP2PRequesttUseCase struct {
	config *config.Config
	logger *slog.Logger
	repo   domain.BlockDataDTORepository
}

func NewBlockDataDTOReceiveP2PRequesttUseCase(config *config.Config, logger *slog.Logger, repo domain.BlockDataDTORepository) *BlockDataDTOReceiveP2PRequesttUseCase {
	return &BlockDataDTOReceiveP2PRequesttUseCase{config, logger, repo}
}

func (uc *BlockDataDTOReceiveP2PRequesttUseCase) Execute(ctx context.Context) (peer.ID, string, error) {
	// ctx, cancel := context.WithTimeout(ctx, 45*time.Second) // add a 45-second timeout
	// defer cancel()
	return uc.repo.ReceiveRequestFromNetwork(ctx)
}
