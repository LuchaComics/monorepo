package service

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/domain"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/usecase"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/pkg/httperror"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
)

type CreateTransactionService struct {
	config                   *config.Config
	logger                   *slog.Logger
	getAccountUseCase        *usecase.GetAccountUseCase
	accountDecryptKeyUseCase *usecase.AccountDecryptKeyUseCase
}

func NewCreateTransactionService(
	cfg *config.Config,
	logger *slog.Logger,
	uc1 *usecase.GetAccountUseCase,
	uc2 *usecase.AccountDecryptKeyUseCase,
) *CreateTransactionService {
	return &CreateTransactionService{cfg, logger, uc1, uc2}
}

func (s *CreateTransactionService) Execute(
	fromAccountID, accountWalletPassword string,
	to *common.Address,
	value uint64,
	data []byte,
) (*keystore.Key, error) {
	//
	// STEP 1: Validation.
	//

	e := make(map[string]string)
	if fromAccountID == "" {
		e["from_account_id"] = "missing value"
	}
	if accountWalletPassword == "" {
		e["account_wallet_password"] = "missing value"
	}
	if to == nil {
		e["to"] = "missing value"
	}
	if value == 0 {
		e["value"] = "missing value"
	}
	if len(e) != 0 {
		s.logger.Warn("Failed creating new account",
			slog.Any("error", e))
		return nil, httperror.NewForBadRequest(&e)
	}

	//
	// STEP 2: Get the account and extract the wallet private/public key.
	//

	account, err := s.getAccountUseCase.Execute(fromAccountID)
	if err != nil {
		s.logger.Error("failed getting from database",
			slog.Any("from_account_id", fromAccountID),
			slog.Any("error", err))
		return nil, fmt.Errorf("failed getting from database: %s", err)
	}
	if account == nil {
		return nil, fmt.Errorf("failed getting from database: %s", "d.n.e.")
	}

	key, err := s.accountDecryptKeyUseCase.Execute(account.WalletFilepath, accountWalletPassword)
	if err != nil {
		s.logger.Error("failed getting key",
			slog.Any("from_account_id", fromAccountID),
			slog.Any("error", err))
		return nil, fmt.Errorf("failed getting key: %s", err)
	}
	if key == nil {
		return nil, fmt.Errorf("failed getting key: %s", "d.n.e.")
	}

	//
	// STEP 2
	// Create our pending transaction and sign it with the accounts private key.
	//

	tx := &domain.Transaction{
		ChainID: s.config.Blockchain.ChainID,
		Nonce:   uint64(time.Now().Unix()),
		From:    account.WalletAddress,
		To:      to,
		Value:   value,
		Data:    data,
	}

	stx, signingErr := tx.Sign(key.PrivateKey)
	if signingErr != nil {
		s.logger.Debug("Failed to sign the signed transaction",
			slog.Any("error", signingErr))
		return nil, signingErr
	}

	s.logger.Debug("Pending transaction signed successfully",
		slog.Uint64("nonce", stx.Nonce))

	//TODO: IMPL.
	// //
	// // STEP 3
	// // Submit to the blockchain network to be processed by the consensus
	// // mechanism and verification to be submitted in the blockchain.
	// // We are using `publish-subscribe` pattern with a `message queue` which
	// // will `publish` the message to the broker so the broker will pass it
	// // along to the subscriber which will submit this to the peer-to-peer
	// // network.
	// //
	//
	// ptBytes, err := signedTransaction.Serialize()
	// if err != nil {
	// 	impl.logger.Error("Failed to serialize our signed transaction",
	// 		slog.Any("error", err))
	// 	return nil, err
	// }
	//
	// // Send our pending signed transaction to our distributed mempool nodes
	// // in the blochcian network. The `mempool` topic is used to
	// // send our signed pending transcation to the actively running in background
	// // mempool node subscriber
	// if err := impl.p2pPubSubBroker.Publish(ctx, constants.PubSubMempoolTopicName, ptBytes); err != nil {
	// 	impl.logger.Error("Failed to publish",
	// 		slog.Any("error", err))
	// 	return nil, err
	// }
	//
	// impl.logger.Debug("Pending transaction submitted to blockchain!",
	// 	slog.Uint64("nonce", pt.Nonce))

	return nil, nil
}
