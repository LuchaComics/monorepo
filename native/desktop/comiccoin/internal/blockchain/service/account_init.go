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
	loadGenesisBlockDataUseCase     *usecase.LoadGenesisBlockDataUseCase
	getBlockchainLastestHashUseCase *usecase.GetBlockchainLastestHashUseCase
	getBlockDataUseCase             *usecase.GetBlockDataUseCase
	getAccountUseCase               *usecase.GetAccountUseCase
	getAccountsHashStateUseCase     *usecase.GetAccountsHashStateUseCase
	createAccountUseCase            *usecase.CreateAccountUseCase
	upsertAccountUseCase            *usecase.UpsertAccountUseCase
}

func NewInitAccountsFromBlockchainService(
	cfg *config.Config,
	logger *slog.Logger,
	uc1 *usecase.LoadGenesisBlockDataUseCase,
	uc2 *usecase.GetBlockchainLastestHashUseCase,
	uc3 *usecase.GetBlockDataUseCase,
	uc4 *usecase.GetAccountUseCase,
	uc5 *usecase.GetAccountsHashStateUseCase,
	uc6 *usecase.CreateAccountUseCase,
	uc7 *usecase.UpsertAccountUseCase,
) *InitAccountsFromBlockchainService {
	return &InitAccountsFromBlockchainService{cfg, logger, uc1, uc2, uc3, uc4, uc5, uc6, uc7}
}

func (s *InitAccountsFromBlockchainService) Execute() error {
	//
	// STEP 1:
	// Load up our Genesis block and the coinbase account.
	//
	if err := s.processGenesisBlockData(); err != nil {
		s.logger.Error("Failed initialize genesis block data",
			slog.Any("error", err))
		return fmt.Errorf("Failed initialize genesis block data because of error: %v", err)
	}

	//
	// STEP 2:
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

			// Print the hashstate.
			hashState, err := s.getAccountsHashStateUseCase.Execute()
			if err != nil {
				s.logger.Error("Failed getting hash state of all accounts",
					slog.Any("error", err))
				return err
			}
			s.logger.Debug(fmt.Sprintf("Current blockchain accounts stateroot: %v", hashState))

			return nil
		}
	}
}

func (s *InitAccountsFromBlockchainService) processGenesisBlockData() error {
	//
	// STEP 1:
	// Load up our genesis block from local file.
	//

	genesisBlockData, err := s.loadGenesisBlockDataUseCase.Execute()
	if err != nil {
		s.logger.Error("Failed loading up genesis block from file",
			slog.Any("error", err))
		return fmt.Errorf("Failed loading up genesis block from file: %v", err)
	}
	if genesisBlockData != nil {
		//
		// STEP 2:
		// Initialize our coinbase account
		//

		genesisTx := genesisBlockData.Trans[0]

		s.logger.Debug("processing genesis block tx.",
			slog.Any("from", genesisTx.From),
			slog.Any("to", genesisTx.To),
			slog.Uint64("value", genesisTx.Value),
			slog.Uint64("tip", genesisTx.Tip),
			slog.Any("data", genesisTx.Data),
		)

		if err := s.upsertAccountUseCase.Execute(genesisTx.From, genesisTx.Value, 0); err != nil {
			s.logger.Error("Failed upserting account.",
				slog.Any("error", err))
			return err
		}

	}
	return nil
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

		s.logger.Debug("viewed tx for `from` account",
			slog.Any("tx", blockTx.Hash),
			slog.Any("addr", acc.Address),
			slog.Any("balance", acc.Balance))
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

		s.logger.Debug("viewed tx for `to` account",
			slog.Any("tx", blockTx.Hash),
			slog.Any("addr", acc.Address),
			slog.Any("balance", acc.Balance))
	}

	return nil
}
