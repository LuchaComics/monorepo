package views

import (
	"fmt"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin-core/constants"
)

type PickDataDirectoryView struct {
	window      fyne.Window
	preferences fyne.Preferences

	// nextPageID represents a recieve operation channel used by this view set the
	nextPageID chan int
}

func NewPickDataDirectoryView(
	w fyne.Window,
	pref fyne.Preferences) Viewer {
	v := &PickDataDirectoryView{
		window:      w,
		preferences: pref,
		nextPageID:  make(chan int),
	}

	return v
}

func (view *PickDataDirectoryView) Render() fyne.CanvasObject {
	view.window.Resize(fyne.NewSize(constants.DefaultScreenWidth, constants.DefaultScreenHeight))

	lb1 := widget.NewLabel("Welcome to ComicCoin Core.")
	lb2 := widget.NewLabel("As this is the first time the program is launched, you can choose where ComicCoin Core will store its data.")
	lb2.Wrapping = fyne.TextWrapWord
	lb3 := widget.NewLabel("ComicCoin Core will download and store a copy of the ComicCoin block chain. Approximately 1 MB of data will be stored in this directory. The wallet will also be stored in this directory.")
	lb3.Wrapping = fyne.TextWrapWord

	//
	// Pick default or custom data directory
	//

	// Directory recorded by the user.
	dataDirectory := binding.NewString()
	dataDirectory.Set(view.preferences.String(constants.PreferenceKeyDataDirectory))

	dataDirectoryEntry := widget.NewEntryWithData(dataDirectory)
	dataDirectoryEntry.Disable()
	dataDirectoryEntry.Wrapping = fyne.TextTruncate

	pickDataDirectoryBtn := widget.NewButtonWithIcon("", theme.MoreHorizontalIcon(), func() {
		dialog.ShowFolderOpen(func(uri fyne.ListableURI, err error) {
			if err != nil {
				dialog.ShowError(err, view.window)
				return
			}
			if uri != nil {
				dataDirectoryPath := uri.Path()
				fmt.Println("Selected directory:", dataDirectoryPath)
				dataDirectory.Set(dataDirectoryPath)
			}
		}, view.window)
	})

	isDefaultDataDirectory := binding.NewBool()
	radio := widget.NewRadioGroup([]string{"Use the default data directory.", "Use a custom data directory."}, func(value string) {
		switch value {
		case "Use the default data directory.":
			isDefaultDataDirectory.Set(true)
			dataDirectory.Set(constants.DefaultDataDirectoryPath)
			dataDirectoryEntry.Disable()
			pickDataDirectoryBtn.Disable()

		case "Use a custom data directory.":
			isDefaultDataDirectory.Set(false)
			dataDirectoryEntry.Enable()
			pickDataDirectoryBtn.Enable()
		}
		log.Println("Radio set to", value)
	})
	radio.SetSelected("Use the default data directory.")

	entryBox := container.NewBorder(nil, nil, nil, pickDataDirectoryBtn, dataDirectoryEntry)

	lb4 := widget.NewLabel("When you click OK, ComicCoin Core will begin to download and process the full ComicCoin block chain (1 MB) starting with the earliest transactions in 2024 when ComicCoin initially launched.")
	lb4.Wrapping = fyne.TextWrapWord
	lb5 := widget.NewLabel("This initial synchronisation is very demanding, and may expose hardware problems with your computer that had previously gone unnoticed. Each time you run ComicCoin Core, it will continue downloading where it left off.")
	lb5.Wrapping = fyne.TextWrapWord

	cancelBtn := widget.NewButtonWithIcon("Cancel", theme.CancelIcon(), func() {
		log.Println("tapped CANCEL")
		view.preferences.SetBool(constants.PreferenceKeyHasSetDataDirectory, false)
		view.nextPageID <- constants.PageIDExit
	})
	okBtn := widget.NewButtonWithIcon("OK", theme.ConfirmIcon(), func() {
		log.Println("tapped OK")

		// Set the preference for the data directory.
		dataDirectoryString, err := dataDirectory.Get()
		if err != nil {
			log.Fatalf("failed getting data directory: %f\n", err)
		}
		view.preferences.SetString(constants.PreferenceKeyDataDirectory, dataDirectoryString)
		view.preferences.SetBool(constants.PreferenceKeyHasSetDataDirectory, true)

		view.nextPageID <- constants.PageIDStartupView
	})
	bottom := container.NewHBox(layout.NewSpacer(), okBtn, cancelBtn)

	content := container.NewVBox(lb1, lb2, lb3, radio, entryBox, lb4, lb5, bottom)
	return content
}

func (view *PickDataDirectoryView) WaitUntilReadyToTransition() int {
	return <-view.nextPageID
}
