package handler

import (
	"context"
	"log/slog"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/service"
)

type MempoolReceiveTaskHandler struct {
	config  *config.Config
	logger  *slog.Logger
	service *service.MempoolReceiveService
}

func NewMempoolReceiveTaskHandler(
	cfg *config.Config,
	logger *slog.Logger,
	mempoolReceiveService *service.MempoolReceiveService,
) *MempoolReceiveTaskHandler {
	return &MempoolReceiveTaskHandler{cfg, logger, mempoolReceiveService}
}

type BlockchainMempoolReceiveResponseIDO struct {
}

func (h *MempoolReceiveTaskHandler) Execute(ctx context.Context) error {
	serviceExecErr := h.service.Execute(ctx)
	if serviceExecErr != nil {
		return serviceExecErr
	}
	return nil
}
