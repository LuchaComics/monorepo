package controller

import (
	"context"
	"log/slog"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/adapter/pubsub"
	mempool_c "github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/app/mempool/controller"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/config"
)

// MempoolWorker provides all the functionality of accepting pending
// transactions and submitting them to the miner when threshold is met.

type MempoolWorker interface {
	// RunReceiveFromNetworkOperation function to be called when a pending signed transaction is received from the network.
	RunReceiveFromNetworkOperation(ctx context.Context) error
}

type mempoolWorkerImpl struct {
	config            *config.Config
	logger            *slog.Logger
	localPubSubBroker pubsub.PubSubBroker
	p2pPubSubBroker   pubsub.PubSubBroker
	mempoolController mempool_c.MempoolController
}

func NewWorker(
	cfg *config.Config,
	logger *slog.Logger,
	locpsbroker pubsub.PubSubBroker,
	p2psbroker pubsub.PubSubBroker,
	mempoolC mempool_c.MempoolController,
) MempoolWorker {
	mempool := &mempoolWorkerImpl{
		config:            cfg,
		logger:            logger,
		localPubSubBroker: locpsbroker,
		p2pPubSubBroker:   p2psbroker,
		mempoolController: mempoolC,
	}
	return mempool
}
