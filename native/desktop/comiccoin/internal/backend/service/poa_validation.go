package service

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"time"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/backend/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/backend/domain"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/backend/usecase"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/pkg/kmutexutil"
)

type ProofOfAuthorityValidationService struct {
	config                                       *config.Config
	logger                                       *slog.Logger
	kmutex                                       kmutexutil.KMutexProvider
	receiveProposedBlockDataDTOUseCase           *usecase.ReceiveProposedBlockDataDTOUseCase
	getBlockchainLastestHashUseCase              *usecase.GetBlockchainLastestHashUseCase
	getBlockDataUseCase                          *usecase.GetBlockDataUseCase
	getAccountsHashStateUseCase                  *usecase.GetAccountsHashStateUseCase
	getTokensHashStateUseCase                    *usecase.GetTokensHashStateUseCase
	createBlockDataUseCase                       *usecase.CreateBlockDataUseCase
	setBlockchainLastestHashUseCase              *usecase.SetBlockchainLastestHashUseCase
	setBlockchainLastestTokenIDIfGreatestUseCase *usecase.SetBlockchainLastestTokenIDIfGreatestUseCase
	getAccountUseCase                            *usecase.GetAccountUseCase
	upsertAccountUseCase                         *usecase.UpsertAccountUseCase
	upsertTokenIfPreviousTokenNonceGTEUseCase    *usecase.UpsertTokenIfPreviousTokenNonceGTEUseCase
}

func NewProofOfAuthorityValidationService(
	cfg *config.Config,
	logger *slog.Logger,
	kmutex kmutexutil.KMutexProvider,
	uc1 *usecase.ReceiveProposedBlockDataDTOUseCase,
	uc2 *usecase.GetBlockchainLastestHashUseCase,
	uc3 *usecase.GetBlockDataUseCase,
	uc4 *usecase.GetAccountsHashStateUseCase,
	uc5 *usecase.GetTokensHashStateUseCase,
	uc6 *usecase.CreateBlockDataUseCase,
	uc7 *usecase.SetBlockchainLastestHashUseCase,
	uc8 *usecase.SetBlockchainLastestTokenIDIfGreatestUseCase,
	uc9 *usecase.GetAccountUseCase,
	uc10 *usecase.UpsertAccountUseCase,
	uc11 *usecase.UpsertTokenIfPreviousTokenNonceGTEUseCase,
) *ProofOfAuthorityValidationService {
	return &ProofOfAuthorityValidationService{cfg, logger, kmutex, uc1, uc2, uc3, uc4, uc5, uc6, uc7, uc8, uc9, uc10, uc11}
}

