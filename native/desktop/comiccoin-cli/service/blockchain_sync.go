package service

import (
	"context"
	"errors"
	"log/slog"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin-cli/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin-cli/usecase"
)

type LocalBlockchainSyncService struct {
	config                                                  *config.Config
	logger                                                  *slog.Logger
	getBlockchainStateFromCentralAuthorityByChainIDUseCase  *usecase.GetBlockchainStateFromCentralAuthorityByChainIDUseCase
	getGenesisBlockDataUseCase                              *usecase.GetGenesisBlockDataUseCase
	getGenesisBlockDataFromCentralAuthorityByChainIDUseCase *usecase.GetGenesisBlockDataFromCentralAuthorityByChainIDUseCase
	upsertGenesisBlockDataUseCase                           *usecase.UpsertGenesisBlockDataUseCase
}

func NewLocalBlockchainSyncService(
	config *config.Config,
	logger *slog.Logger,
	uc1 *usecase.GetBlockchainStateFromCentralAuthorityByChainIDUseCase,
	uc2 *usecase.GetGenesisBlockDataUseCase,
	uc3 *usecase.GetGenesisBlockDataFromCentralAuthorityByChainIDUseCase,
	uc4 *usecase.UpsertGenesisBlockDataUseCase,
) *LocalBlockchainSyncService {
	return &LocalBlockchainSyncService{config, logger, uc1, uc2, uc3, uc4}
}

func (s *LocalBlockchainSyncService) Execute(ctx context.Context) error {
	//
	// STEP 1. Get blockchain state.
	//

	blockchainState, err := s.getBlockchainStateFromCentralAuthorityByChainIDUseCase.Execute(ctx, s.config.Blockchain.ChainID)
	if err != nil {
		s.logger.Error("Failed getting from central authority",
			slog.Any("error", err))
		return err
	}
	if blockchainState == nil {
		dneErr := errors.New("Failed fetching from central authority with no results")
		s.logger.Error("Failed getting from central authority",
			slog.Any("error", dneErr))
		return dneErr
	}

	s.logger.Debug("Fetched latest blockchain state",
		slog.Any("resp", blockchainState))

	//
	// STEP 2: Get genesis block data from local database.
	//

	return nil
}
