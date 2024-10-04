package service

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/usecase"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/pkg/blockchain/signature"
)

type BlockchainSyncClientService struct {
	config                                        *config.Config
	logger                                        *slog.Logger
	lastBlockDataHashDTOSendP2PRequestUseCase     *usecase.BlockchainLastestHashDTOSendP2PRequestUseCase
	lastBlockDataHashDTOReceiveP2PResponseUseCase *usecase.BlockchainLastestHashDTOReceiveP2PResponseUseCase
	getBlockchainLastestHashUseCase               *usecase.GetBlockchainLastestHashUseCase
	setBlockchainLastestHashUseCase               *usecase.SetBlockchainLastestHashUseCase
	blockDataDTOSendP2PRequestUseCase             *usecase.BlockDataDTOSendP2PRequestUseCase
	blockDataDTOReceiveP2PResponsetUseCase        *usecase.BlockDataDTOReceiveP2PResponsetUseCase
	createBlockDataUseCase                        *usecase.CreateBlockDataUseCase
	getBlockDataUseCase                           *usecase.GetBlockDataUseCase
}

func NewBlockchainSyncClientService(
	cfg *config.Config,
	logger *slog.Logger,
	uc1 *usecase.BlockchainLastestHashDTOSendP2PRequestUseCase,
	uc2 *usecase.BlockchainLastestHashDTOReceiveP2PResponseUseCase,
	uc3 *usecase.GetBlockchainLastestHashUseCase,
	uc4 *usecase.SetBlockchainLastestHashUseCase,
	uc5 *usecase.BlockDataDTOSendP2PRequestUseCase,
	uc6 *usecase.BlockDataDTOReceiveP2PResponsetUseCase,
	uc7 *usecase.CreateBlockDataUseCase,
	uc8 *usecase.GetBlockDataUseCase,
) *BlockchainSyncClientService {
	return &BlockchainSyncClientService{cfg, logger, uc1, uc2, uc3, uc4, uc5, uc6, uc7, uc8}
}

func (s *BlockchainSyncClientService) Execute(ctx context.Context) error {
	s.logger.Debug("blockchain sync client running...")
	defer s.logger.Debug("blockchain sync client ran")

	//
	// STEP 1:
	// Send a request over the peer-to-peer network.
	//

	err := s.lastBlockDataHashDTOSendP2PRequestUseCase.Execute(ctx)
	if err != nil {
		if strings.Contains(err.Error(), "no peers connected") {
			s.logger.Warn("blockchain sync client waiting for clients to connect...")
			return nil
		}
		s.logger.Error("blockchain sync client failed sending request",
			slog.Any("error", err))
		return err
	}

	//
	// STEP 2:
	// Wait to receive request from the peer-to-peer network.
	//

	receivedHash, err := s.lastBlockDataHashDTOReceiveP2PResponseUseCase.Execute(ctx)
	if err != nil {
		s.logger.Error("blockchain sync client failed receiving response",
			slog.Any("error", err))
		return err
	}
	if receivedHash == "" {
		s.logger.Warn("returned hash is empty")
		return nil
	}

	s.logger.Debug("blockchain sync client received from network",
		slog.Any("network_hash", receivedHash))

	//
	// STEP 3:
	// Get the latest blockchain hash we have in our local database and compare
	// with the received hash and if there are differences then we will need to
	// download from the network the latest blockdata.
	//

	// Note: Do not handle any errors, if we have any errors then continue and
	// fetch the latest hash from network anyway.
	localHash, _ := s.getBlockchainLastestHashUseCase.Execute()

	s.logger.Debug("blockchain sync client looked up local hash",
		slog.Any("local_hash", localHash))

	if localHash != string(receivedHash) {
		s.logger.Warn("local blockchain is outdated and needs updating from network",
			slog.Any("network_hash", receivedHash),
			slog.Any("local_hash", localHash))

		//TODO: FIX BUG AND CONTINUE DEVELOPMENT

		err := s.runDownloadAndSyncBlockchainFromBlockDataHash(ctx, string(receivedHash))
		if err != nil {
			s.logger.Error("blockchain failed to download and sync",
				slog.Any("error", err))
			return err
		}

		return nil
	} else {
		s.logger.Debug("local blockchain is up-to-date")
	}

	return nil
}

