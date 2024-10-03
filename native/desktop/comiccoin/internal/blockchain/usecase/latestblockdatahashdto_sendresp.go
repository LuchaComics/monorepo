package usecase

import (
	"context"
	"log/slog"
	"time"

	"github.com/libp2p/go-libp2p/core/peer"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/domain"
)

type LastBlockDataHashDTOSendP2PResponseUseCase struct {
	config *config.Config
	logger *slog.Logger
	repo   domain.LastBlockDataHashDTORepository
}

func NewLastBlockDataHashDTOSendP2PResponseUseCase(config *config.Config, logger *slog.Logger, repo domain.LastBlockDataHashDTORepository) *LastBlockDataHashDTOSendP2PResponseUseCase {
	return &LastBlockDataHashDTOSendP2PResponseUseCase{config, logger, repo}
}

func (uc *LastBlockDataHashDTOSendP2PResponseUseCase) Execute(ctx context.Context, peerID peer.ID, data domain.LastBlockDataHashDTO) error {
	ctx, cancel := context.WithTimeout(ctx, 45*time.Second) // add a 45-second timeout
	defer cancel()
	return uc.repo.SendResponseToPeer(ctx, peerID, data)
}
