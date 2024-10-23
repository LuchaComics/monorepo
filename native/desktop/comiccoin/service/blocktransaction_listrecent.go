package service

import (
	"fmt"
	"log/slog"

	"github.com/ethereum/go-ethereum/common"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/common/blockchain/signature"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/common/httperror"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/domain"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/usecase"
)

type ListRecentBlockTransactionService struct {
	config                          *config.Config
	logger                          *slog.Logger
	getBlockchainLastestHashUseCase *usecase.GetBlockchainLastestHashUseCase
	getBlockDataUseCase             *usecase.GetBlockDataUseCase
}

func NewListRecentBlockTransactionService(
	cfg *config.Config,
	logger *slog.Logger,
	uc1 *usecase.GetBlockchainLastestHashUseCase,
	uc2 *usecase.GetBlockDataUseCase,
) *ListRecentBlockTransactionService {
	return &ListRecentBlockTransactionService{cfg, logger, uc1, uc2}
}

func (s *ListRecentBlockTransactionService) Execute(address *common.Address, limit int) ([]*domain.BlockTransaction, error) {
	//
	// STEP 1: Validation.
	//

	e := make(map[string]string)
	if address == nil {
		e["address"] = "missing value"
	}
	if len(e) != 0 {
		s.logger.Warn("Failed validating list recent block transaction",
			slog.Any("error", e))
		return nil, httperror.NewForBadRequest(&e)
	}

	//
	// STEP 2: Get the latest block in the blockchain.
	//

	currentHash, err := s.getBlockchainLastestHashUseCase.Execute()
	if err != nil {
		s.logger.Error("failed to get last hash",
			slog.Any("error", err))
		return nil, fmt.Errorf("failed to get last hash: %v", err)
	}

	//
	// STEP 3:
	// Iterate through the blockchain and compute the balance.
	//

	txs := make([]*domain.BlockTransaction, 0)

	for {
		blockData, err := s.getBlockDataUseCase.Execute(currentHash)
		if err != nil {
			s.logger.Error("failed to get block datah",
				slog.String("hash", currentHash))
			return nil, fmt.Errorf("failed to get block data: %v", err)
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
			// // For debugging purposes only.
			// s.logger.Debug("Transaction Comporator",
			// 	slog.Any("tx.From", *tx.From),
			// 	slog.Any("address", *address),
			// 	slog.Any("tx.From == address", *tx.From == *address))

			if *tx.From == *address {
				txs = append(txs, &tx)

				// Reached limit.
				if len(txs) > limit {
					break
				}
			}

			// // For debugging purposes only.
			// s.logger.Debug("Transaction Comporator",
			// 	slog.Any("tx.To", *tx.To),
			// 	slog.Any("address", *address),
			// 	slog.Any("tx.To == address", *tx.To == *address))

			if *tx.To == *address {
				txs = append(txs, &tx)

				// Reached limit.
				if len(txs) > limit {
					break
				}
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

	s.logger.Debug("Fetched",
		slog.Any("txs", txs))

	return txs, nil
}
