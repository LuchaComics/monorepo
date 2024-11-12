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
	dataChan chan *domain.MempoolTransaction
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
	select {
	case newEntry, ok := <-uc.dataChan:
		if !ok {
			return nil, errors.New("channel closed")
		}
		return newEntry, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

// Terminate releases the channel resource
func (uc *MempoolTransactionInsertionDetectorUseCase) Terminate() {
	close(uc.quitChan)
	close(uc.dataChan)
}
