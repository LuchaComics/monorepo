package service

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/domain"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/usecase"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/pkg/kmutexutil"
)

type ProofOfAuthorityValidationService struct {
	config                             *config.Config
	logger                             *slog.Logger
	kmutex                             kmutexutil.KMutexProvider
	receiveProposedBlockDataDTOUseCase *usecase.ReceiveProposedBlockDataDTOUseCase
	getBlockchainLastestHashUseCase    *usecase.GetBlockchainLastestHashUseCase
	getBlockDataUseCase                *usecase.GetBlockDataUseCase
	getAccountsHashStateUseCase        *usecase.GetAccountsHashStateUseCase
	getTokensHashStateUseCase          *usecase.GetTokensHashStateUseCase
	createBlockDataUseCase             *usecase.CreateBlockDataUseCase
	setBlockchainLastestHashUseCase    *usecase.SetBlockchainLastestHashUseCase
	setBlockchainLastestTokenIDUseCase *usecase.SetBlockchainLastestTokenIDUseCase
	getAccountUseCase                  *usecase.GetAccountUseCase
	upsertAccountUseCase               *usecase.UpsertAccountUseCase
	upsertTokenUseCase                 *usecase.UpsertTokenUseCase
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
	uc8 *usecase.SetBlockchainLastestTokenIDUseCase,
	uc9 *usecase.GetAccountUseCase,
	uc10 *usecase.UpsertAccountUseCase,
	uc11 *usecase.UpsertTokenUseCase,
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

	// Load up into our datastructure.
	blockData := &domain.BlockData{
		Hash:            proposedBlockData.Hash,
		Header:          proposedBlockData.Header,
		HeaderSignature: proposedBlockData.HeaderSignature,
		Trans:           proposedBlockData.Trans,
		Validator:       proposedBlockData.Validator,
	}
	block, err := domain.ToBlock(blockData)
	if err != nil {
		s.logger.Error("validator failed converting block data into a block",
			slog.Any("error", err))
		return err
	}

	//
	// STEP 3:
	// Begin by validating the proof of authority before anything else.
	//

	poaValidator := prevBlockData.Validator
	if poaValidator.Verify(blockData.HeaderSignature, blockData.Header) == false {
		s.logger.Error("validator failed validating: authority signature is invalid")
		return fmt.Errorf("validator failed validating: %v", "authority signature is invalid")
	}

	//
	// STEP 4
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
	stateRoot, err := s.getAccountsHashStateUseCase.Execute()
	if err != nil {
		s.logger.Error("validator failed getting state root",
			slog.Any("error", err))
		return err
	}

	s.logger.Info("beginning validation...",
		slog.Any("prev_hash", previousBlock.Hash),
		slog.Uint64("prev_header_number", previousBlock.Header.Number),
		slog.Any("prev_header_prev_hash", previousBlock.Header.PrevBlockHash),
		slog.Any("prev_header_stateroot", previousBlock.Header.StateRoot),
		slog.Any("current_hash", blockData.Hash),
		slog.Uint64("current_header_number", blockData.Header.Number),
		slog.Any("current_header_prev_hash", blockData.Header.PrevBlockHash),
		slog.Any("current_header_stateroot", blockData.Header.StateRoot),
	)

	if err := block.ValidateBlock(previousBlock, stateRoot); err != nil {
		// DEVELOPERS NOTE:
		// Not an error but simply a friendly warning message.
		s.logger.Warn("validator failed validating the proposed block with the previous block",
			slog.Any("error", err))
		return err
	}

	//
	// STEP 5:
	// Save to the blockchain database.
	//

	if err := s.createBlockDataUseCase.Execute(blockData.Hash, blockData.Header, blockData.HeaderSignature, blockData.Trans, blockData.Validator); err != nil {
		s.logger.Error("validator failed saving block data",
			slog.Any("error", err))
		return err
	}

	s.logger.Info("validator add new block to blockchain",
		slog.Any("hash", proposedBlockData.Hash),
		slog.Uint64("number", proposedBlockData.Header.Number),
		slog.Any("previous_hash", proposedBlockData.Header.PrevBlockHash),
	)

	if err := s.setBlockchainLastestHashUseCase.Execute(blockData.Hash); err != nil {
		s.logger.Error("validator failed saving latest hash",
			slog.Any("error", err))
		return err
	}

	if err := s.setBlockchainLastestTokenIDUseCase.Execute(blockData.Header.LatestTokenID); err != nil {
		s.logger.Error("validator failed saving latest hash",
			slog.Any("error", err))
		return err
	}

	s.logger.Debug("validator set latest hash in blockchain",
		slog.Any("hash", proposedBlockData.Hash),
	)

	//
	// STEP 6
	// Update the account in our in-memory database.
	//

	for _, blockTx := range blockData.Trans {
		if blockTx.Type == domain.TransactionTypeCoin {
			if err := s.processAccountForBlockTransaction(blockData, &blockTx); err != nil {
				s.logger.Error("Failed processing transaction",
					slog.Any("error", err))
				return err
			}
		}
	}

	//
	// STEP 7
	// Update the tokens database.
	//

	for _, tx := range blockData.Trans {
		if tx.Type == domain.TransactionTypeToken {
			if err := s.upsertTokenUseCase.Execute(tx.TokenID, tx.TokenMetadataURI); err != nil {
				s.logger.Error("Failed upserting token transaction",
					slog.Any("error", err))
				return err
			}
		}
	}

	return nil
}

// TODO: (1) Create somesort of `processAccountForBlockTransaction` service and (2) replace it with this.
func (s *ProofOfAuthorityValidationService) processAccountForBlockTransaction(blockData *domain.BlockData, blockTx *domain.BlockTransaction) error {
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

		// DEVELOPERS NOTE:
		// Do not update this accounts `Nonce`, we need to only update the
		// `Nonce` to the receiving account, i.e. the `To` account.

		if err := s.upsertAccountUseCase.Execute(acc.Address, acc.Balance, acc.Nonce); err != nil {
			s.logger.Error("Failed upserting account.",
				slog.Any("error", err))
			return err
		}

		s.logger.Debug("New `From` account balance via validator",
			slog.Any("account_address", acc.Address),
			slog.Any("balance", acc.Balance),
			slog.Any("tx_hash", blockTx.Hash),
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

		s.logger.Debug("New `To` account balance via validator",
			slog.Any("account_address", acc.Address),
			slog.Any("balance", acc.Balance),
			slog.Any("tx_hash", blockTx.Hash),
		)
	}

	return nil
}
