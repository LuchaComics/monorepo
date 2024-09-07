package httptransport

import (
	"log/slog"

	collection_c "github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/app/collection/controller"
)

// Handler Creates http request handler
type Handler struct {
	Logger     *slog.Logger
	Controller collection_c.CollectionController
}

// NewHandler Constructor
func NewHandler(loggerp *slog.Logger, c collection_c.CollectionController) *Handler {
	return &Handler{
		Logger:     loggerp,
		Controller: c,
	}
}
