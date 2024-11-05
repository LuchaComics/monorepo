package handler

import (
	"log/slog"
	"net/http"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/config"
)

type IPFSGatewayHTTPHandler struct {
	config *config.Config
	logger *slog.Logger
}

func NewIPFSGatewayHTTPHandler(
	cfg *config.Config,
	logger *slog.Logger,
) *IPFSGatewayHTTPHandler {
	return &IPFSGatewayHTTPHandler{cfg, logger}
}

func (h *IPFSGatewayHTTPHandler) Execute(w http.ResponseWriter, r *http.Request, tokenIDStr string) {

}
