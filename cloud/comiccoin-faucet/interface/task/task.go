package task

import (
	"context"
	"log/slog"
	"time"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/config"
	taskhandler "github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/interface/task/handler"
)

type TaskManager interface {
	Run()
	Shutdown()
}

type taskManagerImpl struct {
	cfg                                              *config.Configuration
	logger                                           *slog.Logger
	attachmentGarbageCollectorTaskHandler            *taskhandler.AttachmentGarbageCollectorTaskHandler
	blockchainSyncWithBlockchainAuthorityTaskHandler *taskhandler.BlockchainSyncWithBlockchainAuthorityTaskHandler
}

func NewTaskManager(
	cfg *config.Configuration,
	logger *slog.Logger,
	t1 *taskhandler.AttachmentGarbageCollectorTaskHandler,
	t2 *taskhandler.BlockchainSyncWithBlockchainAuthorityTaskHandler,

) TaskManager {
	port := &taskManagerImpl{
		cfg:                                   cfg,
		logger:                                logger,
		attachmentGarbageCollectorTaskHandler: t1,
		blockchainSyncWithBlockchainAuthorityTaskHandler: t2,
	}
	return port
}

func (port *taskManagerImpl) Run() {
	port.logger.Info("Running Task Manager")

	go func(task *taskhandler.AttachmentGarbageCollectorTaskHandler, loggerp *slog.Logger) {
		loggerp.Info("Starting attachment garbage collector...")

		for {
			if err := task.Execute(context.Background()); err != nil {
				loggerp.Error("Failed executing attachment garbage collector",
					slog.Any("error", err))
			}
			// port.logger.Debug("Attachment garbage collector will run again in 15 seconds...")
			time.Sleep(15 * time.Second)
		}
	}(port.attachmentGarbageCollectorTaskHandler, port.logger)

	go func(task *taskhandler.BlockchainSyncWithBlockchainAuthorityTaskHandler, loggerp *slog.Logger) {
		loggerp.Info("Starting blockchain sync with the Authority...")

		for {
			if err := task.Execute(context.Background()); err != nil {
				loggerp.Error("Failed executing blockchain sync with the Authority.",
					slog.Any("error", err))
			}
			port.logger.Debug("Blockchain sync with the Authority will rerun again in 15 seconds...")
			time.Sleep(15 * time.Second)
		}
	}(port.blockchainSyncWithBlockchainAuthorityTaskHandler, port.logger)
}

func (port *taskManagerImpl) Shutdown() {
	port.logger.Info("Gracefully shutting down Task Manager")
}
