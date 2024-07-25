package httptransport

import (
	"log/slog"

	customer_c "github.com/LuchaComics/monorepo/cloud/cps-backend/app/customer/controller"
)

// Handler Creates http request handler
type Handler struct {
	Logger     *slog.Logger
	Controller customer_c.CustomerController
}

// NewHandler Constructor
func NewHandler(loggerp *slog.Logger, c customer_c.CustomerController) *Handler {
	return &Handler{
		Logger:     loggerp,
		Controller: c,
	}
}
