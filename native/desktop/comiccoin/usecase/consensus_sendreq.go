package usecase

import (
	"context"
	"log/slog"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/domain"
)

type ConsensusMechanismBroadcastRequestToNetworkUseCase struct {
	config *config.Config
	logger *slog.Logger
	repo   domain.ConsensusRepository
}

func NewConsensusMechanismBroadcastRequestToNetworkUseCase(config *config.Config, logger *slog.Logger, repo domain.ConsensusRepository) *ConsensusMechanismBroadcastRequestToNetworkUseCase {
	return &ConsensusMechanismBroadcastRequestToNetworkUseCase{config, logger, repo}
}

func (uc *ConsensusMechanismBroadcastRequestToNetworkUseCase) Execute(ctx context.Context) error {
	// ctx, cancel := context.WithTimeout(ctx, 45*time.Second) // add a 45-second timeout
	// defer cancel()
	return uc.repo.BroadcastRequestToNetwork(ctx)
}
