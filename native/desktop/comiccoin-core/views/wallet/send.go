package wallet

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

type WalletViewSendTab struct {
	window      fyne.Window
	preferences fyne.Preferences
}

func NewWalletViewSendTab(
	w fyne.Window,
	pref fyne.Preferences) *WalletViewSendTab {
	v := &WalletViewSendTab{
		window:      w,
		preferences: pref,
	}

	return v
}

func (view *WalletViewSendTab) Render() *fyne.Container {

	loadingLabel := widget.NewLabel("Send")
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
