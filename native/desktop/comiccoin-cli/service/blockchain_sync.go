package service

import (
	"context"
	"errors"
	"log/slog"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin-cli/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin-cli/domain"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin-cli/usecase"
)

type LocalBlockchainSyncWithCentralAuthorityService struct {
	config                                                  *config.Config
	logger                                                  *slog.Logger
	getBlockchainStateUseCase                               *usecase.GetBlockchainStateUseCase
	getBlockchainStateFromCentralAuthorityByChainIDUseCase  *usecase.GetBlockchainStateFromCentralAuthorityByChainIDUseCase
	upsertBlockchainStateUseCase                            *usecase.UpsertBlockchainStateUseCase
	getGenesisBlockDataUseCase                              *usecase.GetGenesisBlockDataUseCase
	getGenesisBlockDataFromCentralAuthorityByChainIDUseCase *usecase.GetGenesisBlockDataFromCentralAuthorityByChainIDUseCase
	upsertGenesisBlockDataUseCase                           *usecase.UpsertGenesisBlockDataUseCase
	getBlockDataFromCentralAuthorityByBlockNumberUseCase    *usecase.GetBlockDataFromCentralAuthorityByBlockNumberUseCase
}

func NewLocalBlockchainSyncWithCentralAuthorityService(
	config *config.Config,
	logger *slog.Logger,
	uc1 *usecase.GetBlockchainStateUseCase,
	uc2 *usecase.GetBlockchainStateFromCentralAuthorityByChainIDUseCase,
	uc3 *usecase.UpsertBlockchainStateUseCase,
	uc4 *usecase.GetGenesisBlockDataUseCase,
	uc5 *usecase.GetGenesisBlockDataFromCentralAuthorityByChainIDUseCase,
	uc6 *usecase.UpsertGenesisBlockDataUseCase,
	uc7 *usecase.GetBlockDataFromCentralAuthorityByBlockNumberUseCase,
) *LocalBlockchainSyncWithCentralAuthorityService {
	return &LocalBlockchainSyncWithCentralAuthorityService{config, logger, uc1, uc2, uc3, uc4, uc5, uc6, uc7}
}

func (s *LocalBlockchainSyncWithCentralAuthorityService) Execute(ctx context.Context) error {
	//
	// STEP 1.
	// Get blockchain state from the authority and get the blockchain and get
	// the blockchain state locally. If the values match then we can skip
	// executing this service.
	//

	blockchainStateFromAuthority, err := s.getBlockchainStateFromCentralAuthorityByChainIDUseCase.Execute(ctx, s.config.Blockchain.ChainID)
	if err != nil {
		s.logger.Error("Failed getting from the authority",
			slog.Any("error", err))
		return err
	}
	if blockchainStateFromAuthority == nil {
		dneErr := errors.New("Failed fetching from the authority with no results")
		s.logger.Error("Failed getting from the authority",
			slog.Any("error", dneErr))
		return dneErr
	}

	s.logger.Debug("Fetched latest blockchain state",
		slog.Any("blockchain_state_from_authority", blockchainStateFromAuthority))

	localblockchainState, err := s.getBlockchainStateUseCase.Execute(ctx, s.config.Blockchain.ChainID)
	if err != nil {
		s.logger.Error("Failed getting from local",
			slog.Any("error", err))
		return err
	}
	if localblockchainState != nil {
		if blockchainStateFromAuthority.LatestBlockNumber == localblockchainState.LatestBlockNumber {
			s.logger.Debug("Local blockchain already in sync with the authority.")
			return nil
		} else {
			s.logger.Debug("Local blockchain is out of sync with authority, beginning to update now...")
		}
	} else {
		s.logger.Debug("Local blockchain in empty, beginning to download blockchain from the authority now...")

		//
		// STEP 2:
		// If empty blockchain, start by downloading the genesis block data first.
		//

		genesisBlockDTO, err := s.getGenesisBlockDataFromCentralAuthorityByChainIDUseCase.Execute(ctx, s.config.Blockchain.ChainID)
		if err != nil {
			s.logger.Error("Failed getting genesis block from the authority",
				slog.Any("error", err))
			return err
		}
		if genesisBlockDTO == nil {
			dneErr := errors.New("Failed fetching genesis block from the authority with no results")
			s.logger.Error("Failed getting genesis block from the authority",
				slog.Any("error", dneErr))
			return dneErr
		}

		genesisBlockIDO := genesisBlockDTO.ToIDO()
		if err := s.upsertGenesisBlockDataUseCase.Execute(ctx, genesisBlockIDO); err != nil {
			s.logger.Error("Failed saving genesis block to local",
				slog.Any("error", err))
			return err
		}
		s.logger.Debug("Genesis block saved to local.")

		//
		// STEP 3:
		// Set latest blockchain state to point to the genesis block data.
		//

		localblockchainState = &domain.BlockchainState{
			ChainID:           genesisBlockIDO.Header.ChainID,
			LatestBlockNumber: genesisBlockIDO.Header.Number,
			LatestHash:        genesisBlockIDO.Hash,
			LatestTokenID:     genesisBlockIDO.Header.LatestTokenID,
			AccountHashState:  genesisBlockIDO.Header.StateRoot,
			TokenHashState:    genesisBlockIDO.Header.TokensRoot,
		}

		if err := s.upsertBlockchainStateUseCase.Execute(ctx, localblockchainState); err != nil {
			s.logger.Error("Failed saving blockchain state to local",
				slog.Any("error", err))
			return err
		}
		s.logger.Debug("Initial blockchain state saved to local.")
	}

	//
	// STEP 4:
	// Download all the missing block data from the authority.
	//

	startBlockNumber := localblockchainState.LatestBlockNumber
	endBlockNumber := blockchainStateFromAuthority.LatestBlockNumber

	for blockNumber := startBlockNumber; blockNumber < endBlockNumber; blockNumber++ {
		blockData, err := s.getBlockDataFromCentralAuthorityByBlockNumberUseCase.Execute(ctx, blockNumber)
		if err != nil {
			s.logger.Error("Failed getting block data from the authority by block number",
				slog.Any("error", err))
			return err
		}
		s.logger.Debug("Fetched", slog.Any("block_data", blockData))

		//TODO: Save block data to local database.
	}

	//
	// STEP 6:
	// Update the local blockchain state to be equal to the blockchain state
	// of the network.
	//

	//TODO: Impl.

	s.logger.Debug("Finished sync'ing with the authority")

	return nil
}
