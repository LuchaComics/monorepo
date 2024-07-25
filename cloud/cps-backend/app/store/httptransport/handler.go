package httptransport

import (
	"log/slog"

	store_c "github.com/LuchaComics/monorepo/cloud/cps-backend/app/store/controller"
)

// Handler Creates http request handler
type Handler struct {
	Logger     *slog.Logger
	Controller store_c.StoreController
}

// NewHandler Constructor
func NewHandler(loggerp *slog.Logger, c store_c.StoreController) *Handler {
	return &Handler{
		Logger:     loggerp,
		Controller: c,
	}
}
