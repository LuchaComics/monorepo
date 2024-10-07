package usecase

import (
	"context"
	"log/slog"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/domain"
)

type QueryLatestHashByMajorityVotingConsensusUseCase struct {
	config *config.Config
	logger *slog.Logger
	repo   domain.ConsensusRepository
}

func NewQueryLatestHashByMajorityVotingConsensusUseCase(config *config.Config, logger *slog.Logger, repo domain.ConsensusRepository) *QueryLatestHashByMajorityVotingConsensusUseCase {
	return &QueryLatestHashByMajorityVotingConsensusUseCase{config, logger, repo}
}

func (uc *QueryLatestHashByMajorityVotingConsensusUseCase) Execute(ctx context.Context) (string, error) {
	return uc.repo.QueryLatestHashByConsensus(ctx)
}
