package usecase

import (
	"context"
	"log/slog"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/domain"
)

type ConsensusMechanismReceiveResponseFromNetworkUseCase struct {
	config *config.Config
	logger *slog.Logger
	repo   domain.ConsensusRepository
}

func NewConsensusMechanismReceiveResponseFromNetworkUseCase(config *config.Config, logger *slog.Logger, repo domain.ConsensusRepository) *ConsensusMechanismReceiveResponseFromNetworkUseCase {
	return &ConsensusMechanismReceiveResponseFromNetworkUseCase{config, logger, repo}
}

func (uc *ConsensusMechanismReceiveResponseFromNetworkUseCase) Execute(ctx context.Context) (string, error) {
	// ctx, cancel := context.WithTimeout(ctx, 45*time.Second) // add a 45-second timeout
	// defer cancel()
	return uc.repo.ReceiveMajorityVoteConsensusResponseFromNetwork(ctx)
}
