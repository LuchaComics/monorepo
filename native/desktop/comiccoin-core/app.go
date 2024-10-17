package main

import (
	"context"
	"fmt"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/config"
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
