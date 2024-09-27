package controller

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/ethereum/go-ethereum/common"

	pt_ds "github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/app/signedtransaction/datastore"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/config/constants"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/utils/httperror"
)

type BlockchainSubmitRequestIDO struct {
	// Name of the account
	FromAccountName string `json:"from_account_name"`

	AccountWalletPassword string `json:"account_wallet_password"`

	// Recipient’s public key
	To string `json:"to"`

	// Value is amount of coins being transferred
	Value uint64 `json:"value"`

	// Data is any NFT related data attached
	Data []byte `json:"data"`
}

type BlockchainSubmitResponseIDO struct {
}

func (impl *blockchainControllerImpl) validateSubmitRequest(ctx context.Context, dirtyData *BlockchainSubmitRequestIDO) error {
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
		if dirtyData.Value == 0 {
			e["value"] = "missing value"
		}

	}

	if len(e) != 0 {
		impl.logger.Debug("Failed creating new wallet",
			slog.Any("e", e))
		return httperror.NewForBadRequest(&e)
	}
	return nil
}

func (impl *blockchainControllerImpl) Submit(ctx context.Context, req *BlockchainSubmitRequestIDO) (*BlockchainSubmitResponseIDO, error) {
	if err := impl.validateSubmitRequest(ctx, req); err != nil {
		impl.logger.Warn("Validation failed",
			slog.Any("error", err))
		return nil, err
	}

	if isConnected := impl.p2pPubSubBroker.IsSubscriberConnectedToNetwork(ctx, constants.PubSubMempoolTopicName); !isConnected {
		impl.logger.Error("Not connected to distributed network")
		return nil, httperror.NewForServiceUnavailableWithSingleField("message", "Not connected to distributed network")
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
	// Create our signed transaction and sign it with the accounts private key.
	//

	pt := &pt_ds.Transaction{
		ChainID: impl.config.Blockchain.ChainID,
		Nonce:   uint64(time.Now().Unix()),
		From:    account.WalletAddress,
		To:      common.HexToAddress(req.To),
		Value:   req.Value,
		Data:    req.Data,
	}

	signedTransaction, signingErr := pt.Sign(accountKey.PrivateKey)
	if signingErr != nil {
		impl.logger.Debug("Failed to sign the signed transaction",
			slog.Any("error", signingErr))
		return nil, signingErr
	}

	impl.logger.Debug("Pending transaction signed successfully",
		slog.Uint64("nonce", pt.Nonce))

	//
	// STEP 3
	// Submit to the blockchain network to be processed by the consensus
	// mechanism and verification to be submitted in the blockchain.
	// We are using `publish-subscribe` pattern with a `message queue` which
	// will `publish` the message to the broker so the broker will pass it
	// along to the subscriber which will submit this to the peer-to-peer
	// network.
	//

	ptBytes, err := signedTransaction.Serialize()
	if err != nil {
		impl.logger.Error("Failed to serialize our signed transaction",
			slog.Any("error", err))
		return nil, err
	}

	// Send our pending signed transaction to our distributed mempool nodes
	// in the blochcian network. The `mempool` topic is used to
	// send our signed pending transcation to the actively running in background
	// mempool node subscriber
	if err := impl.p2pPubSubBroker.Publish(ctx, constants.PubSubMempoolTopicName, ptBytes); err != nil {
		impl.logger.Error("Failed to publish",
			slog.Any("error", err))
		return nil, err
	}

	impl.logger.Debug("Pending transaction submitted to blockchain!",
		slog.Uint64("nonce", pt.Nonce))

	return nil, nil

}
