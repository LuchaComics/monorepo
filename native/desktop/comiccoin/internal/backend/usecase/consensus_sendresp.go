package usecase

import (
	"context"
	"log/slog"

	"github.com/libp2p/go-libp2p/core/peer"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/backend/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/backend/domain"
)

type ConsensusMechanismSendResponseToPeerUseCase struct {
	config *config.Config
	logger *slog.Logger
	repo   domain.ConsensusRepository
}

func NewConsensusMechanismSendResponseToPeerUseCase(config *config.Config, logger *slog.Logger, repo domain.ConsensusRepository) *ConsensusMechanismSendResponseToPeerUseCase {
	return &ConsensusMechanismSendResponseToPeerUseCase{config, logger, repo}
}

func (uc *ConsensusMechanismSendResponseToPeerUseCase) Execute(ctx context.Context, peerID peer.ID, blockchainHash string) error {
	// ctx, cancel := context.WithTimeout(ctx, 45*time.Second) // add a 45-second timeout
	// defer cancel()
	return uc.repo.SendResponseToPeer(ctx, peerID, blockchainHash)
}
