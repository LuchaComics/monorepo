package controller

import (
	"context"
	"fmt"
	"log/slog"
	"math/big"
	"net/http"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/utils/httperror"
)

type BlockchainBalanceResponseIDO struct {
	Amount *big.Int `json:"amount"`
}

func (impl *blockchainControllerImpl) GetBalanceByAccountName(ctx context.Context, accountName string) (*BlockchainBalanceResponseIDO, error) {
	account, err := impl.accountStorer.GetByName(ctx, accountName)
	if err != nil {
		impl.logger.Error("failed getting account",
			slog.String("account_name", accountName),
			slog.Any("error", err))
		return nil, httperror.NewForSingleField(http.StatusBadRequest, "message", fmt.Sprintf("error getting account: %v", err))
	}
	if account == nil {
		impl.logger.Error("failed getting account as account d.n.e.",
			slog.Any("account_name", accountName))
		return nil, httperror.NewForSingleField(http.StatusBadRequest, "message", fmt.Sprintf("does not exist for name: %v", accountName))

	}

	balance := new(big.Int)
	currentHash, err := impl.lastHashStorer.Get(ctx)
	if err != nil {
		impl.logger.Error("failed to get last hash",
			slog.Any("account_name", accountName),
			slog.Any("error", err))
		return nil, httperror.NewForSingleField(http.StatusBadRequest, "message", fmt.Sprintf("failed to get last hash: %v", err))
	}
	impl.logger.Debug("lookup last hash",
		slog.Any("account_name", accountName),
		slog.String("last_hash", currentHash))

	// Iterate through all the blocks.
	for {
		block, err := impl.blockStorer.GetByHash(ctx, currentHash)
		if err != nil {
			impl.logger.Error("failed to get block datah",
				slog.String("hash", currentHash))
			return nil, httperror.NewForSingleField(http.StatusBadRequest, "message", fmt.Sprintf("failed to get block data: %v", err))
		}

		// DEVELOPERS NOTE:
		// If we get a nil block then that means we have reached the genesis
		// block so we can abort.
		if block == nil {
			break // Genesis block reached
		}

		for _, tx := range block.Transactions {
			if tx.From == account.WalletAddress {
				balance.Sub(balance, tx.Value)
			}
			if tx.To == account.WalletAddress {
				balance.Add(balance, tx.Value)
			}
		}

		if block.PreviousHash == "" {
			break // Genesis block reached
		}
		currentHash = block.PreviousHash
	}

	return &BlockchainBalanceResponseIDO{
		Amount: balance,
	}, nil
}
