package service

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"time"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/common/kmutexutil"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/domain"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/usecase"
)

type ProofOfWorkValidationService struct {
	config                                       *config.Config
	logger                                       *slog.Logger
	kmutex                                       kmutexutil.KMutexProvider
	receiveProposedBlockDataDTOUseCase           *usecase.ReceiveProposedBlockDataDTOUseCase
	getBlockchainLastestHashUseCase              *usecase.GetBlockchainLastestHashUseCase
	getBlockDataUseCase                          *usecase.GetBlockDataUseCase
	getAccountsHashStateUseCase                  *usecase.GetAccountsHashStateUseCase
	createBlockDataUseCase                       *usecase.CreateBlockDataUseCase
	setBlockchainLastestHashUseCase              *usecase.SetBlockchainLastestHashUseCase
	getAccountUseCase                            *usecase.GetAccountUseCase
	upsertAccountUseCase                         *usecase.UpsertAccountUseCase
	upsertTokenIfPreviousTokenNonceGTEUseCase    *usecase.UpsertTokenIfPreviousTokenNonceGTEUseCase
	getBlockchainLastestTokenIDUseCase           *usecase.GetBlockchainLastestTokenIDUseCase
	setBlockchainLastestTokenIDIfGreatestUseCase *usecase.SetBlockchainLastestTokenIDIfGreatestUseCase
}

func NewProofOfWorkValidationService(
	cfg *config.Config,
	logger *slog.Logger,
	kmutex kmutexutil.KMutexProvider,
	uc1 *usecase.ReceiveProposedBlockDataDTOUseCase,
	uc2 *usecase.GetBlockchainLastestHashUseCase,
	uc3 *usecase.GetBlockDataUseCase,
	uc4 *usecase.GetAccountsHashStateUseCase,
	uc5 *usecase.CreateBlockDataUseCase,
	uc6 *usecase.SetBlockchainLastestHashUseCase,
	uc7 *usecase.GetAccountUseCase,
	uc8 *usecase.UpsertAccountUseCase,
	uc9 *usecase.UpsertTokenIfPreviousTokenNonceGTEUseCase,
	uc10 *usecase.GetBlockchainLastestTokenIDUseCase,
	uc11 *usecase.SetBlockchainLastestTokenIDIfGreatestUseCase,
) *ProofOfWorkValidationService {
	return &ProofOfWorkValidationService{cfg, logger, kmutex, uc1, uc2, uc3, uc4, uc5, uc6, uc7, uc8, uc9, uc10, uc11}
}

func (s *ProofOfWorkValidationService) Execute(ctx context.Context) error {
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

	s.logger.Info("received dto from network",
		slog.Any("hash", proposedBlockData.Hash),
	)

	// Lock the validator's database so we coordinate when we receive, validate
	// and/or save to the database.
	s.kmutex.Acquire("validator-service")
	defer s.kmutex.Release("validator-service")

	//
	// STEP 2:
	// Fetch the previous block we have.
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

	//
	// STEP 3
	// Validate our proposed block data to our blockchain.
	//

	newBlockData := &domain.BlockData{
		Hash:                 proposedBlockData.Hash,
		Header:               proposedBlockData.Header,
		HeaderSignatureBytes: proposedBlockData.HeaderSignatureBytes,
		Trans:                proposedBlockData.Trans,
		Validator:            proposedBlockData.Validator,
	}
	newBlock, err := domain.ToBlock(newBlockData)
	if err != nil {
		s.logger.Error("validator failed converting block data into a block",
			slog.Any("error", err))
		return err
	}

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
	stateRoot, err := s.getAccountsHashStateUseCase.Execute()
	if err != nil {
		s.logger.Error("validator failed getting state root",
			slog.Any("error", err))
		return err
	}
	if err := newBlock.ValidateBlock(previousBlock, stateRoot); err != nil {
		s.logger.Error("validator failed validating the proposed block with the previous block",
			slog.Any("error", err))
		return err
	}

	//
	// STEP 4
	// Validate each transaction in our proposed block data to our blockchain.
	//
	//TODO: IMPL: VALIDATE THE BLOCK TRANSACTIONS

	//
	// STEP 5:
	// Save to the blockchain database.
	//

	createNewBlockErr := s.createBlockDataUseCase.Execute(
		newBlockData.Hash,
		newBlockData.Header,
		newBlockData.HeaderSignatureBytes,
		newBlockData.Trans,
		newBlockData.Validator)
	if createNewBlockErr != nil {
		s.logger.Error("validator failed saving block data",
			slog.Any("error", createNewBlockErr))
		return createNewBlockErr
	}

	s.logger.Info("validator add new block to blockchain",
		slog.Any("hash", newBlockData.Hash),
		slog.Uint64("number", newBlockData.Header.Number),
		slog.Any("previous_hash", newBlockData.Header.PrevBlockHash),
	)

	if err := s.setBlockchainLastestHashUseCase.Execute(newBlockData.Hash); err != nil {
		s.logger.Error("validator failed saving latest hash",
			slog.Any("error", err))
		return err
	}

	s.logger.Debug("validator set latest hash in blockchain",
		slog.Any("hash", newBlockData.Hash),
	)

	//
	// STEP 6
	// Update the account in our in-memory database.
	//

	for _, blockTx := range newBlockData.Trans {
		if blockTx.Type == domain.TransactionTypeCoin {
			if err := s.processAccountForBlockTransaction(&blockTx); err != nil {
				s.logger.Error("Failed processing coin transaction",
					slog.Any("error", err))
				return err
			}
		}

		if blockTx.Type == domain.TransactionTypeToken {
			if err := s.processTokenForBlockTransaction(&blockTx); err != nil {
				s.logger.Error("Failed processing token transaction",
					slog.Any("error", err))
				return err
			}
		}
	}

	return nil
}

func (s *ProofOfWorkValidationService) processAccountForBlockTransaction(blockTx *domain.BlockTransaction) error {
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

func (s *ProofOfWorkValidationService) processTokenForBlockTransaction(blockTx *domain.BlockTransaction) error {
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
		log.Fatalf("processTokenForBlockTransaction: DB corruption b/c of error - you will need to re-create the db!")
	}

	// DEVELOPERS NOTE:
	// This code will execute when we mint new tokens, it will not execute if
	// we are `transfering` or `burning` a token since no new token IDs are
	// created.
	if err := s.setBlockchainLastestTokenIDIfGreatestUseCase.Execute(blockTx.TokenID); err != nil {
		s.logger.Error("validator failed saving latest hash",
			slog.Any("error", err))
		log.Fatalf("processTokenForBlockTransaction: DB corruption b/c of error - you will need to re-create the db!")
	}

	return nil
}
