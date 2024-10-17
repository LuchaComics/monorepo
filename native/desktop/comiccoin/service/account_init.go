package service

import (
	"fmt"
	"log/slog"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/common/blockchain/signature"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/domain"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/usecase"
	"github.com/ethereum/go-ethereum/common"
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

// AccountInfo structure used to keep track all the transaction
// enumerations for the particular account.
type AccountInfo struct {
	Address             *common.Address
	Nonce               uint64
	TotalAmountSent     uint64
	TotalAmountReceived uint64
	AmountsSent         []uint64
	AmountsReceived     []uint64
}

func (s *InitAccountsFromBlockchainService) Execute() error {
	//
	// STEP 1:
	// Load up our Genesis block and the coinbase account.
	//
	coinbaseAccount, err := s.processGenesisBlockData()
	if err != nil {
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

	accountInfos := make(map[common.Address]*AccountInfo, 0)

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
		// transactions within this block). Skip this step if we are the
		// genesis block as that is handled elsewhere.
		//

		if blockDataHash != signature.ZeroHash && blockData.Hash != signature.ZeroHash {
			for _, blockTx := range blockData.Trans {
				if blockTx.Type == domain.TransactionTypeCoin {
					if err := s.processCoinBlockTransaction(&blockTx, accountInfos); err != nil {
						s.logger.Error("Failed processing coin block transaction",
							slog.Any("error", err))
						return err
					}
				}
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
			s.logger.Debug("Initialized accounts successfully",
				slog.Any("account_infos", accountInfos))

			// Iterate through all the transactions.
			for accountAddress, accountSummary := range accountInfos {
				// The total amount when received is subtracted from sent for
				// regular accounts, coinbase account is an exception in which
				// the received value will always be the initial coin supply.
				var balance uint64 = 0

				//
				// CASE 1 OF 2: Coinbase account.
				//

				if accountAddress == *coinbaseAccount.Address {
					balance = coinbaseAccount.Balance - accountSummary.TotalAmountSent
				}

				//
				// CASE 2 OF 2: Regular account.
				//

				if accountAddress != *coinbaseAccount.Address {
					balance = accountSummary.TotalAmountReceived - accountSummary.TotalAmountSent
				}

				// Save our account in our in-memory database.
				if err := s.upsertAccountUseCase.Execute(&accountAddress, balance, accountSummary.Nonce); err != nil {
					s.logger.Error("Failed upserting account",
						slog.Any("error", err))
					return err
				}

				// For debugging purposes only.
				s.logger.Debug("account ready in-memory",
					slog.Any("address", accountAddress),
					slog.Any("balance", balance),
					slog.Any("total_received", accountSummary.TotalAmountReceived),
					slog.Any("total_sent", accountSummary.TotalAmountSent),
					slog.Any("nonce", accountSummary.Nonce))
			}

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

func (s *InitAccountsFromBlockchainService) processGenesisBlockData() (*domain.Account, error) {
	//
	// STEP 1:
	// Load up our genesis block from local file.
	//

	var coinbaseAccount *domain.Account

	genesisBlockData, err := s.loadGenesisBlockDataUseCase.Execute()
	if err != nil {
		s.logger.Error("Failed loading up genesis block from file",
			slog.Any("error", err))
		return nil, fmt.Errorf("Failed loading up genesis block from file: %v", err)
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
			return nil, err
		}

		coinbaseAccount = &domain.Account{
			Address: genesisTx.From,
			Balance: genesisTx.Value,
			Nonce:   0,
		}
	}

	return coinbaseAccount, nil
}

func (s *InitAccountsFromBlockchainService) processCoinBlockTransaction(blockTx *domain.BlockTransaction, accountInfos map[common.Address]*AccountInfo) error {
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
		accountInfo, ok := accountInfos[*blockTx.From]
		if !ok || accountInfo == nil {
			accountInfo = &AccountInfo{
				Address:             blockTx.From,
				Nonce:               0,
				TotalAmountSent:     0,
				TotalAmountReceived: 0,
				AmountsSent:         []uint64{},
				AmountsReceived:     []uint64{},
			}
		}

		accountInfo.Nonce += 1 // Note: We do this to prevent reply attacks. (See notes in either `domain/accounts.go` or `service/genesis_init.go`)
		accountInfo.TotalAmountSent += blockTx.Value
		accountInfo.AmountsSent = append(accountInfo.AmountsSent, blockTx.Value)

		accountInfos[*blockTx.From] = accountInfo
	}

	//
	// STEP 2
	//

	if blockTx.To != nil {
		accountInfo, ok := accountInfos[*blockTx.To]
		if !ok || accountInfo == nil {
			accountInfo = &AccountInfo{
				Address:             blockTx.To,
				Nonce:               0,
				TotalAmountSent:     0,
				TotalAmountReceived: 0,
				AmountsSent:         []uint64{},
				AmountsReceived:     []uint64{},
			}
		}

		accountInfo.TotalAmountReceived += blockTx.Value
		accountInfo.AmountsReceived = append(accountInfo.AmountsReceived, blockTx.Value)

		accountInfos[*blockTx.To] = accountInfo
	}

	return nil
}
