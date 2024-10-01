package service

import (
	"context"
	"log/slog"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/usecase"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/pkg/kmutexutil"
)

type MiningService struct {
	config                                  *config.Config
	logger                                  *slog.Logger
	kmutex                                  kmutexutil.KMutexProvider
	listAllPendingBlockTransactionUseCase   *usecase.ListAllPendingBlockTransactionUseCase
	getLastBlockDataHashUseCase             *usecase.GetLastBlockDataHashUseCase
	setLastBlockDataHashUseCase             *usecase.SetLastBlockDataHashUseCase
	getBlockDataUseCase                     *usecase.GetBlockDataUseCase
	createBlockDataUseCase                  *usecase.CreateBlockDataUseCase
	proofOfWorkUseCase                      *usecase.ProofOfWorkUseCase
	deleteAllPendingBlockTransactionUseCase *usecase.DeleteAllPendingBlockTransactionUseCase
}

func NewMiningService(
	config *config.Config,
	logger *slog.Logger,
	kmutex kmutexutil.KMutexProvider,
	uc1 *usecase.ListAllPendingBlockTransactionUseCase,
	uc2 *usecase.GetLastBlockDataHashUseCase,
	uc3 *usecase.SetLastBlockDataHashUseCase,
	uc4 *usecase.GetBlockDataUseCase,
	uc5 *usecase.CreateBlockDataUseCase,
	uc6 *usecase.ProofOfWorkUseCase,
	uc7 *usecase.DeleteAllPendingBlockTransactionUseCase,
) *MiningService {
	return &MiningService{config, logger, kmutex, uc1, uc2, uc3, uc4, uc5, uc6, uc7}
}

func (s *MiningService) Execute(ctx context.Context) error {
	s.logger.Debug("starting mining service...")
	defer s.logger.Debug("finished mining service")

	//
	// STEP 1:
	// Lock this function - this is important - because it will fetch all the
	// latest pending block transactions, so there needs to be a lockdown of
	// this service that when it runs it will no longer accept any more calls
	// from the system. Therefore we are using a key-based mutex to lock down
	// this service to act as a singleton runtime usage.
	//

	// Lock the mining service until it has completed executing (or errored).
	s.kmutex.Acquire("mining-service")
	defer s.kmutex.Release("mining-service")

	pendingBlockTxs, err := s.listAllPendingBlockTransactionUseCase.Execute()
	if err != nil {
		s.logger.Debug("failed listing pending block transactions",
			slog.Any("error", err))
		return nil
	}
	if len(pendingBlockTxs) <= 0 {
		s.logger.Debug("no pending block transactions for mining service")
	}

	s.logger.Debug("executing mining for pending block transactions",
		slog.Int("count", len(pendingBlockTxs)),
	)

	//TODO: IMPL:
	//-------------------------------
	// Get all pending block txs.
	// Create block data
	// Submit to blockchain network
	//      TODO: Receive purposed blockdata
	//      TODO: Verify purposed blockdata
	//      TODO: Add blockdata to blockchain
	//      TODO: Broadcast to p2p network the new blockdata.
	// Delete all pending block txs.
	//-------------------------------

	return nil
}
