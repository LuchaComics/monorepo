package task

import (
	"context"
	"log/slog"

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
}

func NewTaskManager(
	cfg *config.Configuration,
	logger *slog.Logger,
	t1 *taskhandler.AttachmentGarbageCollectorTaskHandler,

) TaskManager {
	port := &taskManagerImpl{
		cfg:                                   cfg,
		logger:                                logger,
		attachmentGarbageCollectorTaskHandler: t1,
	}
	return port
}

func (port *taskManagerImpl) Run() {
	backgroundCtx := context.Background()
	port.logger.Info("Running Task Manager")

	go func(task *taskhandler.AttachmentGarbageCollectorTaskHandler, loggerp *slog.Logger) {
		loggerp.Info("Starting attachment garbage collection...")

		for {
			if err := task.Execute(backgroundCtx); err != nil {
				loggerp.Error("Failed executing attachment garbage collection",
					slog.Any("error", err))
			}
			// DEVELOPERS NOTE:
			// No need for delays, automatically start executing again.
			port.logger.Debug("Attachment garbage collection will run again ...")
		}
	}(port.attachmentGarbageCollectorTaskHandler, port.logger)
}

func (port *taskManagerImpl) Shutdown() {
	port.logger.Info("Gracefully shutting down Task Manager")
}
