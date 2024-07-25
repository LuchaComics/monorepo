package httptransport

import (
	receipt_c "github.com/LuchaComics/monorepo/cloud/cps-backend/app/receipt/controller"
	"log/slog"
)

// Handler Creates http request handler
type Handler struct {
	Logger     *slog.Logger
	Controller receipt_c.ReceiptController
}

// NewHandler Constructor
func NewHandler(loggerp *slog.Logger, c receipt_c.ReceiptController) *Handler {
	return &Handler{
		Logger:     loggerp,
		Controller: c,
	}
}
