package service

import (
	"context"
	"log/slog"
	"time"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/usecase"
)

// MempoolReceiveService represents a single Memory Pool node in the distributed
// / P2p blockchain network which waits to receive signed transactions and
// saves them locally to be processed by our node.
type MempoolReceiveService struct {
	config                             *config.Config
	logger                             *slog.Logger
	receiveSignedTransactionDTOUseCase *usecase.ReceiveSignedTransactionDTOUseCase
}

func NewMempoolReceiveService(
	cfg *config.Config,
	logger *slog.Logger,
	uc1 *usecase.ReceiveSignedTransactionDTOUseCase,
) *MempoolReceiveService {
	return &MempoolReceiveService{cfg, logger, uc1}
}

func (s *MempoolReceiveService) Execute(ctx context.Context) error {

	//
	// STEP 1
	// Wait to receive data (which also was validated) from the P2P network.
	//

	stx, err := s.receiveSignedTransactionDTOUseCase.Execute(ctx)
	if err != nil {
		s.logger.Error("mempool failed receiving",
			slog.Any("error", err))
		return err
	}
	if stx == nil {
		s.logger.Warn("mempool did not receive anything")
		time.Sleep(1 * time.Second)
		return nil
	}

	s.logger.Debug("received from network",
		slog.Any("nonce", stx.Nonce),
	)

	//
	// STEP 2:
	// Save to our local database.
	//

	return nil
}
