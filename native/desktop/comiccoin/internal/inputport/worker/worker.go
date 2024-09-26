package worker

import (
	"context"
	"log/slog"

	mempool_c "github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/app/mempool/controller"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/inputport"
)

type workerInputPort struct {
	cfg               *config.Config
	logger            *slog.Logger
	mempoolController mempool_c.MempoolController
}

func NewInputPort(
	cfg *config.Config,
	logger *slog.Logger,
	mempool mempool_c.MempoolController,
) inputport.InputPortServer {

	port := &workerInputPort{
		cfg:               cfg,
		logger:            logger,
		mempoolController: mempool,
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
		port.mempoolController.RunSendPendingSignedTransactionsToLocalMineOperation(context.Background())
	}()
}

func (port *workerInputPort) Shutdown() {
	port.logger.Info("Gracefully shutting down background worker")
}
