package usecase

import (
	"log/slog"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/domain"
)

type CastVoteForLatestHashInMajorityVotingConsensusUseCase struct {
	config *config.Config
	logger *slog.Logger
	repo   domain.ConsensusRepository
}

func NewCastVoteForLatestHashInMajorityVotingConsensusUseCase(config *config.Config, logger *slog.Logger, repo domain.ConsensusRepository) *CastVoteForLatestHashInMajorityVotingConsensusUseCase {
	return &CastVoteForLatestHashInMajorityVotingConsensusUseCase{config, logger, repo}
}

func (uc *CastVoteForLatestHashInMajorityVotingConsensusUseCase) Execute(latestHash string) error {
	return uc.repo.CastVoteForLatestHashConsensus(latestHash)
}
