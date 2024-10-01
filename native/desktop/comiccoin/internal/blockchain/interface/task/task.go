package task

import (
	"context"
	"log/slog"
	"time"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/config"
	task "github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/interface/task/handler"
)

type TaskManager interface {
	Run()
	Shutdown()
}

type taskManagerImpl struct {
	cfg                         *config.Config
	logger                      *slog.Logger
	mempoolReceiveTaskHandler   *task.MempoolReceiveTaskHandler
	mempoolBatchSendTaskHandler *task.MempoolBatchSendTaskHandler
	miningTaskHandler           *task.MiningTaskHandler
	validationTaskHandler       *task.ValidationTaskHandler
}

func NewTaskManager(
	cfg *config.Config,
	logger *slog.Logger,
	mempoolReceiveTaskHandler *task.MempoolReceiveTaskHandler,
	mempoolBatchSendTaskHandler *task.MempoolBatchSendTaskHandler,
	miningTaskHandler *task.MiningTaskHandler,
	validationTaskHandler *task.ValidationTaskHandler,
) TaskManager {
	port := &taskManagerImpl{
		cfg:                         cfg,
		logger:                      logger,
		mempoolReceiveTaskHandler:   mempoolReceiveTaskHandler,
		mempoolBatchSendTaskHandler: mempoolBatchSendTaskHandler,
		miningTaskHandler:           miningTaskHandler,
		validationTaskHandler:       validationTaskHandler,
	}
	return port
}

func (port *taskManagerImpl) Run() {
	ctx := context.Background()
	port.logger.Info("Running Task Manager")

	go func() {
		for {
			taskErr := port.mempoolReceiveTaskHandler.Execute(ctx)
			if taskErr != nil {
				port.logger.Error("failed executing mempool receive task, restarting task in 1 minute...", slog.Any("error", taskErr))
				time.Sleep(1 * time.Minute)
			}
		}
	}()
	go func() {
		for {
			taskErr := port.mempoolBatchSendTaskHandler.Execute(ctx)
			if taskErr != nil {
				port.logger.Error("failed executing mempool batch send task, restarting task in 1 minute...", slog.Any("error", taskErr))
				time.Sleep(1 * time.Minute)
			}
			time.Sleep(1 * time.Minute)
		}
	}()
	go func() {
		for {
			taskErr := port.miningTaskHandler.Execute(ctx)
			if taskErr != nil {
				port.logger.Error("failed executing mining task, restarting task in 1 minute...", slog.Any("error", taskErr))
				time.Sleep(1 * time.Minute)
			}
			time.Sleep(1 * time.Minute)
		}
	}()
	go func() {
		for {
			taskErr := port.validationTaskHandler.Execute(ctx)
			if taskErr != nil {
				port.logger.Error("failed executing validation task, restarting task in 1 minute...", slog.Any("error", taskErr))
				time.Sleep(1 * time.Minute)
			}
			time.Sleep(1 * time.Minute)
		}
	}()
}

func (port *taskManagerImpl) Shutdown() {
	port.logger.Info("Gracefully shutting down Task Manager")
}
