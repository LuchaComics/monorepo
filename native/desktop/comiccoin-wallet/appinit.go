package main

import (
	"fmt"
	"log"
	"log/slog"

	"github.com/wailsapp/wails/v2/pkg/runtime"

	pref "github.com/LuchaComics/monorepo/native/desktop/comiccoin-wallet/common/preferences"
)

func (a *App) GetDataDirectoryFromPreferences() string {
	preferences := pref.PreferencesInstance()
	dataDir := preferences.DataDirectory
	return dataDir
}

func (a *App) GetDefaultDataDirectory() string {
	return pref.GetDefaultDataDirectory()
}

func (a *App) GetNFTStoreRemoteAddressFromPreferences() string { //TODO: Refactor `GetNFTStoreRemoteAddressFromPreferences` to `GetNFTStorageAddressFromPreferences`.
	preferences := pref.PreferencesInstance()
	nftStoreRemoteAddress := preferences.NFTStorageAddress
	return nftStoreRemoteAddress
}

// Greet returns a greeting for the given name
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

// Greet returns a greeting for the given name
func (a *App) SaveDataDirectory(newDataDirectory string) error {
	// Defensive code
	if newDataDirectory == "" {
		return fmt.Errorf("failed saving data directory because: %v", "data directory is empty")
	}
	preferences := pref.PreferencesInstance()
	err := preferences.SetDataDirectory(newDataDirectory)
	if err != nil {
		a.logger.Error("Failed setting data directory",
			slog.Any("error", err))
		return err
	}

	// Re-attempt the startup now that we have the data directory set.
	a.logger.Debug("Data directory was set by user",
		slog.Any("data_directory", newDataDirectory))
	a.startup(a.ctx)
	return nil
}

func (a *App) SaveNFTStoreRemoteAddress(nftStorageAddress string) error { //TODO: Refactor `SaveNFTStoreRemoteAddress` to `SetNFTStorageAddress`.
	// Defensive code
	if nftStorageAddress == "" {
		return fmt.Errorf("failed saving nft storage address because: %v", "value is empty")
	}
	preferences := pref.PreferencesInstance()
	err := preferences.SetNFTStorageAddress(nftStorageAddress)
	if err != nil {
		a.logger.Error("Failed setting nft storage address",
			slog.Any("nft_storage_address", nftStorageAddress),
			slog.Any("error", err))
		return err
	}

	// Re-attempt the startup now that we have value set.
	a.logger.Debug("NFT storage address was set by user",
		slog.Any("nft_storage_address", nftStorageAddress))
	return nil
}

func (a *App) ShutdownApp() {
	runtime.Quit(a.ctx)
}

func (a *App) GetIsBlockhainNodeRunning() bool {
	return true //TODO: REMOVE
}