func (s *ProofOfAuthorityValidationService) Execute(ctx context.Context) error {
	// s.logger.Debug("starting validation service...")
	// defer s.logger.Debug("finished validation service")

	//
	// STEP 1
	// Wait to receive data (which also was validated) from the P2P network.
	//

	proposedBlockData, err := s.receiveProposedBlockDataDTOUseCase.Execute(ctx)
	if err != nil {
		s.logger.Error("validator failed receiving dto",
			slog.Any("error", err))
		return err
	}
	if proposedBlockData == nil {
		// Developer Note:
		// If we haven't received anything, that means we haven't connected to
		// the distributed / P2P network, so all we can do at the moment is to
		// pause the execution for 1 second and then retry again.
		time.Sleep(1 * time.Second)
		return nil
	}
	// Lock the validator's database so we coordinate when we receive, validate
	// and/or save to the database.
	s.kmutex.Acquire("validator-service")
	defer s.kmutex.Release("validator-service")

	s.logger.Info("received dto from network",
		slog.Any("hash", proposedBlockData.Hash),
		slog.Any("header_signature", proposedBlockData.HeaderSignature),
	)

	//
	// STEP 2:
	// Defensive code: Check to see if we already have this block in our
	// blockchain and if we do then skip this validation.
	//

	existingPropposedBlockData, err := s.getBlockDataUseCase.Execute(proposedBlockData.Hash)
	if err != nil {
		s.logger.Error("Failed to lookup existing block data",
			slog.Any("error", err))
		return err
	}
	if existingPropposedBlockData != nil {
		// Data already exists! Therefore we can abandon this request from
		// the peer-to-peer network to validate this block.
		s.logger.Warn("purposed block already exists locally, skipping validation...",
			slog.Any("hash", proposedBlockData.Hash),
			slog.Any("header_signature", proposedBlockData.HeaderSignature),
		)
		return nil

	}

	//
	// STEP 3:
	// Fetch the previous block we have and setup whatever
	// variables we will need to assist in our PoA valdiation service.
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
	previousBlock, err := domain.ToBlock(prevBlockData)
	if err != nil {
		s.logger.Error("Error converting block data to block",
			slog.Any("error", err))
		return err
	}

	// Load up into our datastructure.
	newBlockData := &domain.BlockData{
		Hash:            proposedBlockData.Hash,
		Header:          proposedBlockData.Header,
		HeaderSignature: proposedBlockData.HeaderSignature,
		Trans:           proposedBlockData.Trans,
		Validator:       proposedBlockData.Validator,
	}
	newBlock, err := domain.ToBlock(newBlockData)
	if err != nil {
		s.logger.Error("validator failed converting block data into a block",
			slog.Any("error", err))
		return err
	}

	// // Variable used to keep track the most recent `token_id` value.
	// latestTokenID := prevBlockData.Header.LatestTokenID

	//
	// STEP 4:
	// Begin by validating the proof of authority before doing anything else.
	//

	poaValidator := prevBlockData.Validator
	if poaValidator.Verify(newBlockData.HeaderSignature, newBlockData.Header) == false {
		s.logger.Error("validator failed validating: authority signature is invalid")
		return fmt.Errorf("validator failed validating: %v", "authority signature is invalid")
	}

	//
	// STEP 5:
	// Iterate through all the pending transactions and perform various
	// computations...
	//

	for _, blockTx := range newBlockData.Trans {

		//
		// STEP 5 (A):
		// Process coins.
		//

		if blockTx.Type == domain.TransactionTypeCoin {
			if err := s.processAccountForBlockTransaction(&blockTx); err != nil {
				s.logger.Error("Failed processing transaction",
					slog.Any("error", err))
				log.Fatalf("DB corruption b/c of error - you will need to re-create the db!")
			}
		}

		//
		// STEP 5 (B):
		// Process tokens.
		//

		if blockTx.Type == domain.TransactionTypeToken {
			if err := s.processTokenForBlockTransaction(&blockTx); err != nil {
				s.logger.Error("Failed processing token transaction",
					slog.Any("error", err))
				return err
			}
		}
	}

	//
	// STEP 6
	// Afterwards validate our proposed block data to our blockchain.
	//

	// DEVELOPERS NOTE:
	// In essence every node on the network hash an in-memory database of all
	// the accounts (including the account balances) before we add this purposed
	// block to the blockchain; therefore, we can confirm the `StateRoot` is
	// the same on the miners side to confirm that no modifications were made
	// with any of the account balances.
	//
	// To learn more about the state root, read this in-depth articl:
	// https://www.ardanlabs.com/blog/2022/05/blockchain-04-fraud-detection.html
	//
	currentStateRootInThisNode, err := s.getAccountsHashStateUseCase.Execute()
	if err != nil {
		s.logger.Error("validator failed getting state root",
			slog.Any("error", err))
		log.Fatalf("DB corruption b/c of error - you will need to re-create the db!")
	}

	// Ensure tokens are not tampered with.
	currentTokensRootInThisNode, err := s.getTokensHashStateUseCase.Execute()
	if err != nil {
		s.logger.Error("Failed getting tokens hash state",
			slog.Any("error", err))
		log.Fatalf("DB corruption b/c of error - you will need to re-create the db!")
	}
	_ = currentTokensRootInThisNode //TODO: Add this feature when we are ready.

	s.logger.Info("beginning validation...",
		slog.Any("prev_hash", previousBlock.Hash),
		slog.Uint64("prev_header_number", previousBlock.Header.Number),
		slog.Any("prev_header_prev_hash", previousBlock.Header.PrevBlockHash),
		slog.Any("prev_header_stateroot", previousBlock.Header.StateRoot),
		slog.Any("current_hash", newBlockData.Hash),
		slog.Uint64("current_header_number", newBlockData.Header.Number),
		slog.Any("current_header_prev_hash", newBlockData.Header.PrevBlockHash),
		slog.Any("current_header_stateroot", newBlockData.Header.StateRoot),
	)

	if err := newBlock.ValidateBlock(previousBlock, currentStateRootInThisNode); err != nil {
		// DEVELOPERS NOTE:
		// Not an error but simply a friendly warning message.
		s.logger.Warn("validator failed validating the proposed block with the previous block",
			slog.Any("error", err))
		log.Fatalf("DB corruption b/c of error - you will need to re-create the db!")
	}

	//
	// STEP 7:
	// Save to the (local) blockchain database.
	//

	if err := s.createBlockDataUseCase.Execute(newBlockData.Hash, newBlockData.Header, newBlockData.HeaderSignature, newBlockData.Trans, newBlockData.Validator); err != nil {
		s.logger.Error("validator failed saving block data",
			slog.Any("error", err))
		log.Fatalf("DB corruption b/c of error - you will need to re-create the db!")
	}

	s.logger.Info("validator add new block to blockchain",
		slog.Any("hash", proposedBlockData.Hash),
		slog.Uint64("number", proposedBlockData.Header.Number),
		slog.Any("previous_hash", proposedBlockData.Header.PrevBlockHash),
	)

	if err := s.setBlockchainLastestHashUseCase.Execute(newBlockData.Hash); err != nil {
		s.logger.Error("validator failed saving latest hash",
			slog.Any("error", err))
		log.Fatalf("DB corruption b/c of error - you will need to re-create the db!")
	}

	s.logger.Debug("validator set latest hash in blockchain",
		slog.Any("hash", proposedBlockData.Hash),
	)

	return nil
}

func (s *ProofOfAuthorityValidationService) processAccountForBlockTransaction(blockTx *domain.BlockTransaction) error {
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

func (s *ProofOfAuthorityValidationService) processTokenForBlockTransaction(blockTx *domain.BlockTransaction) error {
	// Save our token to the local database ONLY if this transaction
	// is the most recent one. We track "most recent" transaction by
	// the nonce value in the token.
	err := s.upsertTokenIfPreviousTokenNonceGTEUseCase.Execute(
		blockTx.TokenID,
		blockTx.To,
		blockTx.TokenMetadataURI,
		blockTx.TokenNonce)
	if err != nil {
		s.logger.Error("Failed upserting (if previous token nonce GTE then current)",
			slog.Any("error", err))
		log.Fatalf("DB corruption b/c of error - you will need to re-create the db!")
	}

	// DEVELOPERS NOTE:
	// This code will execute when we mint new tokens, it will not execute if
	// we are `transfering` or `burning` a token since no new token IDs are
	// created.
	if err := s.setBlockchainLastestTokenIDIfGreatestUseCase.Execute(blockTx.TokenID); err != nil {
		s.logger.Error("validator failed saving latest hash",
			slog.Any("error", err))
		log.Fatalf("DB corruption b/c of error - you will need to re-create the db!")
	}

	return nil
}
