package service

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/common/httperror"
	authority_domain "github.com/LuchaComics/monorepo/cloud/comiccoin-authority/domain"
	authority_usecase "github.com/LuchaComics/monorepo/cloud/comiccoin-authority/usecase"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/usecase"
)

type BlockchainStateSyncService struct {
	logger                                              *slog.Logger
	getBlockchainStateDTOFromBlockchainAuthorityUseCase *authority_usecase.GetBlockchainStateDTOFromBlockchainAuthorityUseCase
	upsertBlockchainStateUseCase                        *usecase.UpsertBlockchainStateUseCase
}

func NewBlockchainStateSyncService(
	logger *slog.Logger,
	uc1 *authority_usecase.GetBlockchainStateDTOFromBlockchainAuthorityUseCase,
	uc2 *usecase.UpsertBlockchainStateUseCase,
) *BlockchainStateSyncService {
	return &BlockchainStateSyncService{logger, uc1, uc2}
}

// Execute method gets genesis block data from authority if we don't have it locally, else gets genesis from local source.
func (s *BlockchainStateSyncService) Execute(ctx context.Context, chainID uint16) (*authority_domain.BlockchainState, error) {
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
		return nil, httperror.NewForBadRequest(&e)
	}

	//
	// STEP 2: Download from remote service.
	//

	blockchainStateDTO, err := s.getBlockchainStateDTOFromBlockchainAuthorityUseCase.Execute(ctx, chainID)
	if err != nil {
		s.logger.Error("Failed getting blockchain state remotely",
			slog.Any("chain_id", chainID),
			slog.Any("error", err))
		return nil, err
	}
	if blockchainStateDTO == nil {
		err := fmt.Errorf("Blockchain state does not exist for `chain_id`: %v", chainID)
		s.logger.Error("Failed getting blockchain state remotely",
			slog.Any("chain_id", chainID),
			slog.Any("error", err))
		return nil, err
	}

	// Convert from network format data to our local format.
	blockchainState := authority_domain.BlockchainStateDTOToBlockchainState(blockchainStateDTO)

	// Save the genesis block data to local database.
	if err := s.upsertBlockchainStateUseCase.Execute(ctx, blockchainState); err != nil {
		s.logger.Error("Failed upserting blockchain state",
			slog.Any("chain_id", chainID),
			slog.Any("error", err))
		return nil, err
	}

	return blockchainState, nil
}
