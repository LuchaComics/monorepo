package service

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/ethereum/go-ethereum/common"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/common/httperror"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/domain"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/usecase"
)

type TransferCoinService struct {
	config                                *config.Config
	logger                                *slog.Logger
	getAccountUseCase                     *usecase.GetAccountUseCase
	getWalletUseCase                      *usecase.GetWalletUseCase
	walletDecryptKeyUseCase               *usecase.WalletDecryptKeyUseCase
	broadcastMempoolTransactionDTOUseCase *usecase.BroadcastMempoolTransactionDTOUseCase
}

func NewTransferCoinService(
	cfg *config.Config,
	logger *slog.Logger,
	uc1 *usecase.GetAccountUseCase,
	uc2 *usecase.GetWalletUseCase,
	uc3 *usecase.WalletDecryptKeyUseCase,
	uc4 *usecase.BroadcastMempoolTransactionDTOUseCase,
) *TransferCoinService {
	return &TransferCoinService{cfg, logger, uc1, uc2, uc3, uc4}
}

func (s *TransferCoinService) Execute(
	ctx context.Context,
	fromAccountAddress *common.Address,
	accountWalletPassword string,
	to *common.Address,
	value uint64,
	data []byte,
) error {
	//
	// STEP 1: Validation.
	//

	e := make(map[string]string)
	if fromAccountAddress == nil {
		e["from_account_address"] = "missing value"
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
		s.logger.Warn("Failed validating create transaction parameters",
			slog.Any("error", e))
		return httperror.NewForBadRequest(&e)
	}

	//
	// STEP 2: Get the account and extract the wallet private/public key.
	//

	wallet, err := s.getWalletUseCase.Execute(fromAccountAddress)
	if err != nil {
		s.logger.Error("failed getting from database",
			slog.Any("from_account_address", fromAccountAddress),
			slog.Any("error", err))
		return fmt.Errorf("failed getting from database: %s", err)
	}
	if wallet == nil {
		s.logger.Error("failed getting from database",
			slog.Any("from_account_address", fromAccountAddress),
			slog.Any("error", "d.n.e."))
		return fmt.Errorf("failed getting from database: %s", "wallet d.n.e.")
	}

	key, err := s.walletDecryptKeyUseCase.Execute(wallet.Filepath, accountWalletPassword)
	if err != nil {
		s.logger.Error("failed getting key",
			slog.Any("from_account_address", fromAccountAddress),
			slog.Any("error", err))
		return fmt.Errorf("failed getting key: %s", err)
	}
	if key == nil {
		return fmt.Errorf("failed getting key: %s", "d.n.e.")
	}

	//
	// STEP 3:
	// Verify the account has enough balance before proceeding.
	//

	account, err := s.getAccountUseCase.Execute(fromAccountAddress)
	if err != nil {
		s.logger.Error("failed getting account",
			slog.Any("from_account_address", fromAccountAddress),
			slog.Any("error", err))
		return fmt.Errorf("failed getting account: %s", err)
	}
	if account == nil {
		return fmt.Errorf("failed getting account: %s", "d.n.e.")
	}
	if account.Balance <= value {
		s.logger.Warn("insufficient balance in account",
			slog.Any("account_addr", fromAccountAddress),
			slog.Any("account_balance", account.Balance),
			slog.Any("value", value))
		return fmt.Errorf("insufficient balance: %d", account.Balance)
	}

	//
	// STEP 4
	// Create our pending transaction and sign it with the accounts private key.
	//

	tx := &domain.Transaction{
		ChainID: s.config.Blockchain.ChainID,
		Nonce:   uint64(time.Now().Unix()),
		From:    wallet.Address,
		To:      to,
		Value:   value,
		Data:    data,
		Type:    domain.TransactionTypeCoin,
	}

	stx, signingErr := tx.Sign(key.PrivateKey)
	if signingErr != nil {
		s.logger.Debug("Failed to sign the transaction",
			slog.Any("error", signingErr))
		return signingErr
	}

	s.logger.Debug("Pending transaction signed successfully",
		slog.Any("chain_id", stx.ChainID),
		slog.Any("nonce", stx.Nonce),
		slog.Any("from", stx.From),
		slog.Any("to", stx.To),
		slog.Any("value", stx.Value),
		slog.Any("data", stx.Data),
		slog.Any("type", stx.Type),
		slog.Any("token_id", stx.TokenID),
		slog.Any("token_metadata_uri", stx.TokenMetadataURI),
		slog.Any("token_nonce", stx.TokenNonce),
		slog.Any("tx_sig_v", stx.V),
		slog.Any("tx_sig_r", stx.R),
		slog.Any("tx_sig_s", stx.S),
		slog.Uint64("tx_nonce", stx.Nonce))

	mempoolTx := &domain.MempoolTransaction{
		Transaction: stx.Transaction,
		V:           stx.V,
		R:           stx.R,
		S:           stx.S,
	}

	//
	// STEP 3
	// Send our pending signed transaction to our distributed mempool nodes
	// in the blochcian network.
	//

	if err := s.broadcastMempoolTransactionDTOUseCase.Execute(ctx, mempoolTx); err != nil {
		s.logger.Error("Failed to broadcast to the blockchain network",
			slog.Any("error", err))
		return err
	}

	s.logger.Info("Pending signed transaction submitted to blockchain",
		slog.Uint64("tx_nonce", stx.Nonce))

	return nil
}
