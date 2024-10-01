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
	proofOfWorkUseCase                      *usecase.ProofOfWorkUseCase
	deleteAllPendingBlockTransactionUseCase *usecase.DeleteAllPendingBlockTransactionUseCase
}

func NewMiningService(config *config.Config, logger *slog.Logger, kmutex kmutexutil.KMutexProvider, uc1 *usecase.ListAllPendingBlockTransactionUseCase, uc2 *usecase.ProofOfWorkUseCase, uc3 *usecase.DeleteAllPendingBlockTransactionUseCase) *MiningService {
	return &MiningService{config, logger, kmutex, uc1, uc2, uc3}
}

func (s *MiningService) Execute(ctx context.Context) error {
	s.logger.Debug("starting mining service...")
	defer s.logger.Debug("finished mining service")

	// Lock the mining service until it has completed executing (or errored).
	s.kmutex.Acquire("mining-service")
	defer s.kmutex.Release("mining-service")

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
