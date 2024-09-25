package controller

import (
	"context"
	"fmt"
	"log/slog"
	"math/big"

	"github.com/ethereum/go-ethereum/common"

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
		impl.logger.Warn("Validation failed",
			slog.Any("error", err))
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
		impl.logger.Error("Failed creating new wallet",
			slog.Any("error", err))
		return nil, err
	}
	if accountKey == nil {
		impl.logger.Error("Key does not exist")
		return nil, httperror.NewForBadRequestWithSingleField("account", "key does not exist")
	}

	//
	// STEP 2
	// Create our pending transaction and sign it with the accounts private key.
	//

	pt := &pt_ds.PendingTransaction{
		ID:     impl.uuid.NewUUID("pending_transaction"),
		From:   account.WalletAddress,
		To:     common.HexToAddress(req.To),
		Amount: req.Amount,
		Data:   req.Data,
	}
	if signingErr := pt.Sign(accountKey.PrivateKey); signingErr != nil {
		impl.logger.Debug("Failed to sign the pending transaction",
			slog.Any("error", signingErr))
		return nil, signingErr
	}

	//
	// STEP 3
	// Save to our database.
	//

	if insertErr := impl.pendingTransactionStorer.Insert(ctx, pt); insertErr != nil {
		impl.logger.Debug("Failed to insert the pending transaction into the database",
			slog.Any("error", insertErr))
		return nil, insertErr
	}

	//
	// STEP 4
	// Submit to the blockchain network to be processed by the consensus
	// mechanism and verification to be submitted in the blockchain.
	// We are using `publish-subscribe` pattern with a `message queue` which
	// will `publish` the message to the broker so the broker will pass it
	// along to the subscriber which will submit this to the peer-to-peer
	// network.
	//

	ptBytes, err := pt.Serialize()
	if err != nil {
		impl.logger.Error("Failed to serialize our pending transaction",
			slog.Any("error", err))
		return nil, err
	}
	impl.messageQueueBroker.Publish("pending-transactions", ptBytes)

	return nil, nil

}
