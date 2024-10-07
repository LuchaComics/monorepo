package service

import (
	"context"
	"log/slog"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/pkg/kmutexutil"
)

type ConsensusService struct {
	config *config.Config
	logger *slog.Logger
	kmutex kmutexutil.KMutexProvider
}

func NewConsensusService(
	cfg *config.Config,
	logger *slog.Logger,
	kmutex kmutexutil.KMutexProvider,

) *ConsensusService {
	return &ConsensusService{cfg, logger, kmutex}
}

func (s *ConsensusService) Execute(ctx context.Context) error {
	s.logger.Debug("starting consensus service...")
	defer s.logger.Debug("finished consensus service")
	//TODO IMPL.

	return nil
}
