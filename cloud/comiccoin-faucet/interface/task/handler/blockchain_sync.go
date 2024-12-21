package handler

import (
	"context"
	"log/slog"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/config"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/service"
	"go.mongodb.org/mongo-driver/mongo"
)

type BlockchainSyncWithBlockchainAuthorityTaskHandler struct {
	config                                       *config.Configuration
	logger                                       *slog.Logger
	dbClient                                     *mongo.Client
	blockchainSyncWithBlockchainAuthorityService *service.BlockchainSyncWithBlockchainAuthorityService
}

func NewBlockchainSyncWithBlockchainAuthorityTaskHandler(
	config *config.Configuration,
	logger *slog.Logger,
	dbClient *mongo.Client,
	s1 *service.BlockchainSyncWithBlockchainAuthorityService,
) *BlockchainSyncWithBlockchainAuthorityTaskHandler {
	return &BlockchainSyncWithBlockchainAuthorityTaskHandler{config, logger, dbClient, s1}
}

func (h *BlockchainSyncWithBlockchainAuthorityTaskHandler) Execute(ctx context.Context) error {
	session, err := h.dbClient.StartSession()
	if err != nil {
		h.logger.Error("start session error",
			slog.Any("error", err))
		return err
	}
	defer session.EndSession(ctx)

	// Define a transaction function with a series of operations
	transactionFunc := func(sessCtx mongo.SessionContext) (interface{}, error) {
		return nil, h.blockchainSyncWithBlockchainAuthorityService.Execute(sessCtx)
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
