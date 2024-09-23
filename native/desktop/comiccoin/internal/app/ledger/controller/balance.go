package controller

import (
	"context"
	"fmt"
	"log/slog"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

func (impl *ledgerControllerImpl) GetBalanceByAddress(ctx context.Context, address common.Address) (*big.Int, error) {
	balance := new(big.Int)
	currentHash, err := impl.lastHashStorer.Get(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get last hash: %v", err)
	}
	impl.logger.Debug("lookup last hash",
		slog.String("last_hash", currentHash))

	// Iterate through all the blocks.
	for {
		block, err := impl.blockStorer.GetByHash(ctx, currentHash)
		if err != nil {
			impl.logger.Error("failed to get block datah",
				slog.String("hash", currentHash))
			return nil, fmt.Errorf("failed to get block data: %v", err)
		}

		// DEVELOPERS NOTE:
		// If we get a nil block then that means we have reached the genesis
		// block so we can abort.
		if block == nil {
			break // Genesis block reached
		}

		for _, tx := range block.Transactions {
			if tx.From == address {
				balance.Sub(balance, tx.Value)
			}
			if tx.To == address {
				balance.Add(balance, tx.Value)
			}
		}

		if block.PreviousHash == "" {
			break // Genesis block reached
		}
		currentHash = block.PreviousHash
	}

	return balance, nil
}
