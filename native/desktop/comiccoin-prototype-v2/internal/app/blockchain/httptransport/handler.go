package httptransport

import (
	"log/slog"

	blockchain_c "github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/app/blockchain/controller"
)

// Handler Creates http request handler
type Handler struct {
	Logger     *slog.Logger
	Controller blockchain_c.BlockchainController
}

// NewHandler Constructor
func NewHandler(loggerp *slog.Logger, c blockchain_c.BlockchainController) *Handler {
	return &Handler{
		Logger:     loggerp,
		Controller: c,
	}
}
