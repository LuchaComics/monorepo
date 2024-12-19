package handler

import (
	"context"
	"log/slog"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/config"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/service"
	"go.mongodb.org/mongo-driver/mongo"
)

type BlockchainSyncManagerTaskHandler struct {
	config                       *config.Configuration
	logger                       *slog.Logger
	dbClient                     *mongo.Client
	blockchainSyncManagerService *service.BlockchainSyncManagerService
}

func NewBlockchainSyncManagerTaskHandler(
	config *config.Configuration,
	logger *slog.Logger,
	dbClient *mongo.Client,
	s1 *service.BlockchainSyncManagerService,
) *BlockchainSyncManagerTaskHandler {
	return &BlockchainSyncManagerTaskHandler{config, logger, dbClient, s1}
}

func (h *BlockchainSyncManagerTaskHandler) Execute(ctx context.Context) error {
	// DEVELOPERS NOTE:
	// Do not use MongoDB transactions here, this service handles
	// them internally. If you do then you will get errors such as:
	// - "MongoServerError: WriteConflict error: this operation conflicted with another operation. Please retry your operation or multi-document transaction"
	return h.blockchainSyncManagerService.Execute(ctx, h.config.Blockchain.ChainID, h.config.App.TenantID)
}
