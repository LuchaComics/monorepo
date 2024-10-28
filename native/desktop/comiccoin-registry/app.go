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

	nftRepo *repo.NFTRepo
}

// NewApp creates a new App application struct
func NewApp() *App {
	logger := logger.NewLogger()
	return &App{
		logger:  logger,
		nftRepo: nil,
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

	nftByTokenIDDB := disk.NewDiskStorage(dataDir, "nft_by_tokenid", a.logger)
	nftByMetadataURIDB := disk.NewDiskStorage(dataDir, "nft_by_metadatauri", a.logger)
	nftRepo := repo.NewNFTRepo(a.logger, nftByTokenIDDB, nftByMetadataURIDB)
	a.nftRepo = nftRepo
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}
