package service

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"strings"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/domain"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/usecase"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/pkg/blockchain/signature"
)

type MajorityVoteConsensusClientService struct {
	config                                              *config.Config
	logger                                              *slog.Logger
	consensusMechanismBroadcastRequestToNetworkUseCase  *usecase.ConsensusMechanismBroadcastRequestToNetworkUseCase
	consensusMechanismReceiveResponseFromNetworkUseCase *usecase.ConsensusMechanismReceiveResponseFromNetworkUseCase
	getBlockchainLastestHashUseCase                     *usecase.GetBlockchainLastestHashUseCase
	setBlockchainLastestHashUseCase                     *usecase.SetBlockchainLastestHashUseCase
	blockDataDTOSendP2PRequestUseCase                   *usecase.BlockDataDTOSendP2PRequestUseCase
	blockDataDTOReceiveP2PResponsetUseCase              *usecase.BlockDataDTOReceiveP2PResponsetUseCase
	createBlockDataUseCase                              *usecase.CreateBlockDataUseCase
	getBlockDataUseCase                                 *usecase.GetBlockDataUseCase
	getAccountUseCase                                   *usecase.GetAccountUseCase
	upsertAccountUseCase                                *usecase.UpsertAccountUseCase
}

func NewMajorityVoteConsensusClientService(
	cfg *config.Config,
	logger *slog.Logger,
	uc1 *usecase.ConsensusMechanismBroadcastRequestToNetworkUseCase,
	uc2 *usecase.ConsensusMechanismReceiveResponseFromNetworkUseCase,
	uc3 *usecase.GetBlockchainLastestHashUseCase,
	uc4 *usecase.SetBlockchainLastestHashUseCase,
	uc5 *usecase.BlockDataDTOSendP2PRequestUseCase,
	uc6 *usecase.BlockDataDTOReceiveP2PResponsetUseCase,
	uc7 *usecase.CreateBlockDataUseCase,
	uc8 *usecase.GetBlockDataUseCase,
	uc9 *usecase.GetAccountUseCase,
	uc10 *usecase.UpsertAccountUseCase,
) *MajorityVoteConsensusClientService {
	return &MajorityVoteConsensusClientService{cfg, logger, uc1, uc2, uc3, uc4, uc5, uc6, uc7, uc8, uc9, uc10}
}

func (s *MajorityVoteConsensusClientService) Execute(ctx context.Context) error {
	// s.logger.Debug("consensus mechanism running...")
	// defer s.logger.Debug("consensus mechanism ran")

	//
	// STEP 1:
	// Send a request over the peer-to-peer network.
	//

	err := s.consensusMechanismBroadcastRequestToNetworkUseCase.Execute(ctx)
	if err != nil {
		if strings.Contains(err.Error(), "no peers connected") {
			// s.logger.Warn("consensus mechanism waiting for clients to connect...") // For debugging purposes only.
			return nil
		}
		s.logger.Error("consensus mechanism failed sending request",
			slog.Any("error", err))
		return err
	}

	//
	// STEP 2:
	// Wait to receive request from the peer-to-peer network.
	//

	receivedHash, err := s.consensusMechanismReceiveResponseFromNetworkUseCase.Execute(ctx)
	if err != nil {
		s.logger.Error("consensus mechanism failed receiving response",
			slog.Any("error", err))
		return err
	}
	if receivedHash == "" {
		// For debugging purposes only.
		// s.logger.Warn("returned hash is empty")
		return nil
	}

	//
	// STEP 3:
	// Get the latest blockchain hash we have in our local database and compare
	// with the received hash and if there are differences then we will need to
	// download from the network the latest blockdata.
	//

	// Note: Do not handle any errors, if we have any errors then continue and
	// fetch the latest hash from network anyway.
	localHash, _ := s.getBlockchainLastestHashUseCase.Execute()

	if localHash != string(receivedHash) {
		s.logger.Warn("local blockchain is outdated and needs updating from network",
			slog.Any("network_hash", receivedHash),
			slog.Any("local_hash", localHash))

		if err := s.runDownloadAndSyncBlockchainFromBlockDataHash(ctx, string(receivedHash)); err != nil {
			s.logger.Error("blockchain failed to download and sync",
				slog.Any("error", err))
			return err
		}

		// Once our sync has been completed, we can save our latest hash so
		// we won't have to sync again.
		if err := s.setBlockchainLastestHashUseCase.Execute(string(receivedHash)); err != nil {
			s.logger.Error("blockchain failed to save latest hash to database",
				slog.Any("error", err))
			return err
		}

		// Reaching here is success!
		return nil
	} else {
		s.logger.Debug("local blockchain is up-to-date with peer-to-peer network")
	}

	return nil
}

