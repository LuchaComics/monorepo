package wallet

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

type WalletViewOverviewTab struct {
	window      fyne.Window
	preferences fyne.Preferences
}

func NewWalletViewOverviewTab(
	w fyne.Window,
	pref fyne.Preferences) *WalletViewOverviewTab {
	v := &WalletViewOverviewTab{
		window:      w,
		preferences: pref,
	}

	return v
}

func (view *WalletViewOverviewTab) Render() *fyne.Container {

	loadingLabel := widget.NewLabel("Overview")
	loadingLabel.Alignment = fyne.TextAlignCenter
	loadingLabel.TextStyle = fyne.TextStyle{Bold: true}

	// progressBar := widget.NewProgressBar()
	// progressBar.SetValue(0) // Initialize progress bar value

	content := container.NewVBox(
		layout.NewSpacer(),
		loadingLabel,
		// widget.NewSeparator(),
		// progressBar,
		// layout.NewSpacer(),
	)

	return content
}
