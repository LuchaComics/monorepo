// navigator package
package navigator

import (
	"fmt"
	"log"

	"fyne.io/fyne/v2"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin-core/constants"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin-core/views"
)

type Navigator struct {
	window          fyne.Window
	preferences     fyne.Preferences
	app             fyne.App
	pickDataDirView views.Viewer
	startupView     views.Viewer
	walletView      views.Viewer
}

func NewNavigator(
	window fyne.Window,
	pref fyne.Preferences,
	a fyne.App,
	v1 views.Viewer,
	v2 views.Viewer,
	v3 views.Viewer,
) *Navigator {
	return &Navigator{
		window:          window,
		preferences:     pref,
		app:             a,
		pickDataDirView: v1,
		startupView:     v2,
		walletView:      v3,
	}
}

func (nav *Navigator) RunMainRuntimeLoop() {
	var nextPageID int

	//
	// STEP 1:
	// When we load up the application for the first time, we need to check
	// that the user has selected were they want to save their application data.
	// If the user has not selected a location then provide the GUI to allow
	// then to select, else if they already did the selection then skip.
	//

	hasSetDataDir := nav.preferences.Bool(constants.PreferenceKeyHasSetDataDirectory)
	if !hasSetDataDir {
		content := nav.pickDataDirView.Render()
		nav.window.SetContent(content)
		nextPageID = nav.pickDataDirView.WaitUntilReadyToTransition()
	} else {
		nextPageID = constants.PageIDStartupView
	}

	//
	// STEP 2:
	// Run the main runtime loop.
	//

	for {
		// Enforce application window size.
		nav.window.Resize(fyne.NewSize(constants.DefaultScreenWidth, constants.DefaultScreenHeight))

		switch nextPageID {
		case constants.PageIDExit:
			{
				fmt.Println("Closing application...")
				nav.app.Quit()
			}
		case constants.PageIDStartupView:
			{
				fmt.Println("Starting application...")
				nav.window.SetContent(nav.startupView.Render())
				nextPageID = nav.startupView.WaitUntilReadyToTransition()
			}
		case constants.PageIDMainAppView:
			{
				fmt.Println("Starting wallet app...")
				nav.window.SetContent(nav.walletView.Render())
				nextPageID = nav.walletView.WaitUntilReadyToTransition()
			}
		default:
			log.Fatalf("unknown page id: %v\n", nextPageID)
		}
	}
}

// v1 := views.NewPickBlockchainStorageLocationView(w)
//
// w.SetContent(v1)
