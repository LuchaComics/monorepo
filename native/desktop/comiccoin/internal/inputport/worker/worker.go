package worker

import (
	"context"
	"log/slog"

	blockchain_c "github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/app/blockchain/controller"
	mempool_c "github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/app/mempool/controller"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/inputport"
)

type workerInputPort struct {
	cfg                  *config.Config
	logger               *slog.Logger
	mempoolController    mempool_c.MempoolController
	blockchainController blockchain_c.BlockchainController
}

func NewInputPort(
	cfg *config.Config,
	logger *slog.Logger,
	mempool mempool_c.MempoolController,
	blockchain blockchain_c.BlockchainController,
) inputport.InputPortServer {

	port := &workerInputPort{
		cfg:                  cfg,
		logger:               logger,
		mempoolController:    mempool,
		blockchainController: blockchain,
	}

	return port
}

func (port *workerInputPort) Run() {
	// ctx := context.Background()
	port.logger.Info("Running background worker")
	go func() {
		port.mempoolController.RunReceiveFromNetworkOperation(context.Background())
	}()
	go func() {
		port.blockchainController.RunMinerOperationInBackground(context.Background())
	}()
}

func (port *workerInputPort) Shutdown() {
	port.logger.Info("Gracefully shutting down background worker")
}
