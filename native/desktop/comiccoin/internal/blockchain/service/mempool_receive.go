package service

import (
	"context"
	"log/slog"
	"time"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/usecase"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/pkg/kmutexutil"
)

// MempoolReceiveService represents a single Memory Pool node in the distributed
// / P2p blockchain network which waits to receive signed transactions and
// saves them locally to be processed by our node.
type MempoolReceiveService struct {
	config                              *config.Config
	logger                              *slog.Logger
	kmutex                              kmutexutil.KMutexProvider
	receiveMempoolTransactionDTOUseCase *usecase.ReceiveMempoolTransactionDTOUseCase
	createMempoolTransactionUseCase     *usecase.CreateMempoolTransactionUseCase
}

func NewMempoolReceiveService(
	cfg *config.Config,
	logger *slog.Logger,
	kmutex kmutexutil.KMutexProvider,
	uc1 *usecase.ReceiveMempoolTransactionDTOUseCase,
	uc2 *usecase.CreateMempoolTransactionUseCase,
) *MempoolReceiveService {
	return &MempoolReceiveService{cfg, logger, kmutex, uc1, uc2}
}

func (s *MempoolReceiveService) Execute(ctx context.Context) error {

	//
	// STEP 1
	// Wait to receive data (which also was validated) from the P2P network.
	//

	stx, err := s.receiveMempoolTransactionDTOUseCase.Execute(ctx)
	if err != nil {
		s.logger.Error("mempool failed receiving dto",
			slog.Any("error", err))
		return err
	}
	if stx == nil {
		// Developer Note:
		// If we haven't received anything, that means we haven't connected to
		// the distributed / P2P network, so all we can do at the moment is to
		// pause the execution for 1 second and then retry again.
		time.Sleep(1 * time.Second)
		return nil
	}

	s.logger.Info("received dto from network",
		slog.Any("tx_nonce", stx.Nonce),
	)

	//
	// STEP 2:
	// Save to our local database.
	//

	// Lock the mempool's database so we coordinate when we delete the mempool
	// and when we add mempool.
	s.kmutex.Acquire("mempool-service")
	defer s.kmutex.Release("mempool-service")

	if err := s.createMempoolTransactionUseCase.Execute(stx); err != nil {
		s.logger.Error("mempool failed saving mempool transaction",
			slog.Any("error", err))
		return err
	}

	s.logger.Info("saved to mempool",
		slog.Any("tx_nonce", stx.Nonce),
	)

	return nil
}