func (s *MajorityVoteConsensusClientService) runDownloadAndSyncBlockchainFromBlockDataHash(ctx context.Context, blockDataHash string) error {
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
		s.logger.Error("consensus mechanism failed fetching previous block in local database",
			slog.Any("error", err))
		return err
	}
	if blockData != nil {
		// CASE 1 OF 3: Genesis block reached.
		if blockData.Header.PrevBlockHash == signature.ZeroHash {
			s.logger.Debug("consensus mechanism reached genesis block data, sync completed")
			return nil
		}

		// CASE 2 OF 3: Database error
		if blockData.Header.PrevBlockHash == "" {
			err := fmt.Errorf("consensus mechanism has database error with block that has empty prev block hash: %v", blockData.Header.PrevBlockHash)
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
			s.logger.Warn("consensus mechanism aborted sending request because there are no peers connected yet",
				slog.Any("hash", blockDataHash))
			return nil
		}
		s.logger.Error("consensus mechanism failed sending request",
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
		s.logger.Error("consensus mechanism failed receiving response",
			slog.Any("hash", blockDataHash),
			slog.Any("error", err))
		return err
	}
	if receivedBlockData == nil {
		s.logger.Warn("consensus mechanism returned empty data from network",
			slog.Any("hash", blockDataHash))
		return nil
	}

	//
	// STEP 3:
	// Save to our local database.
	//

	if err := s.createBlockDataUseCase.Execute(receivedBlockData.Hash, receivedBlockData.Header, receivedBlockData.Trans); err != nil {
		s.logger.Error("consensus mechanism failed saving to local database.",
			slog.Any("error", err))
		return err
	}

	s.logger.Debug("consensus mechanism downloaded block data from network",
		slog.Any("hash", receivedBlockData.Hash))

	//
	// STEP 4
	// Update the account in our in-memory database.
	//

	for _, blockTx := range receivedBlockData.Trans {
		if err := s.processAccountForBlockTransaction(blockData, &blockTx); err != nil {
			s.logger.Error("Failed processing transaction",
				slog.Any("error", err))
			return err
		}
	}

	//
	// STEP 5:
	// Lookup the `previous_hash` in our local database and if it does not
	// exist then we repeat.
	//

	// CASE 1 OF 3: Genesis block reached.
	if blockDataHash == signature.ZeroHash {
		s.logger.Debug("consensus mechanism reached genesis block data, sync completed")
		return nil
	}

	// CASE 2 OF 3: Database error
	if receivedBlockData.Header.PrevBlockHash == "" {
		err := fmt.Errorf("consensus mechanism has database error with block that has empty prev block hash: %v", receivedBlockData)
		return err
	}

	// CASE 3 OF 3: Non-genesis block reached.
	// Recursively call this function again to perform the sync.
	return s.runDownloadAndSyncBlockchainFromBlockDataHash(ctx, receivedBlockData.Header.PrevBlockHash)
}

// TODO: (1) Create somesort of `processAccountForBlockTransaction` service and (2) replace it with this.
func (s *MajorityVoteConsensusClientService) processAccountForBlockTransaction(blockData *domain.BlockData, blockTx *domain.BlockTransaction) error {
	// DEVELOPERS NOTE:
	// Please remember that when this function executes, there already is an
	// in-memory database of accounts populated and maintained by this node.
	// Therefore the code in this function is executed on a ready database.

	//
	// STEP 1
	//

	if blockTx.From != nil {
		// DEVELOPERS NOTE:
		// We already *should* have a `From` account in our database, so we can
		acc, _ := s.getAccountUseCase.Execute(blockTx.From)
		if acc == nil {
			s.logger.Error("The `From` account does not exist in our database.",
				slog.Any("hash", blockTx.From))
			log.Fatalf("The `From` account does not exist in our database for hash: %v", blockTx.From.String())
			return fmt.Errorf("The `From` account does not exist in our database for hash: %v", blockTx.From.String())
		}
		acc.Balance -= blockTx.Value

		// DEVELOPERS NOTE:
		// Do not update this accounts `Nonce`, we need to only update the
		// `Nonce` to the receiving account, i.e. the `To` account.

		if err := s.upsertAccountUseCase.Execute(acc.Address, acc.Balance, acc.Nonce); err != nil {
			s.logger.Error("Failed upserting account.",
				slog.Any("error", err))
			return err
		}
	}

	//
	// STEP 2
	//

	if blockTx.To != nil {
		// DEVELOPERS NOTE:
		// It is perfectly normal that our account would possibly not exist
		// so we would need to create a new Account record in our local
		// in-memory database.
		acc, _ := s.getAccountUseCase.Execute(blockTx.To)
		if acc == nil {
			if err := s.upsertAccountUseCase.Execute(blockTx.To, 0, 0); err != nil {
				s.logger.Error("Failed creating account.",
					slog.Any("error", err))
				return err
			}
			acc = &domain.Account{
				Address: blockTx.To,

				// Since we are iterating in reverse in the blockchain, we are
				// starting at the latest block data and then iterating until
				// we reach a genesis; therefore, if this account is created then
				// this is their most recent transaction so therefore we want to
				// save the nonce.
				Nonce: blockData.Header.Nonce,
			}
		}
		acc.Balance += blockTx.Value

		if err := s.upsertAccountUseCase.Execute(acc.Address, acc.Balance, acc.Nonce); err != nil {
			s.logger.Error("Failed upserting account.",
				slog.Any("error", err))
			return err
		}
	}

	return nil
}
