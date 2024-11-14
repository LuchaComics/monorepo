package main

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin-registry/common/kmutexutil"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin-registry/common/logger"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin-registry/domain"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin-registry/repo"
)

// App struct
type App struct {
	ctx context.Context

	// Logger instance which provides detailed debugging information along
	// with the console log messages.
	logger       *slog.Logger
	kmutex       kmutexutil.KMutexProvider
	fileBaseRepo domain.FileBaseRepository
}

// NewApp creates a new App application struct
func NewApp() *App {
	logger := logger.NewProvider()
	kmutex := kmutexutil.NewKMutexProvider()
	return &App{
		logger:       logger,
		kmutex:       kmutex,
		fileBaseRepo: nil,
	}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	// Ensure that this function executes only one time and never concurrently.
	a.kmutex.Acquire("startup")
	defer a.kmutex.Release("startup")

	a.ctx = ctx
	a.logger.Debug("Startup beginning...")

	// DEVELOPERS NOTE:
	// Before we startup our app, we need to make sure the `data directory` is
	// set for this application by the user, else stop the app startup
	// proceedure. This is done on purpose because we need the user to specify
	// the location they want to store instead of having one automatically set.
	preferences := PreferencesInstance()
	dataDir := preferences.DataDirectory
	if dataDir == "" {
		a.logger.Debug("Startup halted: need to specify data directory")
		return
	}

	if a.fileBaseRepo == nil {
		a.logger.Debug("Connecting to FileBase...")
		newConfig := preferences.NFTStoreSettings

		filebaseConfig := repo.NewFileBaseRepoConfigurationProvider(
			newConfig["apiVersion"],
			newConfig["accessKeyId"],
			newConfig["secretAccessKey"],
			newConfig["endpoint"],
			newConfig["region"],
			newConfig["s3ForcePathStyle"],
		)
		fileBaseRepo := repo.NewFileBaseRepo(filebaseConfig, a.logger)
		a.fileBaseRepo = fileBaseRepo
		a.logger.Debug("FileBase connected!")
	}

}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}
