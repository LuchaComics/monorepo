package httptransport

import (
	"log/slog"

	account_c "github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/app/account/controller"
)

// Handler Creates http request handler
type Handler struct {
	Logger     *slog.Logger
	Controller account_c.AccountController
}

// NewHandler Constructor
func NewHandler(loggerp *slog.Logger, c account_c.AccountController) *Handler {
	return &Handler{
		Logger:     loggerp,
		Controller: c,
	}
}
