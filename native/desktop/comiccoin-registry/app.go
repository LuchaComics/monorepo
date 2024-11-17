package main

import (
	"context"
	"log/slog"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/common/kmutexutil"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/common/logger"
	disk "github.com/LuchaComics/monorepo/cloud/comiccoin-authority/common/storage/disk/leveldb"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin-registry/domain"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin-registry/repo"
)

// App struct
type App struct {
	ctx context.Context

	// Logger instance which provides detailed debugging information along
	// with the console log messages.
	logger *slog.Logger

	kmutex kmutexutil.KMutexProvider

	tokenRepo *repo.TokenRepo

	nftAssetRepo domain.NFTAssetRepository

	latestTokenIDRepo *repo.LastestTokenIDRepo
}

// NewApp creates a new App application struct
func NewApp() *App {
	logger := logger.NewProvider()
	kmutex := kmutexutil.NewKMutexProvider()
	return &App{
		logger:            logger,
		kmutex:            kmutex,
		tokenRepo:         nil,
		nftAssetRepo:      nil,
		latestTokenIDRepo: nil,
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

	if a.tokenRepo == nil {
		a.logger.Debug("Loading disk storage: token_by_id")
		tokenByTokenIDDB := disk.NewDiskStorage(dataDir, "token_by_id", a.logger)

		a.logger.Debug("Loading disk storage: token_by_metadata_uri")
		tokenByMetadataURIDB := disk.NewDiskStorage(dataDir, "token_by_metadata_uri", a.logger)

		a.logger.Debug("Loading repo: tokenRepo")
		tokenRepo := repo.NewTokenRepo(a.logger, tokenByTokenIDDB, tokenByMetadataURIDB)
		a.tokenRepo = tokenRepo
	}

	if a.nftAssetRepo == nil {
		nftStoreRemoteAddress := preferences.NFTStoreRemoteAddress
		nftStoreAPIKey := preferences.NFTStoreAPIKey
		a.logger.Debug("Loading repo: NFT asset repo",
			slog.Any("nftStoreRemoteAddress", nftStoreRemoteAddress),
			slog.Any("nftStoreAPIKey", nftStoreAPIKey))

		nftAssetRepoConfig := repo.NewNFTAssetRepoConfigurationProvider(nftStoreRemoteAddress, nftStoreAPIKey)
		nftAssetRepo := repo.NewNFTAssetRepo(nftAssetRepoConfig, a.logger)
		a.nftAssetRepo = nftAssetRepo
	}

	if a.latestTokenIDRepo == nil {
		a.logger.Debug("Loading disk storage: latest_token_id")
		latestTokenIDDB := disk.NewDiskStorage(dataDir, "latest_token_id", a.logger)

		a.logger.Debug("Loading repo: latestTokenIDRepo")
		latestTokenIDRepo := repo.NewLastestTokenIDRepo(
			a.logger,
			latestTokenIDDB)
		a.latestTokenIDRepo = latestTokenIDRepo
	}

	a.logger.Debug("startup: finished...")
}
