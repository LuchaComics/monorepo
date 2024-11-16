package service

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/common/httperror"

	auth_usecase "github.com/LuchaComics/monorepo/cloud/comiccoin-authority/usecase"
)

type BlockchainSyncManagerService struct {
	logger                                                               *slog.Logger
	subscribeToBlockchainStateChangeEventsFromBlockchainAuthorityUseCase *auth_usecase.SubscribeToBlockchainStateChangeEventsFromBlockchainAuthorityUseCase
}

func NewBlockchainSyncManagerService(
	logger *slog.Logger,
	uc1 *auth_usecase.SubscribeToBlockchainStateChangeEventsFromBlockchainAuthorityUseCase,
) *BlockchainSyncManagerService {
	return &BlockchainSyncManagerService{logger, uc1}
}

func (s *BlockchainSyncManagerService) Execute(ctx context.Context, chainID uint16) error {
	//
	// STEP 1: Validation.
	//

	e := make(map[string]string)
	if chainID == 0 {
		e["chain_id"] = "missing value"
	}
	if len(e) != 0 {
		s.logger.Warn("Validation failed.",
			slog.Any("error", e))
		return httperror.NewForBadRequest(&e)
	}

	//
	// STEP 2: Get our account from our in-memory database if it exists.
	//

	// Subscribe to the Blockchain Authority to receive `server sent events`
	// when the blockchain changes globally to our local machine.
	ch, err := s.subscribeToBlockchainStateChangeEventsFromBlockchainAuthorityUseCase.Execute(ctx, chainID)
	if err != nil {
		s.logger.Error("failed subscribing",
			slog.Any("chainID", chainID),
			slog.Any("error", err))
		return err
	}

	// Consume data from the channel
	for value := range ch {
		fmt.Printf("Received update: %d\n", value)
	}
	fmt.Println("Subscription closed")

	return nil
}
