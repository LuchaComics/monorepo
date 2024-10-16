package usecase

import (
	"context"
	"log/slog"

	"github.com/libp2p/go-libp2p/core/peer"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/domain"
)

type ConsensusMechanismReceiveRequestFromNetworkUseCase struct {
	config *config.Config
	logger *slog.Logger
	repo   domain.ConsensusRepository
}

func NewConsensusMechanismReceiveRequestFromNetworkUseCase(config *config.Config, logger *slog.Logger, repo domain.ConsensusRepository) *ConsensusMechanismReceiveRequestFromNetworkUseCase {
	return &ConsensusMechanismReceiveRequestFromNetworkUseCase{config, logger, repo}
}

func (uc *ConsensusMechanismReceiveRequestFromNetworkUseCase) Execute(ctx context.Context) (peer.ID, error) {
	// ctx, cancel := context.WithTimeout(ctx, 45*time.Second) // add a 45-second timeout
	// defer cancel()
	return uc.repo.ReceiveRequestFromNetwork(ctx)
}
