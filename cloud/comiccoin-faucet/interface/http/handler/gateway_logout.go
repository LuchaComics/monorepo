package handler

import (
	"log/slog"
	"net/http"
	_ "time/tzdata"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/common/httperror"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/service"
)

type GatewayLogoutHTTPHandler struct {
	logger   *slog.Logger
	dbClient *mongo.Client
	service  *service.GatewayLogoutService
}

func NewGatewayLogoutHTTPHandler(
	logger *slog.Logger,
	dbClient *mongo.Client,
	service *service.GatewayLogoutService,
) *GatewayLogoutHTTPHandler {
	return &GatewayLogoutHTTPHandler{
		logger:   logger,
		dbClient: dbClient,
		service:  service,
	}
}

func (h *GatewayLogoutHTTPHandler) Execute(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if err := h.service.Execute(ctx); err != nil {
		httperror.ResponseError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)

}
