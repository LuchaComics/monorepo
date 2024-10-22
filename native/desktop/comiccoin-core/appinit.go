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
	preferences := PreferencesInstance()
	defaultDataDir := preferences.GetDefaultDataDirectory()
	return defaultDataDir
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
	preferences := PreferencesInstance()
	err := preferences.SetDataDirectory(newDataDirectory)
	if err != nil {
		a.logger.Error("Failed setting data directory",
			slog.Any("error", err))
		return err
	}

	// Re-attempt the startup now that we have the data directory set.
	a.startup(a.ctx)
	return nil
}

func (a *App) ShutdownApp() {
	runtime.Quit(a.ctx)
}

func (a *App) GetIsBlockhainNodeRunning() bool {
	return a.isBlockchainNodeRunning
}
