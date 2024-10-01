package service

import (
	"context"
	"log/slog"
	"time"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/domain"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/usecase"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/pkg/kmutexutil"
)

// MempoolBatchService represents a single Memory Pool node in the distributed
// / P2p blockchain network which waits to receive signed transactions and
// saves them locally to be processed by our node.
type MempoolBatchSendService struct {
	config                                *config.Config
	logger                                *slog.Logger
	kmutex                                kmutexutil.KMutexProvider
	listAllMempoolTransactionDTOUseCase   *usecase.ListAllMempoolTransactionUseCase
	createPendingBlockTransactionUseCase  *usecase.CreatePendingBlockTransactionUseCase
	deleteAllMempoolTransactionDTOUseCase *usecase.DeleteAllMempoolTransactionUseCase
}

func NewMempoolBatchSendService(
	cfg *config.Config,
	logger *slog.Logger,
	kmutex kmutexutil.KMutexProvider,
	uc1 *usecase.ListAllMempoolTransactionUseCase,
	uc2 *usecase.CreatePendingBlockTransactionUseCase,
	uc3 *usecase.DeleteAllMempoolTransactionUseCase,
) *MempoolBatchSendService {
	return &MempoolBatchSendService{cfg, logger, kmutex, uc1, uc2, uc3}
}

func (s *MempoolBatchSendService) Execute(ctx context.Context) error {
	//
	// STEP 1:
	// List all the signed transactions in the local database.
	//

	// Lock the mempool's database so we coordinate when we delete the mempool
	// and when we add mempool.
	s.kmutex.Acquire("mempool-service")
	defer s.kmutex.Release("mempool-service")

	stxs, err := s.listAllMempoolTransactionDTOUseCase.Execute()
	if err != nil {
		s.logger.Error("mempool failed listing mempool transaction",
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
		// Do nothing, just return this function with nothing.
		s.logger.Debug("mempool skipped submitted to minging service")
		return nil
	}

	//
	// STEP 3
	// Queue our signed transactions for the miner.
	//

	s.logger.Debug("mempool submitting to minging service",
		slog.Any("txs_count", len(stxs)),
	)

	for _, stx := range stxs {
		pendingBlockTx := &domain.PendingBlockTransaction{
			MempoolTransaction: *stx,
			TimeStamp:          uint64(time.Now().Unix()),      // Ethereum: The time the transaction was received.
			GasPrice:           s.config.Blockchain.GasPrice,   // Ethereum: The price of one unit of gas to be paid for fees.
			GasUnits:           s.config.Blockchain.UnitsOfGas, // Ethereum: The number of units of gas used for this transaction.
		}
		if createErr := s.createPendingBlockTransactionUseCase.Execute(pendingBlockTx); err != nil {
			s.logger.Error("mempool failed creating pending signed transaction",
				slog.Any("error", createErr))
			return createErr
		}
		s.logger.Debug("mempool submitted to minging service",
			slog.Any("tx_nonce", stx.Nonce),
		)
	}

	//
	// STEP 4
	// Delete all our signed transactions.
	//

	if deleteAllErr := s.deleteAllMempoolTransactionDTOUseCase.Execute(); err != nil {
		s.logger.Error("mempool failed deleting all mempool transaction",
			slog.Any("error", deleteAllErr))
		return nil
	}

	return nil
}
