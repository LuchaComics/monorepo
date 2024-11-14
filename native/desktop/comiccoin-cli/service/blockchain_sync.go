package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/common/httperror"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/domain"
	authority_domain "github.com/LuchaComics/monorepo/cloud/comiccoin-authority/domain"
	auth_usecase "github.com/LuchaComics/monorepo/cloud/comiccoin-authority/usecase"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin-cli/usecase"
)

type BlockchainSyncWithBlockchainAuthorityService struct {
	logger                                               *slog.Logger
	getGenesisBlockDataUseCase                           *usecase.GetGenesisBlockDataUseCase
	upsertGenesisBlockDataUseCase                        *usecase.UpsertGenesisBlockDataUseCase
	getGenesisBlockDataDTOFromBlockchainAuthorityUseCase *auth_usecase.GetGenesisBlockDataDTOFromBlockchainAuthorityUseCase
	getBlockchainStateUseCase                            *usecase.GetBlockchainStateUseCase
	upsertBlockchainStateUseCase                         *usecase.UpsertBlockchainStateUseCase
	getBlockchainStateDTOFromBlockchainAuthorityUseCase  *auth_usecase.GetBlockchainStateDTOFromBlockchainAuthorityUseCase
	getBlockDataUseCase                                  *usecase.GetBlockDataUseCase
	upsertBlockDataUseCase                               *usecase.UpsertBlockDataUseCase
	getBlockDataDTOFromBlockchainAuthorityUseCase        *auth_usecase.GetBlockDataDTOFromBlockchainAuthorityUseCase
}

