package task

import (
	"log/slog"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/config"
)

type TaskManager interface {
	Run()
	Shutdown()
}

type taskManagerImpl struct {
	cfg    *config.Configuration
	logger *slog.Logger
}

func NewTaskManager(
	cfg *config.Configuration,
	logger *slog.Logger,

) TaskManager {
	port := &taskManagerImpl{
		cfg:    cfg,
		logger: logger,
	}
	return port
}

func (port *taskManagerImpl) Run() {
	// ctx := context.Background()
	port.logger.Info("Running Task Manager")

}

func (port *taskManagerImpl) Shutdown() {
	port.logger.Info("Gracefully shutting down Task Manager")
}
