package main

import (
	"fmt"
	"log"
	"log/slog"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

func (a *App) GetDataDirectoryFromPreferences() string {
	preferences := PreferencesInstance()
	dataDir := preferences.DataDirectory
	return dataDir
}

func (a *App) GetDefaultDataDirectory() string {
	return GetDefaultDataDirectory()
}

func (a *App) GetNFTStorageAddressFromPreferences() string {
	preferences := PreferencesInstance()
	nftStoreRemoteAddress := preferences.NFTStorageAddress
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

	// Add the remaining configuration settings.
	if err := preferences.SetChainID(ComicCoinChainID); err != nil {
		a.logger.Error("Failed setting chain ID",
			slog.Any("error", err))
		return err
	}
	if err := preferences.SetAuthorityAddress(ComicCoinAuthorityAddress); err != nil {
		a.logger.Error("Failed setting authority address",
			slog.Any("error", err))
		return err
	}
	if err := preferences.SetNFTStorageAddress(ComicCoinNFTStorageAddress); err != nil {
		a.logger.Error("Failed setting NFT storage address",
			slog.Any("error", err))
		return err
	}

	a.startup(a.ctx)
	return nil
}

func (a *App) SetNFTStorageAddress(nftStorageAddress string) error {
	// Defensive code
	if nftStorageAddress == "" {
		return fmt.Errorf("failed saving nft storage address because: %v", "value is empty")
	}
	preferences := PreferencesInstance()
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

func (a *App) SetAuthorityAddress(authorityAddress string) error {
	// Defensive code
	if authorityAddress == "" {
		return fmt.Errorf("failed saving authority address because: %v", "value is empty")
	}
	preferences := PreferencesInstance()
	err := preferences.SetAuthorityAddress(authorityAddress)
	if err != nil {
		a.logger.Error("Failed setting authority address",
			slog.Any("authority_address", authorityAddress),
			slog.Any("error", err))
		return err
	}

	// Re-attempt the startup now that we have value set.
	a.logger.Debug("Authority address was set by user",
		slog.Any("authority_address", authorityAddress))
	return nil
}

func (a *App) ShutdownApp() {
	runtime.Quit(a.ctx)
}

func (a *App) GetIsBlockhainNodeRunning() bool {
	return true //TODO: REMOVE
}
