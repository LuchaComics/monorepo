package wallet

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

type WalletViewReceiveTab struct {
	window      fyne.Window
	preferences fyne.Preferences
}

func NewWalletViewReceiveTab(
	w fyne.Window,
	pref fyne.Preferences) *WalletViewReceiveTab {
	v := &WalletViewReceiveTab{
		window:      w,
		preferences: pref,
	}

	return v
}

func (view *WalletViewReceiveTab) Render() fyne.CanvasObject {

	loadingLabel := widget.NewLabel("Receive")
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
