package main

import (
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin-core/navigator"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin-core/views"
	view.wallet_view "github.com/LuchaComics/monorepo/native/desktop/comiccoin-core/views/wallet"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin-core/constants"
)

func main() {
	//
	// STEP 1:
	// Initialize the application and set configuration.
	//

	a := app.NewWithID("com.cpscapsule.comiccoin")
	w := a.NewWindow("ComicCoin Core")
	w.Resize(fyne.NewSize(constants.DefaultScreenWidth, constants.DefaultScreenHeight))

	// // DEVELOPERS NOTE:
	// // Uncomment this if you want to clear the preferences on startup.
	// a.Preferences().RemoveValue(constants.PreferenceKeyHasSetDataDirectory)
	// a.Preferences().RemoveValue(constants.PreferenceKeyDataDirectory)

	// Set first-time preferences.
	a.Preferences().BoolWithFallback(constants.PreferenceKeyHasSetDataDirectory, false)
	a.Preferences().StringWithFallback(constants.PreferenceKeyDataDirectory, constants.DefaultDataDirectoryPath)

	//
	// STEP 2:
	// Load up our dependencies.
	//

	// -- Pick data directory page. ---
	pickDataDirView := views.NewPickDataDirectoryView(w, a.Preferences())

	// --- Startup page ---
	startupView := views.NewStartupView(w, a.Preferences())

	// --- Wallet page ---
	walletOverviewTab := view.wallet_view.NewWalletViewOverviewTab(w, a.Preferences())
	walletSendTab := view.wallet_view.NewWalletViewSendTab(w, a.Preferences())
	walletReceiveTab := view.wallet_view.NewWalletViewReceiveTab(w, a.Preferences())
	walletTransactionsTab := view.wallet_view.NewWalletViewTransactionsTab(w, a.Preferences())
	walletMoreTab := view.wallet_view.NewWalletViewMoreTab(w, a.Preferences())
	walletView := view.wallet_view.NewWalletView(
		w,
		a.Preferences(),
		walletOverviewTab,
		walletSendTab,
		walletReceiveTab,
		walletTransactionsTab,
		walletMoreTab,
	)

	navigator := navigator.NewNavigator(
		w,
		a.Preferences(),
		a,
		pickDataDirView,
		startupView,
		walletView,
	)

	//
	// STEP 3:
	// Load up the graceful shutdown code.
	//

	defer func() {
		log.Println("tidy up...")
	}()

	//
	// STEP 4:
	// Execute our application.
	//
	go navigator.RunMainRuntimeLoop()

	//
	// STEP 5:
	// Render the GUI elements based on our application state and interactions.
	//

	w.ShowAndRun()
}
