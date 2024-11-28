package service

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"strings"
	"time"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/common/httperror"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/usecase"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BlockchainSyncManagerService struct {
	logger                                                               *slog.Logger
	blockchainSyncWithBlockchainAuthorityService                         *BlockchainSyncWithBlockchainAuthorityService
	subscribeToBlockchainStateChangeEventsFromBlockchainAuthorityUseCase *usecase.SubscribeToBlockchainStateChangeEventsFromBlockchainAuthorityUseCase
}

func NewBlockchainSyncManagerService(
	logger *slog.Logger,
	s1 *BlockchainSyncWithBlockchainAuthorityService,
	uc1 *usecase.SubscribeToBlockchainStateChangeEventsFromBlockchainAuthorityUseCase,
) *BlockchainSyncManagerService {
	return &BlockchainSyncManagerService{logger, s1, uc1}
}

func (s *BlockchainSyncManagerService) Execute(ctx context.Context, chainID uint16, tenantID primitive.ObjectID) error {
	//
	// STEP 1: Validation.
	//

	e := make(map[string]string)
	if chainID == 0 {
		e["chain_id"] = "missing value"
	}
	if tenantID.IsZero() {
		e["tenant_id"] = "missing value"
	}
	if len(e) != 0 {
		s.logger.Warn("Validation failed.",
			slog.Any("error", e))
		return httperror.NewForBadRequest(&e)
	}

	//
	// STEP 2:
	// On startup sync with Blockchain Authority.
	//

	s.logger.Debug("Syncing...")

	if err := s.blockchainSyncWithBlockchainAuthorityService.Execute(ctx, chainID, tenantID); err != nil {
		log.Fatalf("Failed to sync blockchain: %v\n", err)
	}

	s.logger.Debug("Syncing finished")

	//
	// STEP 3:
	// Once startup sync has been completed, subscribe to the `server sent
	// events` of the Blockchain Authority to get the latest updates about
	// changes with the global blockchain network.
	//

	s.logger.Debug("Waiting to receive from the global blockchain network...",
		slog.Any("chain_id", chainID))

	// Subscribe to the Blockchain Authority to receive `server sent events`
	// when the blockchain changes globally to our local machine.
	ch, err := s.subscribeToBlockchainStateChangeEventsFromBlockchainAuthorityUseCase.Execute(ctx, chainID)
	if err != nil {
		if strings.Contains(err.Error(), "received non-OK HTTP status: 524") {
			s.logger.Warn("Failed subscribing because of timeout, will retry again in 10 seconds...",
				slog.Any("chainID", chainID),
				slog.Any("error", err))
			time.Sleep(10 * time.Second)
			return nil
		}

		s.logger.Error("Failed subscribing...",
			slog.Any("chainID", chainID),
			slog.Any("error", err))
		return err
	}

	// Consume data from the channel
	for value := range ch {
		fmt.Printf("Received update from chain ID: %d\n", value)

		if err := s.blockchainSyncWithBlockchainAuthorityService.Execute(ctx, chainID, tenantID); err != nil {
			log.Fatalf("Failed to sync blockchain: %v\n", err)
		}

		// DEVELOPERS NOTE:
		// Before we finish this runtime loop, and for debugging purposes, let
		// us print this helpful message to communicate to the user that we
		// are waiting for the next request.
		s.logger.Debug("Waiting to receive from the global blockchain network...",
			slog.Any("chain_id", chainID))
	}

	s.logger.Debug("Subscription to blockchain faucet closed",
		slog.Any("chain_id", chainID))

	return nil
}
