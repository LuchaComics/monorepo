package httptransport

import (
	"log/slog"

	nftasset_c "github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/app/nftasset/controller"
)

// Handler Creates http request handler
type Handler struct {
	Logger     *slog.Logger
	Controller nftasset_c.NFTAssetController
}

// NewHandler Constructor
func NewHandler(loggerp *slog.Logger, c nftasset_c.NFTAssetController) *Handler {
	return &Handler{
		Logger:     loggerp,
		Controller: c,
	}
}

func (h *Handler) Shutdown() {
	h.Controller.Shutdown()
}
