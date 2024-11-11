package httptransport

import (
	"log/slog"

	nft_c "github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/app/nft/controller"
)

// Handler Creates http request handler
type Handler struct {
	Logger     *slog.Logger
	Controller nft_c.NFTController
}

// NewHandler Constructor
func NewHandler(loggerp *slog.Logger, c nft_c.NFTController) *Handler {
	return &Handler{
		Logger:     loggerp,
		Controller: c,
	}
}
