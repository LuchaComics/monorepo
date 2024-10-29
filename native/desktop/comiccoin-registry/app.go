package main

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/common/logger"
	disk "github.com/LuchaComics/monorepo/native/desktop/comiccoin/common/storage/disk/leveldb"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin-registry/repo"
)

// App struct
type App struct {
	ctx context.Context

	// Logger instance which provides detailed debugging information along
	// with the console log messages.
	logger *slog.Logger

	tokenRepo *repo.TokenRepo

	ipfsRepo *repo.IPFSRepo

	latestTokenIDRepo *repo.LastestTokenIDRepo
}

// NewApp creates a new App application struct
func NewApp() *App {
	logger := logger.NewLogger()
	return &App{
		logger:    logger,
		tokenRepo: nil,
	}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
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

	tokenByTokenIDDB := disk.NewDiskStorage(dataDir, "token_by_id", a.logger)
	tokenByMetadataURIDB := disk.NewDiskStorage(dataDir, "token_by_metadata_uri", a.logger)
	tokenRepo := repo.NewTokenRepo(a.logger, tokenByTokenIDDB, tokenByMetadataURIDB)
	a.tokenRepo = tokenRepo

	ipfsNode := repo.NewIPFSRepo(a.logger, "http://localhost:5002")
	a.ipfsRepo = ipfsNode

	latestTokenIDDB := disk.NewDiskStorage(dataDir, "latest_token_id", a.logger)
	latestTokenIDRepo := repo.NewLastestTokenIDRepo(
		a.logger,
		latestTokenIDDB)
	a.latestTokenIDRepo = latestTokenIDRepo
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}
