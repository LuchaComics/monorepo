package task

import (
	"context"
	"log/slog"
	"time"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/config"
	task "github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/interface/task/handler"
	taskmnghandler "github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/interface/task/handler"
)

type TaskManager interface {
	Run()
	Shutdown()
}

type taskManagerImpl struct {
	cfg                           *config.Config
	logger                        *slog.Logger
	mempoolReceiveTaskHandler     *task.MempoolReceiveTaskHandler
	mempoolBatchSendTaskHandler   *task.MempoolBatchSendTaskHandler
	miningTaskHandler             *task.MiningTaskHandler
	validationTaskHandler         *task.ValidationTaskHandler
	consensusTaskHandler          *task.ConsensusTaskHandler
	blockDataDTOServerTaskHandler *task.BlockDataDTOServerTaskHandler
}

func NewTaskManager(
	cfg *config.Config,
	logger *slog.Logger,
	mempoolReceiveTaskHandler *task.MempoolReceiveTaskHandler,
	mempoolBatchSendTaskHandler *task.MempoolBatchSendTaskHandler,
	miningTaskHandler *task.MiningTaskHandler,
	validationTaskHandler *task.ValidationTaskHandler,
	consensusTaskHandler *task.ConsensusTaskHandler,
	blockDataDTOServerTaskHandler *task.BlockDataDTOServerTaskHandler,
) TaskManager {
	port := &taskManagerImpl{
		cfg:                           cfg,
		logger:                        logger,
		mempoolReceiveTaskHandler:     mempoolReceiveTaskHandler,
		mempoolBatchSendTaskHandler:   mempoolBatchSendTaskHandler,
		miningTaskHandler:             miningTaskHandler,
		validationTaskHandler:         validationTaskHandler,
		consensusTaskHandler:          consensusTaskHandler,
		blockDataDTOServerTaskHandler: blockDataDTOServerTaskHandler,
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

	go func(consensus *taskmnghandler.ConsensusTaskHandler, loggerp *slog.Logger) {
		ctx := context.Background()
		for {
			if err := consensus.Execute(ctx); err != nil {
				loggerp.Error("consensus error", slog.Any("error", err))
			}
			time.Sleep(5 * time.Second)
			loggerp.Error("executing consensus mechanism again")
		}
	}(port.consensusTaskHandler, port.logger)

	go func(server *taskmnghandler.BlockDataDTOServerTaskHandler, loggerp *slog.Logger) {
		ctx := context.Background()
		for {
			if err := server.Execute(ctx); err != nil {
				loggerp.Error("blockdatabto upload server error",
					slog.Any("error", err))
				time.Sleep(10 * time.Second)
				continue
			}
			time.Sleep(5 * time.Second)
			loggerp.Debug("shared local blockchain with network")
			break
		}
	}(port.blockDataDTOServerTaskHandler, port.logger)
}

func (port *taskManagerImpl) Shutdown() {
	port.logger.Info("Gracefully shutting down Task Manager")
}
