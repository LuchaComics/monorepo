package usecase

import (
	"context"
	"errors"
	"log"
	"log/slog"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/config"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/domain"
)

type MempoolTransactionInsertionDetectorUseCase struct {
	config   *config.Configuration
	logger   *slog.Logger
	repo     domain.MempoolTransactionRepository
	dataChan <-chan domain.MempoolTransaction
	quitChan chan struct{}
}

func NewMempoolTransactionInsertionDetectorUseCase(config *config.Configuration, logger *slog.Logger, repo domain.MempoolTransactionRepository) *MempoolTransactionInsertionDetectorUseCase {
	dataChan, quitChan, err := repo.GetInsertionChangeStreamChannel(context.Background())
	if err != nil {
		log.Fatalf("NewMempoolTransactionInsertionDetectorUseCase: Failed initializing use-case: %v\n", err)
	}
	return &MempoolTransactionInsertionDetectorUseCase{config, logger, repo, dataChan, quitChan}
}

func (uc *MempoolTransactionInsertionDetectorUseCase) Execute(ctx context.Context) (*domain.MempoolTransaction, error) {
	// uc.logger.Debug("Waiting to receive...")
	select {
	case mempoolTx, ok := <-uc.dataChan:
		if !ok {
			// uc.logger.Debug("Error receiving.")
			return nil, errors.New("channel closed")
		}
		// uc.logger.Debug("Data received.",
		// 	slog.Any("dataChan", newEntry))

		if err := mempoolTx.Validate(uc.config.Blockchain.ChainID, true); err != nil {
			uc.logger.Warn("Validation failed for get",
				slog.Any("error", err),
				slog.Any("Transaction", mempoolTx.Transaction),
				slog.Any("tx_sig_v", mempoolTx.V),
				slog.Any("tx_sig_r", mempoolTx.R),
				slog.Any("tx_sig_s", mempoolTx.S))
			return nil, err
		}

		return &mempoolTx, nil
	case <-ctx.Done():
		// uc.logger.Debug("Done receiving.")
		return nil, ctx.Err()
	}
}

// Terminate releases the channel resource
func (uc *MempoolTransactionInsertionDetectorUseCase) Terminate() {
	uc.logger.Debug("Closing change stream connection...")
	close(uc.quitChan)
}
