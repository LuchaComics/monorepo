package service

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/common/blockchain/merkle"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/common/kmutexutil"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/domain"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/usecase"
)

type ProofOfAuthorityMiningService struct {
	config                                    *config.Config
	logger                                    *slog.Logger
	kmutex                                    kmutexutil.KMutexProvider
	getKeyService                             *GetKeyService
	storageTransactionOpenUseCase             *usecase.StorageTransactionOpenUseCase
	storageTransactionCommitUseCase           *usecase.StorageTransactionCommitUseCase
	storageTransactionDiscardUseCase          *usecase.StorageTransactionDiscardUseCase
	getAccountUseCase                         *usecase.GetAccountUseCase
	getAccountsHashStateUseCase               *usecase.GetAccountsHashStateUseCase
	getTokenUseCase                           *usecase.GetTokenUseCase
	getTokensHashStateUseCase                 *usecase.GetTokensHashStateUseCase
	listAllPendingBlockTransactionUseCase     *usecase.ListAllPendingBlockTransactionUseCase
	getBlockchainLastestHashUseCase           *usecase.GetBlockchainLastestHashUseCase
	getBlockDataUseCase                       *usecase.GetBlockDataUseCase
	proofOfWorkUseCase                        *usecase.ProofOfWorkUseCase
	createBlockDataUseCase                    *usecase.CreateBlockDataUseCase
	broadcastProposedBlockDataDTOUseCase      *usecase.BroadcastProposedBlockDataDTOUseCase
	deleteAllPendingBlockTransactionUseCase   *usecase.DeleteAllPendingBlockTransactionUseCase
	upsertTokenIfPreviousTokenNonceGTEUseCase *usecase.UpsertTokenIfPreviousTokenNonceGTEUseCase
	upsertAccountUseCase                      *usecase.UpsertAccountUseCase
	setBlockchainLastestHashUseCase           *usecase.SetBlockchainLastestHashUseCase
	getBlockchainLastestTokenIDUseCase        *usecase.GetBlockchainLastestTokenIDUseCase
	setBlockchainLastestTokenIDIfGTUseCase    *usecase.SetBlockchainLastestTokenIDIfGTUseCase
}

func NewProofOfAuthorityMiningService(
	config *config.Config,
	logger *slog.Logger,
	kmutex kmutexutil.KMutexProvider,
	getKeyService *GetKeyService,
	uc1 *usecase.StorageTransactionOpenUseCase,
	uc2 *usecase.StorageTransactionCommitUseCase,
	uc3 *usecase.StorageTransactionDiscardUseCase,
	uc4 *usecase.GetAccountUseCase,
	uc5 *usecase.GetAccountsHashStateUseCase,
	uc6 *usecase.GetTokenUseCase,
	uc7 *usecase.GetTokensHashStateUseCase,
	uc8 *usecase.ListAllPendingBlockTransactionUseCase,
	uc9 *usecase.GetBlockchainLastestHashUseCase,
	uc10 *usecase.GetBlockDataUseCase,
	uc11 *usecase.ProofOfWorkUseCase,
	uc12 *usecase.CreateBlockDataUseCase,
	uc13 *usecase.BroadcastProposedBlockDataDTOUseCase,
	uc14 *usecase.DeleteAllPendingBlockTransactionUseCase,
	uc15 *usecase.UpsertTokenIfPreviousTokenNonceGTEUseCase,
	uc16 *usecase.UpsertAccountUseCase,
	uc17 *usecase.SetBlockchainLastestHashUseCase,
	uc18 *usecase.GetBlockchainLastestTokenIDUseCase,
	uc19 *usecase.SetBlockchainLastestTokenIDIfGTUseCase,
) *ProofOfAuthorityMiningService {
	return &ProofOfAuthorityMiningService{config, logger, kmutex, getKeyService, uc1, uc2, uc3, uc4, uc5, uc6, uc7, uc8, uc9, uc10, uc11, uc12, uc13, uc14, uc15, uc16, uc17, uc18, uc19}
}

