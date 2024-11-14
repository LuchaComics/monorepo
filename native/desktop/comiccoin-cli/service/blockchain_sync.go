package service

import (
	"context"
	"log/slog"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/common/httperror"
)

type BlockchainSyncWithBlockchainAuthorityService struct {
	logger                     *slog.Logger
	genesisBlockDataGetOrSync  *GenesisBlockDataGetOrSyncService
	blockchainStateSyncService *BlockchainStateSyncService
}

func NewBlockchainSyncWithBlockchainAuthorityService(
	logger *slog.Logger,
	s1 *GenesisBlockDataGetOrSyncService,
	s2 *BlockchainStateSyncService,
) *BlockchainSyncWithBlockchainAuthorityService {
	return &BlockchainSyncWithBlockchainAuthorityService{logger, s1, s2}
}

func (s *BlockchainSyncWithBlockchainAuthorityService) Execute(ctx context.Context, chainID uint16) error {
	//
	// STEP 1: Validation.
	//

	e := make(map[string]string)
	if chainID == 0 {
		e["chainID"] = "missing value"
	}
	if len(e) != 0 {
		s.logger.Warn("Validation failed",
			slog.Any("error", e))
		return httperror.NewForBadRequest(&e)
	}

	//
	// STEP 2:
	// Get genesis block data from authority if we don't have it locally, else
	// get it locally.
	//

	genesis, err := s.genesisBlockDataGetOrSync.Execute(ctx, chainID)
	if err != nil {
		s.logger.Error("Failed getting genesis block",
			slog.Any("chain_id", chainID),
			slog.Any("error", err))
		return err
	}
	_ = genesis

	//
	// STEP 3:
	// Refresh our local blockchain state with what exists currently on the
	// blockchain network.
	//

	latestBlockchainState, err := s.blockchainStateSyncService.Execute(ctx, chainID)
	if err != nil {
		s.logger.Error("Failed syncing blockchain state",
			slog.Any("chain_id", chainID),
			slog.Any("error", err))
		return err
	}
	_ = latestBlockchainState

	return nil
}
