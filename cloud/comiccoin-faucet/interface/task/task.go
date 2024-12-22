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
	cfg                                                                 *config.Configuration
	logger                                                              *slog.Logger
	attachmentGarbageCollectorTaskHandler                               *taskhandler.AttachmentGarbageCollectorTaskHandler
	blockchainSyncWithBlockchainAuthorityTaskHandler                    *taskhandler.BlockchainSyncWithBlockchainAuthorityTaskHandler
	blockchainSyncWithBlockchainAuthorityViaServerSentEventsTaskHandler *taskhandler.BlockchainSyncWithBlockchainAuthorityViaServerSentEventsTaskHandler
}

func NewTaskManager(
	cfg *config.Configuration,
	logger *slog.Logger,
	t1 *taskhandler.AttachmentGarbageCollectorTaskHandler,
	t2 *taskhandler.BlockchainSyncWithBlockchainAuthorityTaskHandler,
	t3 *taskhandler.BlockchainSyncWithBlockchainAuthorityViaServerSentEventsTaskHandler,

) TaskManager {
	port := &taskManagerImpl{
		cfg:                                   cfg,
		logger:                                logger,
		attachmentGarbageCollectorTaskHandler: t1,
		blockchainSyncWithBlockchainAuthorityTaskHandler:                    t2,
		blockchainSyncWithBlockchainAuthorityViaServerSentEventsTaskHandler: t3,
	}
	return port
}

func (port *taskManagerImpl) Run() {
	port.logger.Info("Running Task Manager")

	// DEVELOPERS NOTE:
	// On startup of the Task Manager, we want to immediately sync with
	// the Global Blockchain Network to make sure we download the latest
	// data in case we are behind. After the successful one-time sync then
	// the Task Manager will load up another task to continously run in the
	// background and sync with the Global Blockchain Network.
	port.logger.Info("Running one-time blockchain sync")
	if err := port.blockchainSyncWithBlockchainAuthorityTaskHandler.Execute(context.Background()); err != nil {
		port.logger.Error("Failed executing blockchain sync with the Authority.",
			slog.Any("error", err))
	}
	port.logger.Info("Finished running one-time blockchain sync ")

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

	//------------------
	// DEPRECATED CODE:
	//------------------
	// go func(task *taskhandler.BlockchainSyncWithBlockchainAuthorityTaskHandler, loggerp *slog.Logger) {
	// 	loggerp.Info("Starting blockchain sync with the Authority...")
	//
	// 	for {
	// 		if err := task.Execute(context.Background()); err != nil {
	// 			loggerp.Error("Failed executing blockchain sync with the Authority.",
	// 				slog.Any("error", err))
	// 		}
	// 		// port.logger.Debug("Blockchain sync with the Authority will rerun again in 15 seconds...")
	// 		time.Sleep(15 * time.Second)
	// 	}
	// }(port.blockchainSyncWithBlockchainAuthorityTaskHandler, port.logger)

	go func(task *taskhandler.BlockchainSyncWithBlockchainAuthorityViaServerSentEventsTaskHandler, loggerp *slog.Logger) {
		loggerp.Info("Starting blockchain sync with the Authority...")

		for {
			if err := task.Execute(context.Background()); err != nil {
				loggerp.Error("Failed executing blockchain sync with the Authority via SSE.",
					slog.Any("error", err))
			}
			port.logger.Debug("Blockchain sync with the Authority will rerun again in 10 seconds...")
			time.Sleep(10 * time.Second)
		}
	}(port.blockchainSyncWithBlockchainAuthorityViaServerSentEventsTaskHandler, port.logger)
}

func (port *taskManagerImpl) Shutdown() {
	port.logger.Info("Gracefully shutting down Task Manager")
}
