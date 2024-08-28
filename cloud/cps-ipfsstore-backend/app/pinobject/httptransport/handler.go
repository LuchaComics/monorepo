package httptransport

import (
	"log/slog"

	pinobject_c "github.com/LuchaComics/monorepo/cloud/cps-ipfsstore-backend/app/pinobject/controller"
)

// Handler Creates http request handler
type Handler struct {
	Logger     *slog.Logger
	Controller pinobject_c.PinObjectController
}

// NewHandler Constructor
func NewHandler(loggerp *slog.Logger, c pinobject_c.PinObjectController) *Handler {
	return &Handler{
		Logger:     loggerp,
		Controller: c,
	}
}

func (h *Handler) Shutdown() {
	h.Controller.Shutdown()
}
