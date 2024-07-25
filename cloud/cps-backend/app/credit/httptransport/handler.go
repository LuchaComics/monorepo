package httptransport

import (
	credit_c "github.com/LuchaComics/monorepo/cloud/cps-backend/app/credit/controller"
	"log/slog"
)

// Handler Creates http request handler
type Handler struct {
	Logger     *slog.Logger
	Controller credit_c.CreditController
}

// NewHandler Constructor
func NewHandler(loggerp *slog.Logger, c credit_c.CreditController) *Handler {
	return &Handler{
		Logger:     loggerp,
		Controller: c,
	}
}