func (s *BlockchainSyncClientService) runDownloadAndSyncBlockchainFromBlockDataHash(ctx context.Context, blockDataHash string) error {
	// Algorithm:
	// 1. Fetch from network the blockdata for `network_hash`
	// 2. Save blockdata to local database
	// 3. Lookup `previous_hash` in local database and check if we have it.
	// 4. If record d.n.e. locally.
	// 4a. Download blockdata from network.
	// 4b. Save to local database.
	// 4c. Lookup `previous_hash` in local database and check if we have it.
	// 4d. If record d.n.e. locally then start again at step (4a)
	// 4e. If record exists then finish
	// 5. Else finish

	//
	// STEP 1:
	// Check to see if we have the data already, and if we do then proceed to the next one.
	//

	blockData, err := s.getBlockDataUseCase.Execute(blockDataHash)
	if err != nil {
		s.logger.Error("blockchain sync client failed fetching previous block in local database",
			slog.Any("error", err))
		return err
	}
	if blockData != nil {
		// CASE 1 OF 3: Genesis block reached.
		if blockData.Header.PrevBlockHash == signature.ZeroHash {
			s.logger.Debug("blockchain sync client reached genesis block data, sync completed")
			return nil
		}

		// CASE 2 OF 3: Database error
		if blockData.Header.PrevBlockHash == "" {
			err := fmt.Errorf("blockchain sync client has database error with block that has empty prev block hash: %v", blockData.Header.PrevBlockHash)
			return err
		}

		// CASE 3 OF 3: Non-genesis block reached.
		// Recursively call this function again to perform the sync.
		return s.runDownloadAndSyncBlockchainFromBlockDataHash(ctx, blockData.Header.PrevBlockHash)
	}

	//
	//
	// STEP 1:
	// Send a request over the peer-to-peer network.
	//

	if err := s.blockDataDTOSendP2PRequestUseCase.Execute(ctx, blockDataHash); err != nil {
		if strings.Contains(err.Error(), "no peers connected") {
			s.logger.Warn("blockchain sync client waiting for clients to connect...",
				slog.Any("hash", blockDataHash))
			return nil
		}
		s.logger.Error("blockchain sync client failed sending request",
			slog.Any("hash", blockDataHash),
			slog.Any("error", err))
		return err
	}

	//
	// STEP 2:
	// Wait to receive request from the peer-to-peer network.
	//

	receivedBlockData, err := s.blockDataDTOReceiveP2PResponsetUseCase.Execute(ctx)
	if err != nil {
		s.logger.Error("blockchain sync client failed receiving response",
			slog.Any("hash", blockDataHash),
			slog.Any("error", err))
		return err
	}
	if receivedBlockData == nil {
		s.logger.Warn("blockchain sync client returned empty data from network",
			slog.Any("hash", blockDataHash))
		return nil
	}

	//
	// STEP 3:
	// Save to our local database.
	//

	//TODO: FUTURE IMPROVEMENT: Security / validation / etc.

	if err := s.createBlockDataUseCase.Execute(receivedBlockData.Hash, receivedBlockData.Header, receivedBlockData.Trans); err != nil {
		s.logger.Error("blockchain sync client failed saving to local database.",
			slog.Any("error", err))
		return err
	}

	s.logger.Debug("blockchain sync client downloaded block data from network",
		slog.Any("hash", receivedBlockData.Hash))

	//
	// STEP 3:
	// Lookup the `previous_hash` in our local database and if it does not
	// exist then we repeat.
	//

	// CASE 1 OF 3: Genesis block reached.
	if blockDataHash == signature.ZeroHash {
		s.logger.Debug("blockchain sync client reached genesis block data, sync completed")
		return nil
	}

	// CASE 2 OF 3: Database error
	if receivedBlockData.Header.PrevBlockHash == "" {
		err := fmt.Errorf("blockchain sync client has database error with block that has empty prev block hash: %v", receivedBlockData)
		return err
	}

	// CASE 3 OF 3: Non-genesis block reached.
	// Recursively call this function again to perform the sync.
	return s.runDownloadAndSyncBlockchainFromBlockDataHash(ctx, receivedBlockData.Header.PrevBlockHash)
}
