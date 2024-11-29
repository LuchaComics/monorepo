package main

import (
	"fmt"
	"log"
	"log/slog"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/common/httperror"
	"github.com/wailsapp/wails/v2/pkg/runtime"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin-registry/repo"
)

func (a *App) GetDataDirectoryFromPreferences() string {
	preferences := PreferencesInstance()
	dataDir := preferences.DataDirectory
	return dataDir
}

func (a *App) GetDefaultDataDirectory() string {
	preferences := PreferencesInstance()
	defaultDataDir := preferences.GetDefaultDataDirectory()
	return defaultDataDir
}

func (a *App) GetNFTStoreAPIKeyFromPreferences() string {
	preferences := PreferencesInstance()
	nftStoreAPIKey := preferences.NFTStoreAPIKey
	return nftStoreAPIKey
}

func (a *App) GetNFTStoreRemoteAddressFromPreferences() string {
	preferences := PreferencesInstance()
	nftStoreRemoteAddress := preferences.NFTStoreRemoteAddress
	return nftStoreRemoteAddress
}

func (a *App) GetDataDirectoryFromDialog() string {
	// Initialize Wails runtime
	result, err := runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "Please select were to save the blockchain",
	})
	if err != nil {
		a.logger.Error("Failed opening directory dialog",
			slog.Any("error", err))
		log.Fatalf("%v", err)
	}
	return result
}

func (a *App) SaveDataDirectory(newDataDirectory string) error {
	// Defensive code
	if newDataDirectory == "" {
		return fmt.Errorf("failed saving data directory because: %v", "data directory is empty")
	}
	preferences := PreferencesInstance()
	err := preferences.SetDataDirectory(newDataDirectory)
	if err != nil {
		a.logger.Error("Failed setting data directory",
			slog.Any("error", err))
		return err
	}

	// Re-attempt the startup now that we have the data directory set.
	a.logger.Debug("Data directory was set by user",
		slog.Any("data_directory", newDataDirectory))

	return nil
}

func (a *App) SaveNFTStoreAPIKey(nftStoreAPIKey string) error {
	// Defensive code
	if nftStoreAPIKey == "" {
		return fmt.Errorf("failed saving nft store api key because: %v", "value is empty")
	}
	preferences := PreferencesInstance()
	err := preferences.SetNFTStoreAPIKey(nftStoreAPIKey)
	if err != nil {
		a.logger.Error("Failed setting nft store api key",
			slog.Any("error", err))
		return err
	}

	// Re-attempt the startup now that we have the data directory set.
	a.logger.Debug("NFT store api key was set by user",
		slog.Any("api_key", nftStoreAPIKey))
	return nil
}

func (a *App) SaveNFTStoreRemoteAddress(nftStoreRemoteAddress string) error {
	// Defensive code
	if nftStoreRemoteAddress == "" {
		return fmt.Errorf("failed saving nft store remote address because: %v", "value is empty")
	}
	preferences := PreferencesInstance()
	err := preferences.SetNFTStoreRemoteAddress(nftStoreRemoteAddress)
	if err != nil {
		a.logger.Error("Failed setting nft store remote address",
			slog.Any("error", err))
		return err
	}

	// Re-attempt the startup now that we have value set.
	a.logger.Debug("NFT store remote address was set by user",
		slog.Any("api_key", nftStoreRemoteAddress))
	return nil
}

func (a *App) SaveNFTStoreConfigVariables(nftStoreAPIKey string, nftStoreRemoteAddress string) error {
	//
	// STEP 1:
	// Validation.
	//

	e := make(map[string]string)
	if nftStoreAPIKey == "" {
		e["nftStoreAPIKey"] = "missing value"
	}
	if nftStoreRemoteAddress == "" {
		e["nftStoreRemoteAddress"] = "missing value"
	}
	if len(e) != 0 {
		// If any fields are missing, log an error and return a bad request error.
		a.logger.Warn("Failed validating",
			slog.Any("error", e))
		return httperror.NewForBadRequest(&e)
	}

	//
	// STEP 2:
	// Confirm the remote address works.
	//

	tempNFTAssetRepo := repo.NewNFTAssetRepoWithConfiguration(a.logger, nftStoreRemoteAddress, nftStoreAPIKey)
	version, err := tempNFTAssetRepo.Version(a.ctx)
	if err != nil {
		return httperror.NewForBadRequestWithSingleField("nftStoreRemoteAddress", fmt.Sprintf("%v", err))
	}
	if version != "1.0" {
		return httperror.NewForBadRequestWithSingleField("nftStoreRemoteAddress", fmt.Sprintf("Wrong version: %v", version))
	}

	//
	// STEP 3:
	// Save to preferences.
	//

	if err := a.SaveNFTStoreRemoteAddress(nftStoreRemoteAddress); err != nil {
		a.logger.Error("Failed setting nft store remote address",
			slog.Any("error", err))
		return err
	}
	if err := a.SaveNFTStoreAPIKey(nftStoreAPIKey); err != nil {
		a.logger.Error("Failed setting nft store api key",
			slog.Any("error", err))
		return err
	}
	return nil

}

// StartupApp method is to be used once all missing environment / configuration
// variables have been set and thus we can begin running the core of the
// application which is dependent on these variables set.
func (a *App) StartupApp() {
	a.startup(a.ctx)
}

func (a *App) ShutdownApp() {
	runtime.Quit(a.ctx)
}