func NewBlockchainSyncWithBlockchainAuthorityService(
	logger *slog.Logger,
	uc1 *usecase.GetGenesisBlockDataUseCase,
	uc2 *usecase.UpsertGenesisBlockDataUseCase,
	uc3 *auth_usecase.GetGenesisBlockDataDTOFromBlockchainAuthorityUseCase,
	uc4 *usecase.GetBlockchainStateUseCase,
	uc5 *usecase.UpsertBlockchainStateUseCase,
	uc6 *auth_usecase.GetBlockchainStateDTOFromBlockchainAuthorityUseCase,
	uc7 *usecase.GetBlockDataUseCase,
	uc8 *usecase.UpsertBlockDataUseCase,
	uc9 *auth_usecase.GetBlockDataDTOFromBlockchainAuthorityUseCase,
) *BlockchainSyncWithBlockchainAuthorityService {
	return &BlockchainSyncWithBlockchainAuthorityService{logger, uc1, uc2, uc3, uc4, uc5, uc6, uc7, uc8, uc9}
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
	// Step 2:
	// Get our genesis block, and if it doesn't exist then we need to fetch it
	// from the blockchain authority for the particular `chainID`.
	//

	genesis, err := s.getGenesisBlockDataUseCase.Execute(ctx, chainID)
	if err != nil {
		s.logger.Error("Failed getting genesis block locally",
			slog.Any("chain_id", chainID),
			slog.Any("error", err))
		return err
	}
	if genesis == nil {
		s.logger.Debug("Genesis block d.n.e, fetching it now ...")
		genesisDTO, err := s.getGenesisBlockDataDTOFromBlockchainAuthorityUseCase.Execute(ctx, chainID)
		if err != nil {
			s.logger.Error("Failed getting genesis block remotely",
				slog.Any("chain_id", chainID),
				slog.Any("error", err))
			return err
		}
		if genesisDTO == nil {
			err := fmt.Errorf("Genesis block data does not exist for `chain_id`: %v", chainID)
			s.logger.Error("Failed getting genesis block remotely",
				slog.Any("chain_id", chainID),
				slog.Any("error", err))
			return err
		}

		// Convert from network format data to our local format.
		genesis = authority_domain.GenesisBlockDataDTOToGenesisBlockData(genesisDTO)

		// Save the genesis block data to local database.
		if err := s.upsertGenesisBlockDataUseCase.Execute(ctx, genesis.Hash, genesis.Header, genesis.HeaderSignatureBytes, genesis.Trans, genesis.Validator); err != nil {
			s.logger.Error("Failed upserting genesis",
				slog.Any("chain_id", chainID),
				slog.Any("error", err))
			return err
		}

		s.logger.Debug("Genesis block saved to local database from global blockchain network",
			slog.Any("chain_id", chainID))
	}

	//
	// STEP 3:
	// Get the blockchain state we have *locally* and *remotely* and compare
	// the differences, if our local blockchain state matches what is on the
	// global blockchain network then we are done synching (because there is
	// nothin left to sync). If we don't even have a blockchain state then we need to
	// proceed to download the entire blockchain immediately. If there is any
	// discrepency between the global and local state then we proceed with
	// this function and update our local blockchain with the available data
	// on the global blockchain network.
	//

	globalBlockchainStateDTO, err := s.getBlockchainStateDTOFromBlockchainAuthorityUseCase.Execute(ctx, chainID)
	if err != nil {
		s.logger.Error("Failed getting global blockchain state",
			slog.Any("chain_id", chainID),
			slog.Any("error", err))
		return err
	}
	if globalBlockchainStateDTO == nil {
		err := fmt.Errorf("Failed getting global blockchain state for chainID: %v", chainID)
		s.logger.Error("Failed getting global blockchain state",
			slog.Any("chain_id", chainID),
			slog.Any("error", err))
		return err
	}
	// Convert from network format data to our local format.
	globalBlockchainState := authority_domain.BlockchainStateDTOToBlockchainState(globalBlockchainStateDTO)

	// Fetch our local blockchain state.
	localBlockchainState, err := s.getBlockchainStateUseCase.Execute(ctx, chainID)
	if err != nil {
		s.logger.Error("Failed getting local blockchain state",
			slog.Any("chain_id", chainID),
			slog.Any("error", err))
		return err
	}

	// If our local blockchain state is empty then create it with using the genesis block.
	if localBlockchainState == nil {
		localBlockchainState = &domain.BlockchainState{
			ChainID:                genesis.Header.ChainID,
			LatestBlockNumberBytes: genesis.Header.NumberBytes,
			LatestHash:             genesis.Hash,
			LatestTokenIDBytes:     genesis.Header.LatestTokenIDBytes,
			AccountHashState:       genesis.Header.StateRoot,
			TokenHashState:         genesis.Header.TokensRoot,
		}
		if err := s.upsertBlockchainStateUseCase.Execute(ctx, localBlockchainState); err != nil {
			s.logger.Error("Failed upserting local blockchain state from genesis block data",
				slog.Any("chain_id", chainID),
				slog.Any("error", err))
			return err
		}
		s.logger.Debug("Local blockchain state set to genesis block data",
			slog.Any("chain_id", chainID))
	} else {
		if localBlockchainState.LatestHash == globalBlockchainState.LatestHash {
			s.logger.Debug("Local blockchain is in sync with global blockchain network",
				slog.Any("chain_id", chainID))
			return nil
		}
		s.logger.Debug("Local blockchain state is out of sync with global blockchain network",
			slog.Any("chain_id", chainID))
	}

	//
	// STEP 4:
	// Proceed to download all the missing block data from the global blockchain
	// network so our local blockchain will be in-sync.
	//

	if err := s.syncWithGlobalBlockchainNetwork(ctx, localBlockchainState, globalBlockchainState); err != nil {
		if localBlockchainState.LatestHash == globalBlockchainState.LatestHash {
			s.logger.Debug("Failed to sync with the global blockchain network",
				slog.Any("chain_id", chainID))
			return nil
		}
	}

	//
	// STEP 5:
	// Update our blockchain state to match the global blockchain network's state.
	//

	// TODO: IMPL.

	return nil
}

func (s *BlockchainSyncWithBlockchainAuthorityService) syncWithGlobalBlockchainNetwork(ctx context.Context, localBlockchainState, globalBlockchainState *domain.BlockchainState) error {
	s.logger.Debug("Beginning to sync with global blockchain network...")
	s.logger.Debug("Finished syncing with global blockchain network")
	return errors.New("HALT BY PROGRAMMER")
	return nil
}
