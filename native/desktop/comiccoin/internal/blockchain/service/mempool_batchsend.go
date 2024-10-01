package service

import (
	"context"
	"log/slog"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/usecase"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/pkg/kmutexutil"
)

// MempoolBatchService represents a single Memory Pool node in the distributed
// / P2p blockchain network which waits to receive signed transactions and
// saves them locally to be processed by our node.
type MempoolBatchSendService struct {
	config                               *config.Config
	logger                               *slog.Logger
	kmutex                               kmutexutil.KMutexProvider
	listAllSignedTransactionDTOUseCase   *usecase.ListAllSignedTransactionUseCase
	deleteAllSignedTransactionDTOUseCase *usecase.DeleteAllSignedTransactionUseCase
}

func NewMempoolBatchSendService(
	cfg *config.Config,
	logger *slog.Logger,
	kmutex kmutexutil.KMutexProvider,
	uc1 *usecase.ListAllSignedTransactionUseCase,
	uc2 *usecase.DeleteAllSignedTransactionUseCase,
) *MempoolBatchSendService {
	return &MempoolBatchSendService{cfg, logger, kmutex, uc1, uc2}
}

func (s *MempoolBatchSendService) Execute(ctx context.Context) error {
	//
	// STEP 1:
	// List all the signed transactions in the local database.
	//

	// Lock the mempool's database so we coordinate when we delete the mempool
	// and when we add mempool.
	s.kmutex.Acquire("mempool")
	defer s.kmutex.Release("mempool")

	stxs, err := s.listAllSignedTransactionDTOUseCase.Execute()
	if err != nil {
		s.logger.Error("mempool failed listing signed transaction",
			slog.Any("error", err))
		return err
	}

	//
	// STEP 2
	//

	// DEVELOPERS NOTE:
	// 1. We can implement complex algorithms about sorting based on rewards but
	//    we are not implementing transaction fees so it doesn't matter. We
	//    just get the latest.
	// 2. To keep things simple, we will just check to see if we meet the
	//    transaction requirement per block and if we meet it then we can
	//    send the transactions to the miner.
	if len(stxs) < int(s.config.Blockchain.TransPerBlock) {
		return nil
	}

	//
	// STEP 3
	// Queue our signed transactions for the miner.
	//

	//TODO: IMPL.

	//
	// STEP 4
	// Delete all our signed transactions.
	//

	//TODO: IMPL.
	// if err := s.deleteAllSignedTransactionDTOUseCase.Execute(); err != nil {
	// 	return nil
	// }

	return nil
}
