package worker

import (
	"context"
	"log"
	"log/slog"

	pubsub "github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/adapter/pubsub"
	blockchain_c "github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/app/blockchain/controller"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/config"
)

// BlockchainWorker provides all the functionality that can be performed
// on the `ComicCoin` cryptocurrency.

type BlockchainWorker interface {
	RunMinerOperation(ctx context.Context)
}

type blockchainWorkerImpl struct {
	config            *config.Config
	logger            *slog.Logger
	localPubSubBroker pubsub.PubSubBroker
	p2pPubSubBroker   pubsub.PubSubBroker
	Controller        blockchain_c.BlockchainController
}

func NewWorker(
	cfg *config.Config,
	logger *slog.Logger,
	locpsbroker pubsub.PubSubBroker,
	p2psbroker pubsub.PubSubBroker,
	c blockchain_c.BlockchainController,
) BlockchainWorker {
	// Defensive code to protect the programmer from any errors.
	if cfg.Blockchain.Difficulty <= 0 {
		log.Fatal("cannot have blochain difficulty less then or equal to zero")
	}
	if cfg.Blockchain.TransPerBlock <= 0 {
		log.Fatal("cannot have blochain transactions per block less then or equal to zero")
	}

	impl := &blockchainWorkerImpl{
		config:            cfg,
		logger:            logger,
		localPubSubBroker: locpsbroker,
		p2pPubSubBroker:   p2psbroker,
		Controller:        c,
	}
	return impl
}
