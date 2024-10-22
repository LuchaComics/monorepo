package wallet

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
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

type item struct {
	Name  string
	Email string
}

func (view *WalletViewOverviewTab) Render() *fyne.Container {

	// balanceLabel := widget.NewLabel("Total Balance: 100.00 CC")
	// tokensLabel := widget.NewLabel("Total Tokens: 500")

	// Create a list of items
	items := []item{
		{"John Doe", "john@example.com"},
		{"Jane Doe", "jane@example.com"},
		{"Bob Smith", "bob@example.com"},
	}

	// Create a new list
	list := widget.NewList(
		func() int {
			return len(items)
		},
		func() fyne.CanvasObject {
			return container.NewHBox(
				widget.NewLabel(""),
				widget.NewLabel(""),
			)
		},
		func(i widget.ListItemID, obj fyne.CanvasObject) {
			item := items[i]
			hbox := obj.(*fyne.Container)
			hbox.Objects[0].(*widget.Label).SetText(item.Name)
			hbox.Objects[1].(*widget.Label).SetText(item.Email)
			// hbox.Resize(fyne.NewSize(300, 400))
		},
	)

	// Create a header
	header := widget.NewLabel("List of Items")

	// Create a border layout with the header at the top and the list below
	border := container.NewBorder(header, nil, nil, nil, list)

	// border.Resize(fyne.NewSize(300, 400))
	return border
}
