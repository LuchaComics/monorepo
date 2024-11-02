package main

import (
	"embed"
	"fmt"
	"net/http"
	"os"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/mac"
	// "github.com/wailsapp/wails/v2/pkg/options/windows"
	// "github.com/wailsapp/wails/v2/pkg/options/linux"
)

//go:embed all:frontend/dist
var assets embed.FS

// FileLoader structure represents the dynamic asset handler to allow our
// application to access files in the computer via HTTP URL requests. For more
// information see https://wails.io/docs/guides/dynamic-assets/.
type FileLoader struct {
	http.Handler
}

// NewFileLoader method is a constructor of the dynamic asset handler.
func NewFileLoader() *FileLoader {
	return &FileLoader{}
}

// ServeHTTP method is required by `go wails` for the dynamic asset handler.
func (h *FileLoader) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	// DEVELOPERS NOTE:
	// See `Dynamic Assets` via https://wails.io/docs/guides/dynamic-assets/.
	var err error
	requestedFilename := req.URL.Path
	// log.Println("FileLoader ---> Requesting file:", requestedFilename)
	fileData, err := os.ReadFile(requestedFilename)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		// log.Println("FileLoader ---> err:", err)
		res.Write([]byte(fmt.Sprintf("Could not load file %s", requestedFilename)))
	}

	res.Write(fileData)
}

func main() {
	// Create an instance of the app structure
	app := NewApp()

	// Create application with options
	err := wails.Run(&options.App{
		Title:  "ComicCoin Core",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			// Static asset handler - This loads up data that come preloaded with our application like logos, etc.
			Assets: assets,

			// Dynamic asset handler - This enables our app to access dynamically created data by the user.
			Handler: NewFileLoader(),
		},
		// BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1}, //FUTURE: Adjustable feature. Currently the default is white background. (https://wails.io/docs/reference/options/#backgroundcolour)
		OnStartup:  app.startup,
		OnShutdown: app.shutdown,
		Mac: &mac.Options{
			WebviewIsTransparent: true,
		},
		Bind: []interface{}{
			app,
		},
		// DisableResize: true,
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
