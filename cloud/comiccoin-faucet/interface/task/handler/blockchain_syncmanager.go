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
	session, err := h.dbClient.StartSession()
	if err != nil {
		h.logger.Error("start session error",
			slog.Any("error", err))
		return err
	}
	defer session.EndSession(ctx)

	// Define a transaction function with a series of operations
	transactionFunc := func(sessCtx mongo.SessionContext) (interface{}, error) {
		return nil, h.blockchainSyncManagerService.Execute(sessCtx, h.config.Blockchain.ChainID, h.config.App.TenantID)
	}

	// Start a transaction
	_, txErr := session.WithTransaction(ctx, transactionFunc)
	if txErr != nil {
		h.logger.Error("session failed error",
			slog.Any("error", txErr))
		return txErr
	}
	return nil
}
