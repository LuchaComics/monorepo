package views

import (
	"fyne.io/fyne/v2"
)

type Viewer interface {
	// Return the container necessary to render all the GUI elements required
	// for this particular view.
	Render() fyne.CanvasObject

	// Wait (block execution flow) until the view is ready to transition to a
	// different view and return the `ID` of what that view should be.
	WaitUntilReadyToTransition() int
}

type TabViewer interface {
	// Return the container necessary to render all the GUI elements required
	// for this particular tab view.
	Render() fyne.CanvasObject
}
