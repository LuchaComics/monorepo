package service

import (
	"fmt"
	"log/slog"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/domain"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/usecase"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/pkg/blockchain/signature"
)

// Service will iterate through every single block in the blockchain and
// populate the entire in-memory database of the accounts and their balances.
type InitAccountsFromBlockchainService struct {
	config                          *config.Config
	logger                          *slog.Logger
	getBlockchainLastestHashUseCase *usecase.GetBlockchainLastestHashUseCase
	getBlockDataUseCase             *usecase.GetBlockDataUseCase
	getAccountUseCase               *usecase.GetAccountUseCase
	createAccountUseCase            *usecase.CreateAccountUseCase
	upsertAccountUseCase            *usecase.UpsertAccountUseCase
}

func NewInitAccountsFromBlockchainService(
	cfg *config.Config,
	logger *slog.Logger,
	uc1 *usecase.GetBlockchainLastestHashUseCase,
	uc2 *usecase.GetBlockDataUseCase,
	uc3 *usecase.GetAccountUseCase,
	uc4 *usecase.CreateAccountUseCase,
	uc5 *usecase.UpsertAccountUseCase,
) *InitAccountsFromBlockchainService {
	return &InitAccountsFromBlockchainService{cfg, logger, uc1, uc2, uc3, uc4, uc5}
}

func (s *InitAccountsFromBlockchainService) Execute() error {
	//
	// STEP 1:
	// Start by looking up the latest hash in the blockchain. If nothing is
	// returned then no need to worry as this means this is a new node that was
	// created and the blockchain will be downloaded from the distributed / peer
	// -to-peer network later.
	//

	blockDataHash, err := s.getBlockchainLastestHashUseCase.Execute()
	if err != nil {
		s.logger.Error("Failed initialize accounts from blockchain",
			slog.Any("error", err))
		return fmt.Errorf("Failed initialize accounts from blockchain because of error: %v", err)
	}
	if blockDataHash == "" {
		s.logger.Debug("No local blockchain exists yet, skipping initialization of accounts")
		return nil
	}

	for {
		blockData, err := s.getBlockDataUseCase.Execute(blockDataHash)
		if err != nil {
			s.logger.Error("Failed looking up block data in blockchain",
				slog.Any("error", err))
			return fmt.Errorf("Failed looking up block data in blockchain because of error: %v", err)
		}
		if blockData == nil {
			return fmt.Errorf("Block data does not exist for hash: %v", blockDataHash)
		}

		//
		// STEP 3:
		// Process the account for this particular block data (and all the
		// transactions within this block).
		//

		for _, blockTx := range blockData.Trans {
			if err := s.processBlockTransaction(blockData, &blockTx); err != nil {
				s.logger.Error("Failed processing transaction",
					slog.Any("error", err))
				return err
			}
		}

		//
		// STEP 4:
		// Get the previous block data hash and terminate if we found the
		// genesis block.
		//

		blockDataHash = blockData.Header.PrevBlockHash

		// Check if the genesis block has been reached and if so then exit.
		if blockDataHash == signature.ZeroHash && blockData.Hash == signature.ZeroHash {
			s.logger.Debug("Initialized accounts successfully")
			return nil
		}
	}
}

func (s *InitAccountsFromBlockchainService) processBlockTransaction(blockData *domain.BlockData, blockTx *domain.BlockTransaction) error {
	s.logger.Debug("processing block tx.",
		slog.Any("from", blockTx.From),
		slog.Any("to", blockTx.To),
		slog.Uint64("value", blockTx.Value),
		slog.Uint64("tip", blockTx.Tip),
		slog.Any("data", blockTx.Data),
	)

	//
	// STEP 1
	//

	if blockTx.From != nil {
		acc, _ := s.getAccountUseCase.Execute(blockTx.From)
		if acc == nil {
			if err := s.createAccountUseCase.Execute(blockTx.From); err != nil {
				s.logger.Error("Failed creating account.",
					slog.Any("error", err))
				return err
			}
			acc = &domain.Account{
				Address: blockTx.From,

				// Since we are iterating in reverse in the blockchain, we are
				// starting at the latest block data and then iterating until
				// we reach a genesis; therefore, if this account is created then
				// this is their most recent transaction so therefore we want to
				// save the nonce.
				Nonce: blockData.Header.Nonce,
			}
		}
		acc.Balance -= blockTx.Value

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
		acc, _ := s.getAccountUseCase.Execute(blockTx.To)
		if acc == nil {
			if err := s.createAccountUseCase.Execute(blockTx.To); err != nil {
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
