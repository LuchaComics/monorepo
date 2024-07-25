package httptransport

import (
	userpurchase_c "github.com/LuchaComics/monorepo/cloud/cps-backend/app/userpurchase/controller"
	"log/slog"
)

// Handler Creates http request handler
type Handler struct {
	Logger     *slog.Logger
	Controller userpurchase_c.UserPurchaseController
}

// NewHandler Constructor
func NewHandler(loggerp *slog.Logger, c userpurchase_c.UserPurchaseController) *Handler {
	return &Handler{
		Logger:     loggerp,
		Controller: c,
	}
}
