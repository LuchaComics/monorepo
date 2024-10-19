package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"

	//"fyne.io/fyne/v2/layout"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin-core/mvc"
)

func main() {
	a := app.New()
	w := a.NewWindow("ComicCoin core")
	w.Resize(fyne.NewSize(680, 480))

	v1 := mvc.NewPickBlockchainStorageLocationViewController(w)

	w.SetContent(v1)
	w.ShowAndRun()
}
