package wallet

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin-core/constants"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin-core/views"
)

type WalletView struct {
	window      fyne.Window
	preferences fyne.Preferences

	overviewTab views.TabViewer
	sendTab     views.TabViewer
	receiveTab  views.TabViewer
	txsTab      views.TabViewer
	moreTab     views.TabViewer

	// nextPageID represents a recieve operation channel used by this view set the
	nextPageID chan int
}

func NewWalletView(
	w fyne.Window,
	pref fyne.Preferences,
	overviewTab views.TabViewer,
	sendTab views.TabViewer,
	receiveTab views.TabViewer,
	txsTab views.TabViewer,
	moreTab views.TabViewer,
) views.Viewer {
	v := &WalletView{
		window:      w,
		preferences: pref,
		nextPageID:  make(chan int),
		overviewTab: overviewTab,
		sendTab:     sendTab,
		receiveTab:  receiveTab,
		txsTab:      txsTab,
		moreTab:     moreTab,
	}

	return v
}

func (view *WalletView) Render() *fyne.Container {
	view.window.Resize(fyne.NewSize(constants.DefaultScreenWidth, constants.DefaultScreenHeight))

	tabs := container.NewAppTabs()
	tabs.Append(container.NewTabItemWithIcon("Overview", theme.HomeIcon(), view.overviewTab.Render()))
	tabs.Append(container.NewTabItemWithIcon("Send", theme.MailSendIcon(), view.sendTab.Render()))
	tabs.Append(container.NewTabItemWithIcon("Receive", theme.DownloadIcon(), view.receiveTab.Render()))
	tabs.Append(container.NewTabItemWithIcon("Transactions", theme.ListIcon(), view.txsTab.Render()))
	tabs.Append(container.NewTabItemWithIcon("", theme.MoreHorizontalIcon(), view.moreTab.Render()))

	tabs.SetTabLocation(container.TabLocationTop)

	return container.New(layout.NewVBoxLayout(), tabs)
}

func (view *WalletView) WaitUntilReadyToTransition() int {
	return <-view.nextPageID
}
