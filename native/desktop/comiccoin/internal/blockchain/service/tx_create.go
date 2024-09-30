package service

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/ethereum/go-ethereum/common"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/domain"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/usecase"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/pkg/httperror"
)

type CreateTransactionService struct {
	config                               *config.Config
	logger                               *slog.Logger
	getAccountUseCase                    *usecase.GetAccountUseCase
	accountDecryptKeyUseCase             *usecase.AccountDecryptKeyUseCase
	broadcastSignedTransactionDTOUseCase *usecase.BroadcastSignedTransactionDTOUseCase
}

func NewCreateTransactionService(
	cfg *config.Config,
	logger *slog.Logger,
	uc1 *usecase.GetAccountUseCase,
	uc2 *usecase.AccountDecryptKeyUseCase,
	uc3 *usecase.BroadcastSignedTransactionDTOUseCase,
) *CreateTransactionService {
	return &CreateTransactionService{cfg, logger, uc1, uc2, uc3}
}

func (s *CreateTransactionService) Execute(
	ctx context.Context,
	fromAccountID string,
	accountWalletPassword string,
	to *common.Address,
	value uint64,
	data []byte,
) error {
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
		return httperror.NewForBadRequest(&e)
	}

	//
	// STEP 2: Get the account and extract the wallet private/public key.
	//

	account, err := s.getAccountUseCase.Execute(fromAccountID)
	if err != nil {
		s.logger.Error("failed getting from database",
			slog.Any("from_account_id", fromAccountID),
			slog.Any("error", err))
		return fmt.Errorf("failed getting from database: %s", err)
	}
	if account == nil {
		return fmt.Errorf("failed getting from database: %s", "d.n.e.")
	}

	key, err := s.accountDecryptKeyUseCase.Execute(account.WalletFilepath, accountWalletPassword)
	if err != nil {
		s.logger.Error("failed getting key",
			slog.Any("from_account_id", fromAccountID),
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

	//TODO: IMPL.

	//
	// STEP 4
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
		return signingErr
	}

	s.logger.Debug("Pending transaction signed successfully",
		slog.Uint64("nonce", stx.Nonce))

	//
	// STEP 3
	// Send our pending signed transaction to our distributed mempool nodes
	// in the blochcian network.
	//

	if err := s.broadcastSignedTransactionDTOUseCase.Execute(ctx, &stx); err != nil {
		s.logger.Error("Failed to broadcast to the blockchain network",
			slog.Any("error", err))
		return err
	}

	s.logger.Debug("Pending transaction submitted to blockchain!",
		slog.Uint64("nonce", stx.Nonce))

	return nil
}
