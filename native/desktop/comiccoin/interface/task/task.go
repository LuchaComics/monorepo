package task

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/config/constants"
	task "github.com/LuchaComics/monorepo/native/desktop/comiccoin/interface/task/handler"
	taskmnghandler "github.com/LuchaComics/monorepo/native/desktop/comiccoin/interface/task/handler"
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
	proofOfWorkMiningTaskHandler           *task.ProofOfWorkMiningTaskHandler
	proofOfAuthorityMiningTaskHandler      *task.ProofOfAuthorityMiningTaskHandler
	proofOfWorkValidationTaskHandler       *task.ProofOfWorkValidationTaskHandler
	proofOfAuthorityValidationTaskHandler  *task.ProofOfAuthorityValidationTaskHandler
	blockDataDTOServerTaskHandler          *task.BlockDataDTOServerTaskHandler
	majorityVoteConsensusServerTaskHandler *task.MajorityVoteConsensusServerTaskHandler
	majorityVoteConsensusClientTaskHandler *task.MajorityVoteConsensusClientTaskHandler
	issuedTokenClientServiceTaskHandler    *task.IssuedTokenClientServiceTaskHandler
}

func NewTaskManager(
	cfg *config.Config,
	logger *slog.Logger,
	mempoolReceiveTaskHandler *task.MempoolReceiveTaskHandler,
	mempoolBatchSendTaskHandler *task.MempoolBatchSendTaskHandler,
	proofOfWorkMiningTaskHandler *task.ProofOfWorkMiningTaskHandler,
	proofOfAuthorityMiningTaskHandler *task.ProofOfAuthorityMiningTaskHandler,
	proofOfWorkValidationTaskHandler *task.ProofOfWorkValidationTaskHandler,
	proofOfAuthorityValidationTaskHandler *task.ProofOfAuthorityValidationTaskHandler,
	blockDataDTOServerTaskHandler *task.BlockDataDTOServerTaskHandler,
	majorityVoteConsensusServerTaskHandler *task.MajorityVoteConsensusServerTaskHandler,
	majorityVoteConsensusClientTaskHandler *task.MajorityVoteConsensusClientTaskHandler,
	issuedTokenClientServiceTaskHandler *task.IssuedTokenClientServiceTaskHandler,
) TaskManager {
	port := &taskManagerImpl{
		cfg:                                    cfg,
		logger:                                 logger,
		mempoolReceiveTaskHandler:              mempoolReceiveTaskHandler,
		mempoolBatchSendTaskHandler:            mempoolBatchSendTaskHandler,
		proofOfWorkMiningTaskHandler:           proofOfWorkMiningTaskHandler,
		proofOfAuthorityMiningTaskHandler:      proofOfAuthorityMiningTaskHandler,
		proofOfWorkValidationTaskHandler:       proofOfWorkValidationTaskHandler,
		proofOfAuthorityValidationTaskHandler:  proofOfAuthorityValidationTaskHandler,
		blockDataDTOServerTaskHandler:          blockDataDTOServerTaskHandler,
		majorityVoteConsensusServerTaskHandler: majorityVoteConsensusServerTaskHandler,
		majorityVoteConsensusClientTaskHandler: majorityVoteConsensusClientTaskHandler,
		issuedTokenClientServiceTaskHandler:    issuedTokenClientServiceTaskHandler,
	}
	return port
}

func (port *taskManagerImpl) Run() {
	ctx := context.Background()
	port.logger.Info("Running Task Manager")

	go func() {
		port.logger.Info("Running mempool (receive) service...")
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
		if port.cfg.Blockchain.EnableMiner {
			if port.cfg.Blockchain.ConsensusProtocol == constants.ConsensusPoW {
				port.logger.Info("Running PoW mining service...")
				for {
					taskErr := port.proofOfWorkMiningTaskHandler.Execute(ctx)
					if taskErr != nil {
						port.logger.Error("failed executing mining task, restarting task in 1 minute...", slog.Any("error", taskErr))
						time.Sleep(1 * time.Minute)
					}
					time.Sleep(1 * time.Second)
				}
			} else if port.cfg.Blockchain.ConsensusProtocol == constants.ConsensusPoA {
				port.logger.Info("Running PoA mining service...")
				for {
					taskErr := port.proofOfAuthorityMiningTaskHandler.Execute(ctx)
					if taskErr != nil {
						port.logger.Error("failed executing mining task, restarting task in 1 minute...", slog.Any("error", taskErr))
						time.Sleep(1 * time.Minute)
					}
					time.Sleep(1 * time.Second)
				}
			} else {
				port.logger.Info("Skipped running the mining service...")
			}
		}
	}()

	go func() {
		if port.cfg.Blockchain.ConsensusProtocol == constants.ConsensusPoW {
			port.logger.Info("Running PoW validation service...")
			for {
				taskErr := port.proofOfWorkValidationTaskHandler.Execute(ctx)
				if taskErr != nil {
					port.logger.Error("failed executing validation task, restarting task in 1 minute...", slog.Any("error", taskErr))
					time.Sleep(1 * time.Minute)
				}
				time.Sleep(1 * time.Second)
			}
		} else if port.cfg.Blockchain.ConsensusProtocol == constants.ConsensusPoA {
			port.logger.Info("Running PoA validation service...")
			for {
				taskErr := port.proofOfAuthorityValidationTaskHandler.Execute(ctx)
				if taskErr != nil {
					port.logger.Error("failed executing validation task, restarting task in 1 minute...", slog.Any("error", taskErr))
					time.Sleep(1 * time.Minute)
				}
				time.Sleep(1 * time.Second)
			}
		} else {
			port.logger.Info("Skipped running the mining service...")
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

	go func(server *taskmnghandler.IssuedTokenClientServiceTaskHandler, loggerp *slog.Logger) {
		loggerp.Info("Running issued token dto server...")
		ctx := context.Background()
		for {
			if err := server.Execute(ctx); err != nil {
				loggerp.Error("issued token server error",
					slog.Any("error", err))
				time.Sleep(10 * time.Second)
				continue
			}
			// DEVELOPERS NOTE:
			// No need for delays, automatically start executing again.
			port.logger.Debug("issued token dto server executing again ...")
		}
	}(port.issuedTokenClientServiceTaskHandler, port.logger)
}

func (port *taskManagerImpl) Shutdown() {
	port.logger.Info("Gracefully shutting down Task Manager")
}
