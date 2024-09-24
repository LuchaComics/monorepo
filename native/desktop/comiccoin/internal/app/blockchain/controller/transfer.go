package controller

import (
	"context"
	"fmt"
	"log/slog"
	"math/big"

	pt_ds "github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/app/pendingtransaction/datastore"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/utils/httperror"
)

type BlockchainTransferRequestIDO struct {
	// Name of the account
	FromAccountName string `json:"from_account_name"`

	AccountWalletPassword string `json:"account_wallet_password"`

	// Recipient’s public key
	To string `json:"to"`

	// Amount of coins being transferred
	Amount *big.Int `json:"amount"`

	// Data is any NFT related data attached
	Data []byte `json:"data"`
}

type BlockchainTransferResponseIDO struct {
}

func (impl *blockchainControllerImpl) validateTransferRequest(ctx context.Context, dirtyData *BlockchainTransferRequestIDO) error {
	e := make(map[string]string)

	if dirtyData == nil {
		e["from_account_name"] = "missing value"
		e["to"] = "missing value"
		e["amount"] = "missing value"
	} else {
		if dirtyData.FromAccountName == "" {
			e["from_account_name"] = "missing value"
		} else {
			account, err := impl.accountStorer.GetByName(context.Background(), dirtyData.FromAccountName)
			if err != nil {
				e["from_account_name"] = fmt.Sprintf("failed getting account: %v", err)
			}
			if account == nil {
				e["from_account_name"] = "account does not exist"
			}

			//TODO: Check if account has enough balance
		}
		if dirtyData.AccountWalletPassword == "" {
			e["account_wallet_password"] = "missing value"
		}
		if dirtyData.To == "" {
			e["to"] = "missing value"
		}
		if dirtyData.Amount == nil {
			e["amount"] = "missing value"
		}

	}

	if len(e) != 0 {
		impl.logger.Debug("Failed creating new wallet",
			slog.Any("e", e))
		return httperror.NewForBadRequest(&e)
	}
	return nil
}

func (impl *blockchainControllerImpl) Transfer(ctx context.Context, req *BlockchainTransferRequestIDO) (*BlockchainTransferResponseIDO, error) {
	if err := impl.validateTransferRequest(ctx, req); err != nil {
		return nil, err
	}
	impl.logger.Debug("Submitting transfer request to the blockchain network...",
		slog.Any("req", req))

	//
	// STEP 1.
	// Get all the related records.
	//

	account, _ := impl.accountStorer.GetByName(ctx, req.FromAccountName)
	accountKey, err := impl.accountStorer.GetKeyByNameAndPassword(ctx, req.FromAccountName, req.AccountWalletPassword)
	if err != nil {
		impl.logger.Debug("Failed creating new wallet",
			slog.Any("error", err))
		return nil, err
	}
	if accountKey == nil {
		return nil, httperror.NewForBadRequestWithSingleField("account", "key does not exist")
	}

	pt := &pt_ds.PendingTransaction{} // Continue later here

	return nil, nil

}
