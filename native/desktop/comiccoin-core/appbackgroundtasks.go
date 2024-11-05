package main

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/config/constants"
	taskmnghandler "github.com/LuchaComics/monorepo/native/desktop/comiccoin/interface/task/handler"
)

func (a *App) startBackgroundTasks() {
	ctx := a.ctx
	a.logger.Info("Running Task Manager")

	go func() {
		a.logger.Info("Running mempool (receive) service...")
		for {
			taskErr := a.mempoolReceiveTaskHandler.Execute(ctx)
			if taskErr != nil {
				a.logger.Error("failed executing mempool receive task, restarting task in 1 minute...", slog.Any("error", taskErr))
				time.Sleep(1 * time.Second)
			}
		}
	}()

	go func() {
		a.logger.Info("Running mempool (send) service...")
		for {
			taskErr := a.mempoolBatchSendTaskHandler.Execute(ctx)
			if taskErr != nil {
				a.logger.Error("failed executing mempool batch send task, restarting task in 1 minute...", slog.Any("error", taskErr))
				time.Sleep(1 * time.Minute)
			}
			time.Sleep(1 * time.Second)
		}
	}()
	go func() {
		if a.config.Blockchain.EnableMiner {
			if a.config.Blockchain.ConsensusProtocol == constants.ConsensusPoW {
				a.logger.Info("Running PoW mining service...")
				for {
					taskErr := a.proofOfWorkMiningTaskHandler.Execute(ctx)
					if taskErr != nil {
						a.logger.Error("failed executing mining task, restarting task in 1 minute...", slog.Any("error", taskErr))
						time.Sleep(1 * time.Minute)
					}
					time.Sleep(1 * time.Second)
				}
			} else if a.config.Blockchain.ConsensusProtocol == constants.ConsensusPoA {
				a.logger.Info("Running PoA mining service...")
				for {
					taskErr := a.proofOfAuthorityMiningTaskHandler.Execute(ctx)
					if taskErr != nil {
						a.logger.Error("failed executing mining task, restarting task in 1 minute...", slog.Any("error", taskErr))
						time.Sleep(1 * time.Minute)
					}
					time.Sleep(1 * time.Second)
				}
			} else {
				a.logger.Info("Skipped running the mining service...")
			}
		}
	}()

	go func() {
		if a.config.Blockchain.ConsensusProtocol == constants.ConsensusPoW {
			a.logger.Info("Running PoW validation service...")
			for {
				taskErr := a.proofOfWorkValidationTaskHandler.Execute(ctx)
				if taskErr != nil {
					a.logger.Error("failed executing validation task, restarting task in 1 minute...", slog.Any("error", taskErr))
					time.Sleep(1 * time.Minute)
				}
				time.Sleep(1 * time.Second)
			}
		} else if a.config.Blockchain.ConsensusProtocol == constants.ConsensusPoA {
			a.logger.Info("Running PoA validation service...")
			for {
				taskErr := a.proofOfAuthorityValidationTaskHandler.Execute(ctx)
				if taskErr != nil {
					a.logger.Error("failed executing validation task, restarting task in 1 minute...", slog.Any("error", taskErr))
					time.Sleep(1 * time.Minute)
				}
				time.Sleep(1 * time.Second)
			}
		} else {
			a.logger.Info("Skipped running the mining service...")
		}
	}()

	go func(consensus *taskmnghandler.MajorityVoteConsensusServerTaskHandler, loggerp *slog.Logger) {
		a.logger.Info("Running consensus server...")
		ctx := context.Background()
		for {
			if err := consensus.Execute(ctx); err != nil {
				loggerp.Error("consensus error", slog.Any("error", err))
			}
			// DEVELOPERS NOTE:
			// No need for delays, automatically start executing again.
			loggerp.Debug("blockchain consensus serving done, excuting again ...")
		}
	}(a.majorityVoteConsensusServerTaskHandler, a.logger)

	go func(consensus *taskmnghandler.MajorityVoteConsensusClientTaskHandler, loggerp *slog.Logger) {
		loggerp.Info("Running consensus client...")
		ctx := context.Background()
		for {
			if err := consensus.Execute(ctx); err != nil {
				loggerp.Error("consensus error", slog.Any("error", err))
			}
			time.Sleep(time.Duration(a.config.Blockchain.ConsensusPollingDelayInMinutes) * time.Minute)
			loggerp.Debug(fmt.Sprintf("blockchain consensus client ran, will run again in %v minutes...", a.config.Blockchain.ConsensusPollingDelayInMinutes))
		}
	}(a.majorityVoteConsensusClientTaskHandler, a.logger)

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
			a.logger.Debug("block data dto server executing again ...")
		}
	}(a.blockDataDTOServerTaskHandler, a.logger)

	go func(client *taskmnghandler.SignedIssuedTokenClientServiceTaskHandler, loggerp *slog.Logger) {
		loggerp.Info("Running issued token dto client...")
		ctx := context.Background()
		for {
			if err := client.Execute(ctx); err != nil {
				loggerp.Error("issued token client error",
					slog.Any("error", err))
				time.Sleep(10 * time.Second)
				continue
			}
			// DEVELOPERS NOTE:
			// No need for delays, automatically start executing again.
			loggerp.Debug("issued token dto client executing again ...")
		}
	}(a.signedIssuedTokenClientServiceTaskHandler, a.logger)
}

func (a *App) stopBackgroundTasks() {
	a.logger.Info("Gracefully shutting down Task Manager")
}
