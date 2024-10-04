package service

import (
	"fmt"
	"log/slog"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/domain"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/usecase"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/pkg/blockchain/signature"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/pkg/httperror"
)

type GetAccountBalanceService struct {
	config                          *config.Config
	logger                          *slog.Logger
	getBlockchainLastestHashUseCase *usecase.GetBlockchainLastestHashUseCase
	getBlockDataUseCase             *usecase.GetBlockDataUseCase
}

func NewGetAccountBalanceService(
	cfg *config.Config,
	logger *slog.Logger,
	uc1 *usecase.GetBlockchainLastestHashUseCase,
	uc2 *usecase.GetBlockDataUseCase,
) *GetAccountBalanceService {
	return &GetAccountBalanceService{cfg, logger, uc1, uc2}
}

func (s *GetAccountBalanceService) Execute(account *domain.Account) (uint64, error) {
	//
	// STEP 1: Validation.
	//

	e := make(map[string]string)
	if account == nil {
		e["account"] = "missing value"
	}
	if len(e) != 0 {
		s.logger.Warn("Failed getting account balance",
			slog.Any("error", e))
		return 0, httperror.NewForBadRequest(&e)
	}

	//
	// STEP 2: Get the latest block in the blockchain.
	//

	currentHash, err := s.getBlockchainLastestHashUseCase.Execute()
	if err != nil {
		s.logger.Error("failed to get last hash",
			slog.Any("error", err))
		return 0, fmt.Errorf("failed to get last hash: %v", err)
	}

	//
	// STEP 3:
	// Iterate through the blockchain and compute the balance.
	//

	var balanceTotal uint64 = 0

	for {
		blockData, err := s.getBlockDataUseCase.Execute(currentHash)
		if err != nil {
			s.logger.Error("failed to get block datah",
				slog.String("hash", currentHash))
			return 0, fmt.Errorf("failed to get block data: %v", err)
		}

		// DEVELOPERS NOTE:
		// If we get a nil block then that means we have reached the genesis
		// block so we can abort.
		if blockData == nil {
			break // Genesis block reached
		}

		// DEVELOPERS NOTE:
		// Every block can have one or many transactions, therefore we will
		// need to iterate through all of them for our computation.
		for _, tx := range blockData.Trans {
			if tx.From == account.WalletAddress {
				balanceTotal -= tx.Value
			}
			if tx.To == account.WalletAddress {
				balanceTotal += tx.Value
			}
		}

		// DEVELOPERS NOTE:
		// To traverse the blockchain, we want to go always iterate through the
		// previous block, unless we reached the first block called the genesis
		// block; therefore, keep looking at the previous blocks hash and set
		// it as the current hash so when we re-run this loop, we are processing
		// for a new block.
		if blockData.Header.PrevBlockHash == signature.ZeroHash {
			break // Genesis block reached
		}
		currentHash = blockData.Header.PrevBlockHash
	}

	return balanceTotal, nil
}
