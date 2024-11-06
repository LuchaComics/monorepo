package handler

import (
	"log/slog"
	"net/http"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/config"
)

type GetVersionHTTPHandler struct {
	config *config.Config
	logger *slog.Logger
}

func NewGetVersionHTTPHandler(
	cfg *config.Config,
	logger *slog.Logger,
) *GetVersionHTTPHandler {
	return &GetVersionHTTPHandler{cfg, logger}
}

func (h *GetVersionHTTPHandler) Execute(w http.ResponseWriter, r *http.Request) {

}
