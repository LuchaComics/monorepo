package service

import (
	"context"
	"log/slog"
	"time"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/usecase"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/pkg/kmutexutil"
)

// MiningValidationService represents (TODO)
type MiningValidationService struct {
	config                             *config.Config
	logger                             *slog.Logger
	kmutex                             kmutexutil.KMutexProvider
	receiveProposedBlockDataDTOUseCase *usecase.ReceiveProposedBlockDataDTOUseCase
	createBlockDataUseCase             *usecase.CreateBlockDataUseCase
}

func NewMiningValidationService(
	cfg *config.Config,
	logger *slog.Logger,
	kmutex kmutexutil.KMutexProvider,
	uc1 *usecase.ReceiveProposedBlockDataDTOUseCase,
	uc2 *usecase.CreateBlockDataUseCase,
) *MiningValidationService {
	return &MiningValidationService{cfg, logger, kmutex, uc1, uc2}
}

func (s *MiningValidationService) Execute(ctx context.Context) error {

	//
	// STEP 1
	// Wait to receive data (which also was validated) from the P2P network.
	//

	proposedBlockData, err := s.receiveProposedBlockDataDTOUseCase.Execute(ctx)
	if err != nil {
		s.logger.Error("validator failed receiving dto",
			slog.Any("error", err))
		return err
	}
	if proposedBlockData == nil {
		// Developer Note:
		// If we haven't received anything, that means we haven't connected to
		// the distributed / P2P network, so all we can do at the moment is to
		// pause the execution for 1 second and then retry again.
		time.Sleep(1 * time.Second)
		return nil
	}

	s.logger.Debug("received dto from network",
		slog.Any("hash", proposedBlockData.Hash),
	)

	// Lock the validator's database so we coordinate when we receive, validate
	// and/or save to the database.
	s.kmutex.Acquire("validator-service")
	defer s.kmutex.Release("validator-service")

	//
	// STEP 2
	// Validate our proposed block data to our blockchain.
	//

	//TODO: IMPL.

	//
	// STEP 3:
	// Save to our local database.
	//

	if err := s.createBlockDataUseCase.Execute(proposedBlockData.Hash, proposedBlockData.Header, proposedBlockData.Trans); err != nil {
		s.logger.Error("validator failed saving block data",
			slog.Any("error", err))
		return err
	}

	s.logger.Debug("saved to validator",
		slog.Any("hash", proposedBlockData.Hash),
	)

	return nil
}
