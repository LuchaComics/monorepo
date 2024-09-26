package controller

import (
	"context"
	"log/slog"

	dpubsub "github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/adapter/distributedpubsub"
	pt_ds "github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/app/signedtransaction/datastore"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/provider/uuid"
)

// MempoolController provides all the functionality of accepting pending
// transactions and submitting them to the miner when threshold is met.

type MempoolController interface {
	// RunReceiveFromNetworkOperation function to be called when a pending signed transaction is received from the network.
	RunReceiveFromNetworkOperation(ctx context.Context) error
}

type mempoolControllerImpl struct {
	config                  *config.Config
	logger                  *slog.Logger
	uuid                    uuid.Provider
	pubSubBroker            dpubsub.PublishSubscribeBroker
	signedTransactionStorer pt_ds.SignedTransactionStorer
}

func NewController(
	cfg *config.Config,
	logger *slog.Logger,
	uuid uuid.Provider,
	psbroker dpubsub.PublishSubscribeBroker,
	pt pt_ds.SignedTransactionStorer,
) MempoolController {
	mempool := &mempoolControllerImpl{
		config:                  cfg,
		logger:                  logger,
		uuid:                    uuid,
		pubSubBroker:            psbroker,
		signedTransactionStorer: pt,
	}
	return mempool
}
