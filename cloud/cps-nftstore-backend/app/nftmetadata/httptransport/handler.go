package httptransport

import (
	"log/slog"

	nftmetadata_c "github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/app/nftmetadata/controller"
)

// Handler Creates http request handler
type Handler struct {
	Logger     *slog.Logger
	Controller nftmetadata_c.NFTMetadataController
}

// NewHandler Constructor
func NewHandler(loggerp *slog.Logger, c nftmetadata_c.NFTMetadataController) *Handler {
	return &Handler{
		Logger:     loggerp,
		Controller: c,
	}
}
