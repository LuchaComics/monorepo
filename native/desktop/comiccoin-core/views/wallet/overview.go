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

func (view *WalletViewOverviewTab) Render() fyne.CanvasObject {
	balanceLabel := widget.NewLabel("Total Balance: 100.00 CC")
	tokensLabel := widget.NewLabel("Total Tokens: 500")

	type item struct {
		Name  string
		Email string
	}

	// Create a list of items
	items := []item{
		{"John Doe", "john@example.com"},
		{"Jane Doe", "jane@example.com"},
		{"Bob Smith", "bob@example.com"},
		{"Bob Smith", "bob@example.com"},
		{"John Doe", "john@example.com"},
		{"Bob Smith", "bob@example.com"},
		{"Bob Smith", "bob@example.com"},
		{"Bob Smith", "bob@example.com"},
		{"Bob Smith", "bob@example.com"},
		{"Bob Smith", "bob@example.com"},
		{"John Doe", "john@example.com"},
		{"Bob Smith", "bob@example.com"},
		{"Bob Smith", "bob@example.com"},
		{"Bob Smith", "bob@example.com"},
		{"Bob Smith", "bob@example.com"},
		{"Bob Smith", "bob@example.com"},
		{"John Doe", "john@example.com"},
		{"Bob Smith", "bob@example.com"},
		{"Bob Smith", "bob@example.com"},
		{"Bob Smith", "bob@example.com"},
		{"Bob Smith", "bob@example.com"},
		{"Bob Smith", "bob@example.com"},
		{"John Doe", "john@example.com"},
		{"Bob Smith", "bob@example.com"},
		{"Bob Smith", "bob@example.com"},
		{"Bob Smith", "bob@example.com"},
		{"Bob Smith", "bob@example.com"},
		{"Bob Smith", "bob@example.com"},
		{"Bob Smith", "bob@example.com"},
		// ...
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
		},
	)

	// Create a header
	header := widget.NewLabel("Recent Transactions")

	v := container.NewVBox(
		balanceLabel,
		tokensLabel,
		header,
		container.NewScroll(list),
	)

	return v
}