func (s *ProofOfAuthorityMiningService) Execute(ctx context.Context) error {
	// s.logger.Debug("starting mining service...")
	// defer s.logger.Debug("finished mining service")

	//
	// STEP 1:
	// Lock this function - this is important - because it will fetch all the
	// latest pending block transactions, so there needs to be a lockdown of
	// this service that when it runs it will no longer accept any more calls
	// from the system. Therefore we are using a key-based mutex to lock down
	// this service to act as a singleton runtime usage.
	//

	// Lock the mining service until it has completed executing (or errored).
	s.kmutex.Acquire("mining-service")
	defer s.kmutex.Release("mining-service")

	pendingBlockTxs, err := s.listAllPendingBlockTransactionUseCase.Execute()
	if err != nil {
		s.logger.Error("failed listing pending block transactions",
			slog.Any("error", err))
		return nil
	}
	if len(pendingBlockTxs) <= 0 {
		// s.logger.Debug("skipped mining: has no pending block transactions")
		return nil
	}

	s.logger.Info("PoA mining service executing mining for pending block transactions",
		slog.Int("count", len(pendingBlockTxs)),
	)

	//
	// STEP 2:
	// Lookup the most recent block (data) in our blockchain and setup whatever
	// variables we will need to assist in our PoA mining service.
	//

	prevBlockDataHash, err := s.getBlockchainLastestHashUseCase.Execute()
	if err != nil {
		s.logger.Error("Failed to get last hash in database",
			slog.Any("error", err))
		return fmt.Errorf("Failed to get last hash in database: %v", err)
	}
	if prevBlockDataHash == "" {
		s.logger.Error("Blockchain not initialized error")
		return fmt.Errorf("Error: %v", "blockchain not initialized")
	}
	prevBlockData, err := s.getBlockDataUseCase.Execute(string(prevBlockDataHash))
	if err != nil {
		s.logger.Error("Error getting block data frin database",
			slog.Any("error", err))
		return fmt.Errorf("Failed to get lastest block data from database: %v", err)
	}
	if prevBlockData == nil {
		s.logger.Error("Block data does not exist in database",
			slog.Any("hash", prevBlockDataHash))
		return fmt.Errorf("Block data does not exist in database for hash: %v", prevBlockDataHash)
	}
	poaValidator := prevBlockData.Validator

	// Apply whatever fees we request by the authority...
	gasPrice := uint64(s.config.Blockchain.GasPrice)
	unitsOfGas := uint64(s.config.Blockchain.UnitsOfGas)

	// Variable used to create the transactions to store on the blockchain.
	trans := make([]domain.BlockTransaction, 0)

	//
	// STEP (PRE) 3
	// Start a transaction in the database and if any errors occur then
	// we will need to discard the transaction. On success then we commit
	// the storage transaction.
	//

	s.logger.Debug("PoA mining service starting storage transaction...")
	if err := s.storageTransactionOpenUseCase.Execute(); err != nil {
		s.logger.Error("failed opening storage transaction",
			slog.Any("error", err))
		return nil
	}
	s.logger.Debug("PoA mining service started storage transaction")

	//
	// STEP 3:
	// Iterate through all the pending transactions and update our node's local
	// database. Afterwords create the block transaction which we will include
	// in our blockchain `block` and then perform our mining.
	//

	for _, pendingBlockTx := range pendingBlockTxs {
		//
		// STEP 3 (A):
		// VALIDATION
		//

		if err := s.verifyPendingBlockTransaction(pendingBlockTx); err != nil {
			s.logger.Error("Failed verifying the pending block transaction",
				slog.Any("error", err))
			s.storageTransactionDiscardUseCase.Execute()
			s.logger.Debug("PoA mining service discarded storage transaction")
			return err
		}

		//
		// STEP 3 (B):
		// Process tokens.
		//

		if pendingBlockTx.Type == domain.TransactionTypeToken {
			if err := s.processTokenForPendingBlockTransaction(pendingBlockTx); err != nil {
				s.logger.Error("Failed processing token in pending block transaction",
					slog.Any("error", err))
				s.storageTransactionDiscardUseCase.Execute()
				s.logger.Debug("PoA mining service discarded storage transaction")
				return err
			}
		}

		//
		// STEP 3 (C):
		// Process coins.
		//

		if pendingBlockTx.Type == domain.TransactionTypeCoin {
			if err := s.processAccountForPendingBlockTransaction(pendingBlockTx); err != nil {
				s.logger.Error("Failed processing account in pending block transaction",
					slog.Any("error", err))
				s.storageTransactionDiscardUseCase.Execute()
				s.logger.Debug("PoA mining service discarded storage transaction")
				return err
			}
		}

		// Create our block.
		blockTx := domain.BlockTransaction{
			SignedTransaction: *pendingBlockTx.MempoolTransaction.ToSignedTransaction(),
			TimeStamp:         uint64(time.Now().UTC().UnixMilli()),
			GasPrice:          gasPrice,
			GasUnits:          unitsOfGas,
		}
		trans = append(trans, blockTx)
	}

	// Construct a merkle tree from the transaction for this block. The root
	// of this tree will be part of the block to be mined.
	tree, err := merkle.NewTree(trans)
	if err != nil {
		s.logger.Error("Failed creating merkle tree",
			slog.Any("error", err))
		s.storageTransactionDiscardUseCase.Execute()
		s.logger.Debug("PoA mining service discarded storage transaction")
		return err
	}

	// Query the local database and get the most recent token ID.
	latestTokenID, err := s.getBlockchainLastestTokenIDUseCase.Execute()
	if err != nil {
		s.logger.Error("Failed getting blockchains latest token id",
			slog.Any("error", err))
		s.storageTransactionDiscardUseCase.Execute()
		s.logger.Debug("PoA mining service discarded storage transaction")
		return err
	}

	// Iterate through all the accounts from this local machine, sort them and
	// then hash all of them - this hash represents our `stateRoot` which is
	// in essence a snapshot of the current accounts and their balances. Why is
	// this important?
	//
	// At the start of creating a new block to be mined, a hash of this map is
	// created and stored in the block under the StateRoot field. This allows
	// each node to validate the current state of the peer’s accounts database
	// as part of block validation.
	//
	// It’s critically important that the order of the account balances are
	// exact when hashing the data. The Go spec does not define the order of
	// map iteration and leaves it up to the compiler. Since Go 1.0,
	// the compiler has chosen to have map iteration be random. This function
	//  sorts the accounts and their balances into a slice first and then
	// performs a hash of that slice.
	//
	// When a new block is received, the node can take a hash of their current
	// accounts database and match that to the StateRoot field in the
	// block header. If these hash values don’t match, then there is fraud
	//  going on by the peer and their block would be rejected.	//
	//
	// SPECIAL THANKS TO:
	// https://www.ardanlabs.com/blog/2022/05/blockchain-04-fraud-detection.html
	//
	stateRoot, err := s.getAccountsHashStateUseCase.Execute()
	if err != nil {
		s.logger.Error("Failed getting accounts hash state",
			slog.Any("error", err))
		s.storageTransactionDiscardUseCase.Execute()
		s.logger.Debug("PoA mining service discarded storage transaction")
		return err
	}

	// Ensure tokens are not tampered with.
	tokensRoot, err := s.getTokensHashStateUseCase.Execute()
	if err != nil {
		s.logger.Error("Failed getting tokens hash state",
			slog.Any("error", err))
		s.storageTransactionDiscardUseCase.Execute()
		s.logger.Debug("PoA mining service discarded storage transaction")
		return err
	}

	// Construct the genesis block.
	block := domain.Block{
		Header: &domain.BlockHeader{
			Number:        prevBlockData.Header.Number + 1,
			PrevBlockHash: string(prevBlockDataHash),
			TimeStamp:     uint64(time.Now().UTC().UnixMilli()),
			Beneficiary:   prevBlockData.Header.Beneficiary,
			Difficulty:    s.config.Blockchain.Difficulty,
			MiningReward:  s.config.Blockchain.MiningReward,
			StateRoot:     stateRoot,
			TransRoot:     tree.RootHex(), //
			Nonce:         0,              // Will be identified by the PoW algorithm.
			LatestTokenID: latestTokenID,  // Ensure our blockchain state has always the latest token ID recorded.
			TokensRoot:    tokensRoot,
		},
		HeaderSignatureBytes: []byte{}, // Will be identified by the PoA algorithm in this function!
		MerkleTree:           tree,
	}

	//
	// STEP 4:
	// Execute the proof of work to find our nonce to meet the hash difficulty.
	//

	nonce, powErr := s.proofOfWorkUseCase.Execute(ctx, &block, s.config.Blockchain.Difficulty)
	if powErr != nil {
		s.logger.Error("Failed to mine block header",
			slog.Any("error", err))
		s.storageTransactionDiscardUseCase.Execute()
		s.logger.Debug("PoA mining service discarded storage transaction")
		return err
	}

	block.Header.Nonce = nonce

	// Convert into saving for our database and transmitting over network.
	blockData := domain.NewBlockData(block)

	//
	// STEP 5
	// Our proof-of-authority signs this block data's header.
	//

	coinbaseAccountKey, err := s.getKeyService.Execute(
		s.config.Blockchain.ProofOfAuthorityAccountAddress,
		s.config.Blockchain.ProofOfAuthorityWalletPassword)
	if err != nil {
		s.logger.Error("Failed getting account wallet key",
			slog.Any("error", err))
		s.storageTransactionDiscardUseCase.Execute()
		s.logger.Debug("PoA mining service discarded storage transaction")
		return err
	}
	if coinbaseAccountKey == nil {
		err := fmt.Errorf("failed getting account wallet key: %v", "does not exist")
		s.logger.Error("Failed getting account wallet key",
			slog.Any("poa_address", s.config.Blockchain.ProofOfAuthorityAccountAddress),
			slog.Any("poa_password", s.config.Blockchain.ProofOfAuthorityWalletPassword),
			slog.Any("error", err))
		return err
	}

	coinbasePrivateKey := coinbaseAccountKey.PrivateKey
	blockDataHeaderSignatureBytes, err := poaValidator.Sign(coinbasePrivateKey, blockData.Header)
	if err != nil {
		s.logger.Error("Failed to sign block header",
			slog.Any("error", err))
		s.storageTransactionDiscardUseCase.Execute()
		return err
	}
	blockData.HeaderSignatureBytes = blockDataHeaderSignatureBytes
	blockData.Validator = poaValidator

	s.logger.Info("PoA mining completed",
		slog.String("hash", blockData.Hash),
		slog.Uint64("block_number", blockData.Header.Number),
		slog.String("prev_block_hash", blockData.Header.PrevBlockHash),
		// slog.Uint64("timestamp", blockData.Header.TimeStamp),
		// slog.String("beneficiary", blockData.Header.Beneficiary.String()),
		// slog.Uint64("difficulty", uint64(blockData.Header.Difficulty)),
		// slog.Uint64("mining_reward", blockData.Header.MiningReward),
		slog.String("state_root", blockData.Header.StateRoot),
		slog.String("trans_root", blockData.Header.TransRoot),
		// slog.Uint64("nonce", blockData.Header.Nonce),
		// slog.Uint64("latest_token_id", blockData.Header.LatestTokenID),
		// slog.Any("trans", blockData.Trans),
		// slog.Any("header_signature_bytes", blockData.HeaderSignatureBytes))
	)

	//
	// STEP 6
	// Delete purposed block data as it has already been processed.
	//

	if err := s.deleteAllPendingBlockTransactionUseCase.Execute(); err != nil {
		s.logger.Error("Failed deleting all pending block transactions",
			slog.Any("error", err))
		s.storageTransactionDiscardUseCase.Execute()
		s.logger.Debug("PoA mining service discarded storage transaction")
		return err
	}

	//
	// STEP 7:
	// Save to (local) blockchain database
	//

	if err := s.createBlockDataUseCase.Execute(blockData.Hash, blockData.Header, blockData.HeaderSignatureBytes, blockData.Trans, blockData.Validator); err != nil {
		s.logger.Error("PoA mining service failed saving block data",
			slog.Any("error", err))
		s.storageTransactionDiscardUseCase.Execute()
		s.logger.Debug("PoA mining service discarded storage transaction")
		return err
	}

	s.logger.Info("PoA mining service added new block to blockchain",
		slog.Any("hash", blockData.Hash),
		slog.Uint64("block_number", blockData.Header.Number),
		slog.String("state_root", blockData.Header.StateRoot),
		slog.Any("previous_hash", blockData.Header.PrevBlockHash),
		slog.Uint64("previous_block_number", prevBlockData.Header.Number),
		slog.String("previous_state_root", prevBlockData.Header.StateRoot),
	)

	if err := s.setBlockchainLastestHashUseCase.Execute(blockData.Hash); err != nil {
		s.logger.Error("PoA mining service failed saving latest hash",
			slog.Any("error", err))
		s.storageTransactionDiscardUseCase.Execute()
		s.logger.Debug("PoA mining service discarded storage transaction")
		return err
	}

	s.logger.Debug("PoA mining service set latest hash in blockchain",
		slog.Any("hash", blockData.Hash),
	)

	// Commit our latest changes to the database.
	s.logger.Debug("PoA mining service committing storage transaction...")
	if err := s.storageTransactionCommitUseCase.Execute(); err != nil {
		s.logger.Error("failed to commit storage transaction",
			slog.Any("error", err))
		return nil
	}
	s.logger.Debug("PoA mining service committed storage transaction")

	//
	// STEP 8:
	// Save to (distributed / peer-to-peer) blockchain database.
	//
	// Broadcast to the distributed / P2P blockchain network our new proposed
	// block data. In addition we will send this to ourselves as well.
	//
	// Developers Note:
	// When each peer (including our own) gets the broadcast message, it will
	// perform the validation on our most recent newly mined block and then
	// proceed to save to the local blockchain database if the validation
	// was successful.
	//

	// Convert to our new datastructure.
	purposeBlockData := &domain.ProposedBlockData{
		Hash:                 blockData.Hash,
		Header:               blockData.Header,
		HeaderSignatureBytes: blockData.HeaderSignatureBytes,
		Trans:                blockData.Trans,
		Validator:            blockData.Validator,
	}

	if err := s.broadcastProposedBlockDataDTOUseCase.Execute(ctx, purposeBlockData); err != nil {
		s.logger.Error("Failed to broadcast to peer-to-peer network the new block",
			slog.Any("error", err))
		return err
	}

	s.logger.Info("PoA mining service broadcasted new block data to propose to the network",
		slog.Uint64("nonce", block.Header.Nonce))

	return nil
}

