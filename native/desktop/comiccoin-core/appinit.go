package main

import (
	"fmt"
	"log"

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
		log.Fatalf("%v", err)
	}
	return result
}

// Greet returns a greeting for the given name
func (a *App) SaveDataDirectory(newDataDirectory string) error {
	if newDataDirectory == "" {
		return fmt.Errorf("failed saving data directory because: %v", "data directory is empty")
	}
	fmt.Println("saved", newDataDirectory)
	preferences := PreferencesInstance()
	return preferences.SetDataDirectory(newDataDirectory)
}

func (a *App) ShutdownApp() {
	runtime.Quit(a.ctx)
}
