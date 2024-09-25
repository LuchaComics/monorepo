package controller

import (
	"context"
	"log/slog"

	mqb "github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/adapter/messagequeuebroker"
	pt_ds "github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/app/signedtransaction/datastore"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/provider/uuid"
)

// MempoolController provides all the functionality of accepting pending
// transactions and submitting them to the miner when threshold is met.

type MempoolController interface {
	// ReadyToDistribute function blocks your execution flow until a new pending signed transaction is received from a message queue and the pending signed transaction is ready to be distributed to the network.
	ReadyToDistribute(ctx context.Context) (*pt_ds.SignedTransaction, error)

	// Receive function to be called when a pending signed transaction is received from the network.
	Receive(ctx context.Context, pendingTx *pt_ds.SignedTransaction) error
}

type mempoolControllerImpl struct {
	logger                  *slog.Logger
	uuid                    uuid.Provider
	messageQueueBroker      mqb.MessageQueueBroker
	signedTransactionStorer pt_ds.SignedTransactionStorer
}

func NewController(
	cfg *config.Config,
	logger *slog.Logger,
	uuid uuid.Provider,
	broker mqb.MessageQueueBroker,
	pt pt_ds.SignedTransactionStorer,
) MempoolController {
	return &mempoolControllerImpl{
		logger:                  logger,
		uuid:                    uuid,
		messageQueueBroker:      broker,
		signedTransactionStorer: pt,
	}
}
