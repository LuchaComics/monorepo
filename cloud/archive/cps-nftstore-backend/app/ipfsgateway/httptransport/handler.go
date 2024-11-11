package httptransport

import (
	"log/slog"

	pinobject_c "github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/app/ipfsgateway/controller"
)

// Handler Creates http request handler
type Handler struct {
	Logger     *slog.Logger
	Controller pinobject_c.IpfsGatewayController
}

// NewHandler Constructor
func NewHandler(loggerp *slog.Logger, c pinobject_c.IpfsGatewayController) *Handler {
	return &Handler{
		Logger:     loggerp,
		Controller: c,
	}
}

func (h *Handler) Shutdown() {
	h.Controller.Shutdown()
}
