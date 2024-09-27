package handler

import (
	"log/slog"
	"net/http"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/account/service"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/config"
)

type CreateAccountHTTPHandler struct {
	config               *config.Config
	logger               *slog.Logger
	createAccountService *service.CreateAccountService
}

func NewCreateAccountHTTPHandler(
	cfg *config.Config,
	logger *slog.Logger,
	s *service.CreateAccountService,
) *CreateAccountHTTPHandler {
	return &CreateAccountHTTPHandler{cfg, logger, s}
}

func (h *CreateAccountHTTPHandler) Execute(w http.ResponseWriter, r *http.Request) {

}
