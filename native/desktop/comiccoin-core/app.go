package main

import (
	"context"
	"fmt"
	"log"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/config"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx    context.Context
	pageID string
	config *config.Config
}

// NewApp creates a new App application struct
func NewApp() *App {
	cfg := &config.Config{}

	return &App{
		config: cfg,
		pageID: "welcome",
	}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	fmt.Println("Starting now...")
}

func (b *App) shutdown(ctx context.Context) {
	fmt.Println("Shutting down now...")
}

// Greet returns a greeting for the given name
func (a *App) GetPageID() string {
	return a.pageID
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	a.pageID = "home"
	return fmt.Sprintf("Hello %s, It's show time!", name)
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
func (a *App) SaveDataDirectory(dataDirectory string) error {
	fmt.Println("saved", dataDirectory)
	return nil
}
