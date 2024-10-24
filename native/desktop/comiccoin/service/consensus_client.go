package service

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/common/blockchain/signature"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/domain"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/usecase"
)

type MajorityVoteConsensusClientService struct {
	config                                              *config.Config
	logger                                              *slog.Logger
	storageTransactionOpenUseCase                       *usecase.StorageTransactionOpenUseCase
	storageTransactionCommitUseCase                     *usecase.StorageTransactionCommitUseCase
	storageTransactionDiscardUseCase                    *usecase.StorageTransactionDiscardUseCase
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
	getAccountsHashStateUseCase                         *usecase.GetAccountsHashStateUseCase
	getTokensHashStateUseCase                           *usecase.GetTokensHashStateUseCase
}

func NewMajorityVoteConsensusClientService(
	cfg *config.Config,
	logger *slog.Logger,
	uc1 *usecase.StorageTransactionOpenUseCase,
	uc2 *usecase.StorageTransactionCommitUseCase,
	uc3 *usecase.StorageTransactionDiscardUseCase,
	uc4 *usecase.ConsensusMechanismBroadcastRequestToNetworkUseCase,
	uc5 *usecase.ConsensusMechanismReceiveResponseFromNetworkUseCase,
	uc6 *usecase.GetBlockchainLastestHashUseCase,
	uc7 *usecase.SetBlockchainLastestHashUseCase,
	uc8 *usecase.BlockDataDTOSendP2PRequestUseCase,
	uc9 *usecase.BlockDataDTOReceiveP2PResponsetUseCase,
	uc10 *usecase.CreateBlockDataUseCase,
	uc11 *usecase.GetBlockDataUseCase,
	uc12 *usecase.GetAccountUseCase,
	uc13 *usecase.UpsertAccountUseCase,
	uc14 *usecase.GetAccountsHashStateUseCase,
	uc15 *usecase.GetTokensHashStateUseCase,
) *MajorityVoteConsensusClientService {
	return &MajorityVoteConsensusClientService{cfg, logger, uc1, uc2, uc3, uc4, uc5, uc6, uc7, uc8, uc9, uc10, uc11, uc12, uc13, uc14, uc15}
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

	// Make sure that if we receive signature for the genesis block that
	// this function aborts because there's no need to sync genesis block
	// because all the peers in the network already have it. This code gets
	// called when the other consensus server is joining the network and
	// they don't have anything.
	if receivedHash == signature.ZeroHash {
		// For debugging purposes only.
		s.logger.Debug("consensus mechanism detected the other consensus server joined peer-to-peer network with only the genesis block so our consensus client will abort syncing")
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

		//
		// STEP 4:
		// Lookup both the received and local hashes, and perform the following:
		// 1. If no received hash d.n.e. then continue
		// 2. Else if received hash exists locally then compare the
		// `block_number` and if local's is greater then abort as it means
		// the other consensus server is out-of-date.
		//

		localBlockData, err := s.getBlockDataUseCase.Execute(localHash)
		if err != nil {
			s.logger.Error("consensus mechanism failed getting from local database",
				slog.Any("error", err))
			return err
		}
		receivedBlockData, err := s.getBlockDataUseCase.Execute(receivedHash)
		if err != nil {
			s.logger.Error("consensus mechanism failed getting from local database",
				slog.Any("error", err))
			return err
		}

		// Defensive code: Check if consensus server has outdated blockchain.
		if receivedBlockData != nil {
			if localBlockData.Header.Number >= receivedBlockData.Header.Number {
				s.logger.Debug("local blockchain is up-to-date with peer-to-peer network, however the other sending server has outdated blockchain...")
			}
			s.logger.Debug("local blockchain is out-of-date with peer-to-peer network, proceeding to update...")
		}

		//
		// STEP 5
		// Start a transaction in the database and if any errors occur then
		// we will need to discard the transaction. On success then we commit
		// the storage transaction.
		//

		s.logger.Debug("consensus mechanism starting storage transaction...")
		if err := s.storageTransactionOpenUseCase.Execute(); err != nil {
			s.logger.Error("failed opening storage transaction",
				slog.Any("error", err))
			return nil
		}
		s.logger.Debug("Consensus mechanism started storage transaction")

		s.logger.Warn("local blockchain is outdated and needs updating from network",
			slog.Any("network_hash", receivedHash),
			slog.Any("local_hash", localHash))

		if err := s.runDownloadAndSyncBlockchainFromBlockDataHash(ctx, string(receivedHash)); err != nil {
			s.logger.Error("blockchain failed to download and sync",
				slog.Any("error", err))
			s.storageTransactionDiscardUseCase.Execute()
			return err
		}

		// Once our sync has been completed, we can save our latest hash so
		// we won't have to sync again.
		if err := s.setBlockchainLastestHashUseCase.Execute(string(receivedHash)); err != nil {
			s.logger.Error("blockchain failed to save latest hash to database",
				slog.Any("error", err))
			s.storageTransactionDiscardUseCase.Execute()
			s.logger.Debug("Consensus mechanism discarded storage transaction")
			return err
		}
		s.logger.Debug("local blockchain was updated successfully from the peer-to-peer network")

		// Reaching here is success!
		s.logger.Debug("Consensus mechanism committing storage transaction...")
		if err := s.storageTransactionCommitUseCase.Execute(); err != nil {
			s.logger.Error("failed to commit storage transaction",
				slog.Any("error", err))
			return nil
		}
		s.logger.Debug("Consensus mechanism committed storage transaction")
	} else {
		s.logger.Debug("local blockchain is up-to-date with peer-to-peer network")
	}

	localHash, _ = s.getBlockchainLastestHashUseCase.Execute()
	accountStateRoot, err := s.getAccountsHashStateUseCase.Execute()
	if err != nil {
		s.logger.Error("failed to get account state root",
			slog.Any("error", err))
		return nil
	}
	tokenStateRoot, err := s.getTokensHashStateUseCase.Execute()
	if err != nil {
		s.logger.Error("failed to get token state root",
			slog.Any("error", err))
		return nil
	}

	s.logger.Debug("consensus reached",
		slog.String("hash", localHash),
		slog.String("account_state_root", accountStateRoot),
		slog.String("token_state_root", tokenStateRoot),
	)

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

	if err := s.createBlockDataUseCase.Execute(receivedBlockData.Hash, receivedBlockData.Header, receivedBlockData.HeaderSignatureBytes, receivedBlockData.Trans, receivedBlockData.Validator); err != nil {
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
		if err := s.processAccountForBlockTransaction(&blockTx); err != nil {
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
func (s *MajorityVoteConsensusClientService) processAccountForBlockTransaction(blockTx *domain.BlockTransaction) error {
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
			return fmt.Errorf("The `From` account does not exist in our database for hash: %v", blockTx.From.String())
		}
		acc.Balance -= blockTx.Value
		acc.Nonce += 1 // Note: We do this to prevent reply attacks. (See notes in either `domain/accounts.go` or `service/genesis_init.go`)

		// DEVELOPERS NOTE:
		// Do not update this accounts `Nonce`, we need to only update the
		// `Nonce` to the receiving account, i.e. the `To` account.

		if err := s.upsertAccountUseCase.Execute(acc.Address, acc.Balance, acc.Nonce); err != nil {
			s.logger.Error("Failed upserting account.",
				slog.Any("error", err))
			return err
		}

		s.logger.Debug("New `From` account balance via censensus",
			slog.Any("account_address", acc.Address),
			slog.Any("balance", acc.Balance),
		)
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

				// Always start by zero, increment by 1 after mining successful.
				Nonce: 0,

				Balance: blockTx.Value,
			}
		} else {
			acc.Balance += blockTx.Value
			acc.Nonce += 1 // Note: We do this to prevent reply attacks. (See notes in either `domain/accounts.go` or `service/genesis_init.go`)
		}

		if err := s.upsertAccountUseCase.Execute(acc.Address, acc.Balance, acc.Nonce); err != nil {
			s.logger.Error("Failed upserting account.",
				slog.Any("error", err))
			return err
		}

		s.logger.Debug("New `To` account balance via censensus",
			slog.Any("account_address", acc.Address),
			slog.Any("balance", acc.Balance),
		)
	}

	//
	// STEP 3
	//

	if blockTx.Type == domain.TransactionTypeToken {
		//TODO: Impl.
	}

	return nil
}
