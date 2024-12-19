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
	cfg                                   *config.Configuration
	logger                                *slog.Logger
	attachmentGarbageCollectorTaskHandler *taskhandler.AttachmentGarbageCollectorTaskHandler
	blockchainSyncManagerTaskHandler      *taskhandler.BlockchainSyncManagerTaskHandler
}

func NewTaskManager(
	cfg *config.Configuration,
	logger *slog.Logger,
	t1 *taskhandler.AttachmentGarbageCollectorTaskHandler,
	t2 *taskhandler.BlockchainSyncManagerTaskHandler,

) TaskManager {
	port := &taskManagerImpl{
		cfg:                                   cfg,
		logger:                                logger,
		attachmentGarbageCollectorTaskHandler: t1,
		blockchainSyncManagerTaskHandler:      t2,
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
			// DEVELOPERS NOTE:
			// No need for delays, automatically start executing again.
			port.logger.Debug("Attachment garbage collector will run again in 15 seconds...")
			time.Sleep(15 * time.Second)
		}
	}(port.attachmentGarbageCollectorTaskHandler, port.logger)

	go func(task *taskhandler.BlockchainSyncManagerTaskHandler, loggerp *slog.Logger) {
		loggerp.Info("Starting blockchain manager...")

		for {
			if err := task.Execute(context.Background()); err != nil {
				loggerp.Error("Failed executing blockchain manager",
					slog.Any("error", err))
			}
			// DEVELOPERS NOTE:
			// No need for delays, automatically start executing again.
			port.logger.Debug("Blockchain manager will run again ...")
		}
	}(port.blockchainSyncManagerTaskHandler, port.logger)
}

func (port *taskManagerImpl) Shutdown() {
	port.logger.Info("Gracefully shutting down Task Manager")
}
