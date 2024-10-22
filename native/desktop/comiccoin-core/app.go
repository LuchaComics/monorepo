package main

import (
	"context"
	"fmt"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/config"
)

// App struct
type App struct {
	ctx    context.Context
	pageID int
	config *config.Config
}

// NewApp creates a new App application struct
func NewApp() *App {
	cfg := &config.Config{}
	return &App{
		config: cfg,
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
