package controller

import (
	"context"
	"log/slog"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/adapter/pubsub"
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
	localPubSubBroker       pubsub.PubSubBroker
	p2pPubSubBroker         pubsub.PubSubBroker
	signedTransactionStorer pt_ds.SignedTransactionStorer
}

func NewController(
	cfg *config.Config,
	logger *slog.Logger,
	uuid uuid.Provider,
	locpsbroker pubsub.PubSubBroker,
	p2psbroker pubsub.PubSubBroker,
	pt pt_ds.SignedTransactionStorer,
) MempoolController {
	mempool := &mempoolControllerImpl{
		config:                  cfg,
		logger:                  logger,
		uuid:                    uuid,
		localPubSubBroker:       locpsbroker,
		p2pPubSubBroker:         p2psbroker,
		signedTransactionStorer: pt,
	}
	return mempool
}