func (s *ProofOfAuthorityMiningService) processAccountForPendingBlockTransaction(blockTx *domain.PendingBlockTransaction) error {
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

		if err := s.upsertAccountUseCase.Execute(acc.Address, acc.Balance, acc.Nonce); err != nil {
			s.logger.Error("Failed upserting account.",
				slog.Any("error", err))
			return err
		}

		s.logger.Debug("New `From` account balance via validator",
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

		s.logger.Debug("New `To` account balance via validator",
			slog.Any("account_address", acc.Address),
			slog.Any("balance", acc.Balance),
		)
	}

	return nil
}

func (s *ProofOfAuthorityMiningService) verifyPendingBlockTransaction(pendingBlockTx *domain.PendingBlockTransaction) error {
	fromAddr, err := pendingBlockTx.FromAddress()
	if err != nil {
		s.logger.Error("Failed getting from address",
			slog.Any("chain_id", s.config.Blockchain.ChainID),
			slog.Any("error", err))
		return err
	}

	// STEP 1: Verify that the signature is correct.
	if err := pendingBlockTx.Validate(s.config.Blockchain.ChainID, true); err != nil {
		s.logger.Error("Failed validating pending block transaction.",
			slog.Any("from_via_sig", fromAddr),
			slog.Any("chain_id", pendingBlockTx.ChainID),
			slog.Any("nonce", pendingBlockTx.Nonce),
			slog.Any("from", pendingBlockTx.From),
			slog.Any("to", pendingBlockTx.To),
			slog.Any("value", pendingBlockTx.Value),
			slog.Any("data", pendingBlockTx.Data),
			slog.Any("type", pendingBlockTx.Type),
			slog.Any("token_id", pendingBlockTx.TokenID),
			slog.Any("token_metadata_uri", pendingBlockTx.TokenMetadataURI),
			slog.Any("token_nonce", pendingBlockTx.TokenNonce),
			slog.Any("tx_sig_v", pendingBlockTx.V),
			slog.Any("tx_sig_r", pendingBlockTx.R),
			slog.Any("tx_sig_s", pendingBlockTx.S),
			slog.Any("error", err))
		return err
	}

	// STEP 2: Verify account exists in our database.
	account, err := s.getAccountUseCase.Execute(pendingBlockTx.From)
	if err != nil {
		s.logger.Error("Failed getting account.",
			slog.Any("chain_id", s.config.Blockchain.ChainID),
			slog.Any("error", err))
		return err
	}
	if account == nil {
		err := fmt.Errorf("Failed validating account: d.n.e. for: %v", pendingBlockTx.From)
		s.logger.Error("Failed validating account",
			slog.Any("chain_id", s.config.Blockchain.ChainID),
			slog.Any("error", err))
		return err
	}

	// STEP 3: Verify account has enough coins (if tx is coin-based)
	if pendingBlockTx.Type == domain.TransactionTypeCoin {
		// If the account is sending, then we need to verify the user has
		// enough coins in the balance.
		if pendingBlockTx.Value > account.Balance {
			err := fmt.Errorf("Insufficient balance in account: Have currently %v in account but transaction is requesting %v", account.Balance, pendingBlockTx.Value)
			s.logger.Error("Failed validating account",
				slog.Any("chain_id", s.config.Blockchain.ChainID),
				slog.Any("error", err))
			return err
		}
	}

	// STEP 4: Verify account belongs to token (if tx is token-based)
	if pendingBlockTx.Type == domain.TransactionTypeToken {
		// Get the token for the particular token ID.
		token, err := s.getTokenUseCase.Execute(pendingBlockTx.TokenID)
		if err != nil {
			s.logger.Error("failed getting token",
				slog.Any("error", err))
			return fmt.Errorf("failed getting token: %s", err)
		}

		// Defensive code.
		if token == nil {
			// Do nothing! This means it hasn't been created yet, meaning
			// it was just minted! So in that case all we have to do is skip
			// this function and this token will be created later on.
			return nil
			// s.logger.Warn("failed getting token",
			// 	slog.Any("token_id", pendingBlockTx.TokenID),
			// 	slog.Any("error", "token does not exist"))
			// return fmt.Errorf("failed getting token: does not exist for ID: %v", pendingBlockTx.TokenID)
		}

		// Verify the account owns the token
		if account.Address.Hex() != token.Owner.Hex() {
			s.logger.Warn("permission failed")
			return fmt.Errorf("permission denied: token address is %v but your address is %v", token.Owner.Hex(), account.Address.Hex())
		}
	}

	return nil
}

func (s *ProofOfAuthorityMiningService) processTokenForPendingBlockTransaction(pendingBlockTx *domain.PendingBlockTransaction) error {
	// Save our token to the local database ONLY if this transaction
	// is the most recent one. We track "most recent" transaction by
	// the nonce value in the token.
	err := s.upsertTokenIfPreviousTokenNonceGTEUseCase.Execute(
		pendingBlockTx.TokenID,
		pendingBlockTx.To,
		pendingBlockTx.TokenMetadataURI,
		pendingBlockTx.TokenNonce)
	if err != nil {
		s.logger.Error("Failed upserting (if previous token nonce GTE then current)",
			slog.Any("error", err))
		return err
	}

	// DEVELOPERS NOTE:
	// This code will execute when we mint new tokens, it will not execute if
	// we are `transfering` or `burning` a token since no new token IDs are
	// created.
	if err := s.setBlockchainLastestTokenIDIfGTUseCase.Execute(pendingBlockTx.TokenID); err != nil {
		s.logger.Error("validator failed saving latest hash",
			slog.Any("error", err))
		return err
	}

	return nil
}
