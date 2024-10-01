package handler

import (
	"context"
	"log/slog"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/service"
)

type MempoolBatchSendTaskHandler struct {
	config  *config.Config
	logger  *slog.Logger
	service *service.MempoolBatchSendService
}

func NewMempoolBatchSendTaskHandler(
	cfg *config.Config,
	logger *slog.Logger,
	mempoolBatchSendService *service.MempoolBatchSendService,
) *MempoolBatchSendTaskHandler {
	return &MempoolBatchSendTaskHandler{cfg, logger, mempoolBatchSendService}
}

func (h *MempoolBatchSendTaskHandler) Execute(ctx context.Context) error {
	serviceExecErr := h.service.Execute(ctx)
	if serviceExecErr != nil {
		return serviceExecErr
	}
	return nil
}
