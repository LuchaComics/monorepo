package httptransport

import (
	"log/slog"

	ledger_c "github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/app/ledger/controller"
)

// Handler Creates http request handler
type Handler struct {
	Logger     *slog.Logger
	Controller ledger_c.LedgerController
}

// NewHandler Constructor
func NewHandler(loggerp *slog.Logger, c ledger_c.LedgerController) *Handler {
	return &Handler{
		Logger:     loggerp,
		Controller: c,
	}
}
