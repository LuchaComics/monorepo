package task

import (
	"context"
	"log/slog"
	"time"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/config"
	taskhandler "github.com/LuchaComics/monorepo/cloud/comiccoin-authority/interface/task/handler"
)

type TaskManager interface {
	Run()
	Shutdown()
}

type taskManagerImpl struct {
	cfg                                           *config.Configuration
	logger                                        *slog.Logger
	proofOfAuthorityConsensusMechanismTaskHandler *taskhandler.ProofOfAuthorityConsensusMechanismTaskHandler
}

func NewTaskManager(
	cfg *config.Configuration,
	logger *slog.Logger,
	task1 *taskhandler.ProofOfAuthorityConsensusMechanismTaskHandler,

) TaskManager {
	port := &taskManagerImpl{
		cfg:    cfg,
		logger: logger,
		proofOfAuthorityConsensusMechanismTaskHandler: task1,
	}
	return port
}

func (port *taskManagerImpl) Run() {
	// ctx := context.Background()
	port.logger.Info("Running Task Manager")
	backgroundCtx := context.Background()

	go func(task *taskhandler.ProofOfAuthorityConsensusMechanismTaskHandler, loggerp *slog.Logger) {
		loggerp.Info("Starting PoA consensus mechanism...")

		for {
			if err := task.Execute(backgroundCtx); err != nil {
				loggerp.Error("Failed executing PoA consensus mechanism",
					slog.Any("error", err))
				time.Sleep(10 * time.Second)
				continue
			}
			// DEVELOPERS NOTE:
			// No need for delays, automatically start executing again.
			port.logger.Debug("poa consensus mechanism will run again ...")
			time.Sleep(60 * time.Second)
		}
	}(port.proofOfAuthorityConsensusMechanismTaskHandler, port.logger)
}

func (port *taskManagerImpl) Shutdown() {
	port.logger.Info("Gracefully shutting down Task Manager")
}
