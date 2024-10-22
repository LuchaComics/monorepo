package wallet

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type WalletViewTransactionsTab struct {
	window      fyne.Window
	preferences fyne.Preferences
}

func NewWalletViewTransactionsTab(
	w fyne.Window,
	pref fyne.Preferences) *WalletViewTransactionsTab {
	v := &WalletViewTransactionsTab{
		window:      w,
		preferences: pref,
	}

	return v
}

func (view *WalletViewTransactionsTab) Render() fyne.CanvasObject {

	type item struct {
		Name  string
		Email string
	}

	// balanceLabel := widget.NewLabel("Total Balance: 100.00 CC")
	// tokensLabel := widget.NewLabel("Total Tokens: 500")

	// Create a list of items
	items := []item{
		{"John Doe", "john@example.com"},
		{"Jane Doe", "jane@example.com"},
		{"Bob Smith", "bob@example.com"},
		{"Bob Smith", "bob@example.com"},
		{"Bob Smith", "bob@example.com"},
		{"Bob Smith", "bob@example.com"},
		{"Bob Smith", "bob@example.com"},
		{"Bob Smith", "bob@example.com"},
		{"Bob Smith", "bob@example.com"},
		{"Bob Smith", "bob@example.com"},
		{"Bob Smith", "bob@example.com"},
		{"Bob Smith", "bob@example.com"},
		{"Bob Smith", "bob@example.com"},
		{"Bob Smith", "bob@example.com"},
		{"Bob Smith", "bob@example.com"},
		{"Bob Smith", "bob@example.com"},
		{"Bob Smith", "bob@example.com"},
		{"Bob Smith", "bob@example.com"},
		{"Bob Smith", "bob@example.com"},
		{"Bob Smith", "bob@example.com"},
		{"Bob Smith", "bob@example.com"},
		{"Bob Smith", "bob@example.com"},
		{"Bob Smith", "bob@example.com"},
		{"Bob Smith", "bob@example.com"},
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

	xxx := container.NewScroll(list)

	// Create a new container with the label and stretch it vertically
	container := container.NewBorder(header, nil, nil, nil, xxx)

	view.window.SetContent(container)
	// view.window.Resize(fyne.NewSize(400, 600))

	// border.Resize(fyne.NewSize(300, 400))
	return container
}
