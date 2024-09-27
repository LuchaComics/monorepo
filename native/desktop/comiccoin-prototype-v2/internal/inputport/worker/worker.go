package worker

import (
	"context"
	"log/slog"

	blockchain_work "github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/app/blockchain/worker"
	mempool_work "github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/app/mempool/worker"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/inputport"
)

type workerInputPort struct {
	cfg              *config.Config
	logger           *slog.Logger
	mempoolWorker    mempool_work.MempoolWorker
	blockchainWorker blockchain_work.BlockchainWorker
}

func NewInputPort(
	cfg *config.Config,
	logger *slog.Logger,
	mempool mempool_work.MempoolWorker,
	blockchain blockchain_work.BlockchainWorker,
) inputport.InputPortServer {

	port := &workerInputPort{
		cfg:              cfg,
		logger:           logger,
		mempoolWorker:    mempool,
		blockchainWorker: blockchain,
	}

	return port
}

func (port *workerInputPort) Run() {
	// ctx := context.Background()
	port.logger.Info("Running background worker")
	go func() {
		port.mempoolWorker.RunReceiveFromNetworkOperation(context.Background())
	}()
	go func() {
		port.blockchainWorker.RunMinerOperation(context.Background())
	}()
}

func (port *workerInputPort) Shutdown() {
	port.logger.Info("Gracefully shutting down background worker")
}
