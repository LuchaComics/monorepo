package usecase

import (
	"context"
	"log/slog"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/domain"
	"github.com/libp2p/go-libp2p/core/peer"
)

type BlockDataDTOSendP2PResponsetUseCase struct {
	config *config.Config
	logger *slog.Logger
	repo   domain.BlockDataDTORepository
}

func NewBlockDataDTOSendP2PResponsetUseCase(config *config.Config, logger *slog.Logger, repo domain.BlockDataDTORepository) *BlockDataDTOSendP2PResponsetUseCase {
	return &BlockDataDTOSendP2PResponsetUseCase{config, logger, repo}
}

func (uc *BlockDataDTOSendP2PResponsetUseCase) Execute(ctx context.Context, peerID peer.ID, data *domain.BlockDataDTO) error {
	// ctx, cancel := context.WithTimeout(ctx, 45*time.Second) // add a 45-second timeout
	// defer cancel()
	return uc.repo.SendResponseToPeer(ctx, peerID, data)
}
