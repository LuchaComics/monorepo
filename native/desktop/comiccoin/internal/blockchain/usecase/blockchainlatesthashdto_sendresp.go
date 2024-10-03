package usecase

import (
	"context"
	"log/slog"

	"github.com/libp2p/go-libp2p/core/peer"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/domain"
)

type BlockchainLastestHashDTOSendP2PResponseUseCase struct {
	config *config.Config
	logger *slog.Logger
	repo   domain.BlockchainLastestHashDTORepository
}

func NewBlockchainLastestHashDTOSendP2PResponseUseCase(config *config.Config, logger *slog.Logger, repo domain.BlockchainLastestHashDTORepository) *BlockchainLastestHashDTOSendP2PResponseUseCase {
	return &BlockchainLastestHashDTOSendP2PResponseUseCase{config, logger, repo}
}

func (uc *BlockchainLastestHashDTOSendP2PResponseUseCase) Execute(ctx context.Context, peerID peer.ID, data domain.BlockchainLastestHashDTO) error {
	// ctx, cancel := context.WithTimeout(ctx, 45*time.Second) // add a 45-second timeout
	// defer cancel()
	return uc.repo.SendResponseToPeer(ctx, peerID, data)
}
