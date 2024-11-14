package service

import (
	"context"
	"log/slog"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/common/httperror"
)

type BlockchainSyncWithBlockchainAuthorityService struct {
	logger                     *slog.Logger
	genesisBlockDataGetService *GenesisBlockDataGetService
}

func NewBlockchainSyncWithBlockchainAuthorityService(
	logger *slog.Logger,
	s1 *GenesisBlockDataGetService,
) *BlockchainSyncWithBlockchainAuthorityService {
	return &BlockchainSyncWithBlockchainAuthorityService{logger, s1}
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

	genesis, err := s.genesisBlockDataGetService.Execute(ctx, chainID)
	if err != nil {
		s.logger.Error("Failed getting genesis block",
			slog.Any("chain_id", chainID),
			slog.Any("error", err))
		return err
	}
	_ = genesis

	return nil
}
