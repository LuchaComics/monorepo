package task

import (
	"context"
	"fmt"
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
	cfg                                    *config.Config
	logger                                 *slog.Logger
	mempoolReceiveTaskHandler              *task.MempoolReceiveTaskHandler
	mempoolBatchSendTaskHandler            *task.MempoolBatchSendTaskHandler
	miningTaskHandler                      *task.MiningTaskHandler
	validationTaskHandler                  *task.ValidationTaskHandler
	blockDataDTOServerTaskHandler          *task.BlockDataDTOServerTaskHandler
	majorityVoteConsensusServerTaskHandler *task.MajorityVoteConsensusServerTaskHandler
	majorityVoteConsensusClientTaskHandler *task.MajorityVoteConsensusClientTaskHandler
}

func NewTaskManager(
	cfg *config.Config,
	logger *slog.Logger,
	mempoolReceiveTaskHandler *task.MempoolReceiveTaskHandler,
	mempoolBatchSendTaskHandler *task.MempoolBatchSendTaskHandler,
	miningTaskHandler *task.MiningTaskHandler,
	validationTaskHandler *task.ValidationTaskHandler,
	blockDataDTOServerTaskHandler *task.BlockDataDTOServerTaskHandler,
	majorityVoteConsensusServerTaskHandler *task.MajorityVoteConsensusServerTaskHandler,
	majorityVoteConsensusClientTaskHandler *task.MajorityVoteConsensusClientTaskHandler,
) TaskManager {
	port := &taskManagerImpl{
		cfg:                                    cfg,
		logger:                                 logger,
		mempoolReceiveTaskHandler:              mempoolReceiveTaskHandler,
		mempoolBatchSendTaskHandler:            mempoolBatchSendTaskHandler,
		miningTaskHandler:                      miningTaskHandler,
		validationTaskHandler:                  validationTaskHandler,
		blockDataDTOServerTaskHandler:          blockDataDTOServerTaskHandler,
		majorityVoteConsensusServerTaskHandler: majorityVoteConsensusServerTaskHandler,
		majorityVoteConsensusClientTaskHandler: majorityVoteConsensusClientTaskHandler,
	}
	return port
}

func (port *taskManagerImpl) Run() {
	ctx := context.Background()
	port.logger.Info("Running Task Manager")

	go func() {
		port.logger.Info("Runningmempool (receive) service...")
		for {
			taskErr := port.mempoolReceiveTaskHandler.Execute(ctx)
			if taskErr != nil {
				port.logger.Error("failed executing mempool receive task, restarting task in 1 minute...", slog.Any("error", taskErr))
				time.Sleep(1 * time.Second)
			}
		}
	}()
	go func() {
		port.logger.Info("Running mempool (send) service...")
		for {
			taskErr := port.mempoolBatchSendTaskHandler.Execute(ctx)
			if taskErr != nil {
				port.logger.Error("failed executing mempool batch send task, restarting task in 1 minute...", slog.Any("error", taskErr))
				time.Sleep(1 * time.Minute)
			}
			time.Sleep(1 * time.Second)
		}
	}()
	go func() {
		port.logger.Info("Running mining service...")
		for {
			taskErr := port.miningTaskHandler.Execute(ctx)
			if taskErr != nil {
				port.logger.Error("failed executing mining task, restarting task in 1 minute...", slog.Any("error", taskErr))
				time.Sleep(1 * time.Minute)
			}
			time.Sleep(1 * time.Second)
		}
	}()
	go func() {
		port.logger.Info("Running validation service...")
		for {
			taskErr := port.validationTaskHandler.Execute(ctx)
			if taskErr != nil {
				port.logger.Error("failed executing validation task, restarting task in 1 minute...", slog.Any("error", taskErr))
				time.Sleep(1 * time.Minute)
			}
			time.Sleep(1 * time.Second)
		}
	}()

	go func(consensus *taskmnghandler.MajorityVoteConsensusServerTaskHandler, loggerp *slog.Logger) {
		port.logger.Info("Running consensus server...")
		ctx := context.Background()
		for {
			if err := consensus.Execute(ctx); err != nil {
				loggerp.Error("consensus error", slog.Any("error", err))
			}
			// DEVELOPERS NOTE:
			// No need for delays, automatically start executing again.
			loggerp.Debug("blockchain consensus serving done, excuting again ...")
		}
	}(port.majorityVoteConsensusServerTaskHandler, port.logger)

	go func(consensus *taskmnghandler.MajorityVoteConsensusClientTaskHandler, loggerp *slog.Logger) {
		loggerp.Info("Running consensus client...")
		ctx := context.Background()
		for {
			if err := consensus.Execute(ctx); err != nil {
				loggerp.Error("consensus error", slog.Any("error", err))
			}
			time.Sleep(time.Duration(port.cfg.Blockchain.ConsensusPollingDelayInMinutes) * time.Minute)
			loggerp.Debug(fmt.Sprintf("blockchain consensus client ran, will run again in %v minutes...", port.cfg.Blockchain.ConsensusPollingDelayInMinutes))
		}
	}(port.majorityVoteConsensusClientTaskHandler, port.logger)

	go func(server *taskmnghandler.BlockDataDTOServerTaskHandler, loggerp *slog.Logger) {
		loggerp.Info("Running block data dto server...")
		ctx := context.Background()
		for {
			if err := server.Execute(ctx); err != nil {
				loggerp.Error("blockdatabto upload server error",
					slog.Any("error", err))
				time.Sleep(10 * time.Second)
				continue
			}
			// DEVELOPERS NOTE:
			// No need for delays, automatically start executing again.
			port.logger.Debug("block data dto server executing again ...")
		}
	}(port.blockDataDTOServerTaskHandler, port.logger)
}

func (port *taskManagerImpl) Shutdown() {
	port.logger.Info("Gracefully shutting down Task Manager")
}
