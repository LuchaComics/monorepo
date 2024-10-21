package views

import (
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin-core/constants"
)

type StartupView struct {
	window      fyne.Window
	preferences fyne.Preferences

	// nextPageID represents a recieve operation channel used by this view set the
	nextPageID chan int
}

func NewStartupView(
	w fyne.Window,
	pref fyne.Preferences) Viewer {
	v := &StartupView{
		window:      w,
		preferences: pref,
		nextPageID:  make(chan int),
	}

	return v
}

func (view *StartupView) Render() *fyne.Container {
	view.window.Resize(fyne.NewSize(400, 400)) // Set the window size

	loadingLabel := widget.NewLabel("Loading ComicCoin v1...")
	loadingLabel.Alignment = fyne.TextAlignCenter
	loadingLabel.TextStyle = fyne.TextStyle{Bold: true}

	progressBar := widget.NewProgressBar()
	progressBar.SetValue(0) // Initialize progress bar value

	content := container.NewVBox(
		layout.NewSpacer(),
		loadingLabel,
		widget.NewSeparator(),
		progressBar,
		layout.NewSpacer(),
	)

	view.window.SetContent(content)

	//TODO: IMPL.

	go func() {
		for i := 0; i <= 100; i++ {
			time.Sleep(1 * time.Millisecond) // Simulate loading process
			progressBar.SetValue(float64(i) / 100)
			if i == 100 {
				// Loading complete, close the window
				view.nextPageID <- constants.PageIDMainAppView
			}
		}
	}()

	return content
}

func (view *StartupView) WaitUntilReadyToTransition() int {
	// go func() {
	// 	//TODO: Write code to wait for synchronization w/ blockchain network.
	// 	time.Sleep(5 * time.Second)
	// 	view.nextPageID <- constants.PageIDTabNavigatorAppView
	// }()
	return <-view.nextPageID
}
