package service

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"math/big"
	"strings"
	"time"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/common/httperror"
	sstring "github.com/LuchaComics/monorepo/cloud/comiccoin-authority/common/security/securestring"
	auth_domain "github.com/LuchaComics/monorepo/cloud/comiccoin-authority/domain"
	auth_usecase "github.com/LuchaComics/monorepo/cloud/comiccoin-authority/usecase"
	"github.com/ethereum/go-ethereum/common"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/domain"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/usecase"
)

type TokenTransferService struct {
	logger                                                  *slog.Logger
	storageTransactionOpenUseCase                           *usecase.StorageTransactionOpenUseCase
	storageTransactionCommitUseCase                         *usecase.StorageTransactionCommitUseCase
	storageTransactionDiscardUseCase                        *usecase.StorageTransactionDiscardUseCase
	listPendingSignedTransactionUseCase                     *usecase.ListPendingSignedTransactionUseCase
	getGenesisBlockDataUseCase                              *usecase.GetGenesisBlockDataUseCase
	upsertPendingSignedTransactionUseCase                   *usecase.UpsertPendingSignedTransactionUseCase
	getAccountUseCase                                       *usecase.GetAccountUseCase
	getWalletUseCase                                        *usecase.GetWalletUseCase
	walletDecryptKeyUseCase                                 *usecase.WalletDecryptKeyUseCase
	getTokenUseCase                                         *usecase.GetTokenUseCase
	submitMempoolTransactionDTOToBlockchainAuthorityUseCase *auth_usecase.SubmitMempoolTransactionDTOToBlockchainAuthorityUseCase
}

func NewTokenTransferService(
	logger *slog.Logger,
	uc1 *usecase.StorageTransactionOpenUseCase,
	uc2 *usecase.StorageTransactionCommitUseCase,
	uc3 *usecase.StorageTransactionDiscardUseCase,
	uc4 *usecase.ListPendingSignedTransactionUseCase,
	uc5 *usecase.GetGenesisBlockDataUseCase,
	uc6 *usecase.UpsertPendingSignedTransactionUseCase,
	uc7 *usecase.GetAccountUseCase,
	uc8 *usecase.GetWalletUseCase,
	uc9 *usecase.WalletDecryptKeyUseCase,
	uc10 *usecase.GetTokenUseCase,
	uc11 *auth_usecase.SubmitMempoolTransactionDTOToBlockchainAuthorityUseCase,
) *TokenTransferService {
	return &TokenTransferService{logger, uc1, uc2, uc3, uc4, uc5, uc6, uc7, uc8, uc9, uc10, uc11}
}

func (s *TokenTransferService) Execute(
	ctx context.Context,
	chainID uint16,
	fromAccountAddress *common.Address,
	accountWalletPassword *sstring.SecureString,
	to *common.Address,
	tokenID *big.Int,
) error {
	s.logger.Debug("Validating...",
		slog.Any("chain_id", chainID),
		slog.Any("from_account_address", fromAccountAddress),
		slog.Any("account_wallet_password", accountWalletPassword),
		slog.Any("to", to),
		slog.Any("tokenID", tokenID),
	)

	//
	// STEP 1: Validation.
	//

	e := make(map[string]string)
	if fromAccountAddress == nil {
		e["from_account_address"] = "missing value"
	}
	if accountWalletPassword == nil {
		e["account_wallet_password"] = "missing value"
	}
	if to == nil {
		e["to"] = "missing value"
	}
	if tokenID == nil {
		e["token_id"] = "missing value"
	}
	pstxs, err := s.listPendingSignedTransactionUseCase.Execute(ctx)
	if err != nil {
		s.logger.Debug("Failed listing pending signed transactions", slog.Any("error", err))
		return err
	}
	if pstxs != nil {
		if len(pstxs) > 0 {
			e["message"] = "already has a pending transaction - please wait for authority to complete request"
		}
	}
	if len(e) != 0 {
		s.logger.Warn("Failed validating create transaction parameters",
			slog.Any("error", e))
		return httperror.NewForBadRequest(&e)
	}

	//
	// STEP 2:
	// Get related records.
	//

	genesis, err := s.getGenesisBlockDataUseCase.Execute(ctx, chainID)
	if err != nil {
		s.logger.Error("failed getting genesis from database",
			slog.Any("chain_id", chainID),
			slog.Any("error", err))
		return err
	}
	if genesis == nil {
		s.logger.Error("failed getting genesis from database",
			slog.Any("chain_id", chainID),
			slog.Any("error", "d.n.e."))
		return fmt.Errorf("failed getting genesis block from database: %s", "genesis d.n.e.")
	}
	txFee := genesis.Header.TransactionFee

	wallet, err := s.getWalletUseCase.Execute(ctx, fromAccountAddress)
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

	key, err := s.walletDecryptKeyUseCase.Execute(ctx, wallet.KeystoreBytes, accountWalletPassword)
	if err != nil {
		s.logger.Error("failed getting key",
			slog.Any("from_account_address", fromAccountAddress),
			slog.Any("error", err))
		return fmt.Errorf("failed getting key: %s", err)
	}
	if key == nil {
		return fmt.Errorf("failed getting key: %s", "d.n.e.")
	}

	tok, err := s.getTokenUseCase.Execute(ctx, tokenID)
	if err != nil {
		s.logger.Error("failed getting token from database",
			slog.Any("from_account_address", fromAccountAddress),
			slog.Any("error", err))
		return fmt.Errorf("failed getting token from database: %s", err)
	}
	if tok == nil {
		s.logger.Error("failed getting token from database",
			slog.Any("from_account_address", fromAccountAddress),
			slog.Any("error", "d.n.e."))
		return fmt.Errorf("failed getting token from database: %s", "token d.n.e.")
	}

	//
	// STEP 3:
	// Verify the account owns the token and has enough balance before proceeding.
	//

	if err := s.storageTransactionOpenUseCase.Execute(); err != nil {
		s.storageTransactionDiscardUseCase.Execute()
		log.Fatalf("Failed to open storage transaction: %v\n", err)
	}

	account, err := s.getAccountUseCase.Execute(ctx, fromAccountAddress)
	if err != nil {
		s.logger.Error("failed getting account",
			slog.Any("from_account_address", fromAccountAddress),
			slog.Any("error", err))
		s.storageTransactionDiscardUseCase.Execute()
		return fmt.Errorf("failed getting account: %s", err)
	}
	if account == nil {
		s.storageTransactionDiscardUseCase.Execute()
		return fmt.Errorf("failed getting account: %s", "d.n.e.")
	}

	// Verify ownership of NFT.
	if strings.ToLower(account.Address.Hex()) != strings.ToLower(tok.Owner.Hex()) {
		s.logger.Warn("you do not own this token",
			slog.Any("account_addr", strings.ToLower(account.Address.Hex())),
			slog.Any("token_addr", strings.ToLower(tok.Owner.Hex())))
		s.storageTransactionDiscardUseCase.Execute()
		return fmt.Errorf("you do not own the token: %v", strings.ToLower(tok.Owner.Hex()))
	}

	// Verify balance covers fees in account
	if account.Balance < txFee {
		s.logger.Warn("insufficient balance in account",
			slog.Any("account_addr", fromAccountAddress),
			slog.Any("account_balance", account.Balance),
			slog.Any("fee", txFee),
			slog.Any("value", txFee),
			slog.Any("total", txFee))
		s.storageTransactionDiscardUseCase.Execute()
		return fmt.Errorf("insufficient balance: %d", account.Balance)
	}

	//
	// STEP 4
	// Create our pending transaction and sign it with the accounts private key.
	//

	tx := &auth_domain.Transaction{
		ChainID:          chainID,
		NonceBytes:       big.NewInt(time.Now().Unix()).Bytes(),
		From:             wallet.Address,
		To:               to,
		Value:            txFee, // Users pay transaction fee for transfering NFTs.
		Data:             []byte{},
		Type:             auth_domain.TransactionTypeToken,
		TokenIDBytes:     tok.IDBytes,
		TokenMetadataURI: tok.MetadataURI,
		TokenNonceBytes:  tok.NonceBytes,
	}

	stx, signingErr := tx.Sign(key.PrivateKey)
	if signingErr != nil {
		s.logger.Debug("Failed to sign the transaction",
			slog.Any("error", signingErr))
		s.storageTransactionDiscardUseCase.Execute()
		return signingErr
	}

	// Defensive Coding.
	if err := stx.Validate(chainID, true); err != nil {
		s.logger.Debug("Failed to validate signature of the signed transaction",
			slog.Any("error", signingErr))
		s.storageTransactionDiscardUseCase.Execute()
		return signingErr
	}

	s.logger.Debug("Transaction signed successfully",
		slog.Any("chain_id", stx.ChainID),
		slog.Any("nonce", stx.GetNonce()),
		slog.Any("from", stx.From),
		slog.Any("to", stx.To),
		slog.Any("value", stx.Value),
		slog.Any("data", stx.Data),
		slog.Any("type", stx.Type),
		slog.Any("token_id", stx.GetTokenID()),
		slog.Any("token_metadata_uri", stx.TokenMetadataURI),
		slog.Any("token_nonce", stx.GetTokenNonce()),
		slog.Any("tx_sig_v_bytes", stx.VBytes),
		slog.Any("tx_sig_r_bytes", stx.RBytes),
		slog.Any("tx_sig_s_bytes", stx.SBytes),
		slog.Any("tx_nonce", stx.GetNonce()))

	mempoolTx := &auth_domain.MempoolTransaction{
		ID:                primitive.NewObjectID(),
		SignedTransaction: stx,
	}

	// Defensive Coding.
	if err := mempoolTx.Validate(chainID, true); err != nil {
		s.logger.Debug("Failed to validate signature of mempool transaction",
			slog.Any("error", signingErr))
		s.storageTransactionDiscardUseCase.Execute()
		return signingErr
	}

	s.logger.Debug("Mempool transaction ready for submission",
		slog.Any("Transaction", stx.Transaction),
		slog.Any("tx_sig_v_bytes", stx.VBytes),
		slog.Any("tx_sig_r_bytes", stx.RBytes),
		slog.Any("tx_sig_s_bytes", stx.SBytes))

	//
	// STEP 5: Save as pending signed transaction to keep track of completion.
	//

	pstx := domain.SignedTransactionToPendingSignedTransaction(&stx)
	if err := s.upsertPendingSignedTransactionUseCase.Execute(ctx, pstx); err != nil {
		s.logger.Debug("Failed saving pending signed transaction",
			slog.Any("error", signingErr))
		s.storageTransactionDiscardUseCase.Execute()
		return err
	}

	//
	// STEP 6
	// Send our pending signed transaction to the authority's mempool to wait
	// in a queue to be processed.
	//

	dto := mempoolTx.ToDTO()

	if err := s.submitMempoolTransactionDTOToBlockchainAuthorityUseCase.Execute(ctx, dto); err != nil {
		s.logger.Error("Failed to broadcast to the blockchain authority",
			slog.Any("error", err))
		s.storageTransactionDiscardUseCase.Execute()
		return err
	}

	s.logger.Info("Pending signed transaction for token transfer submitted to the blockchain authority",
		slog.Any("tx_nonce", stx.GetNonce()))

	if err := s.storageTransactionCommitUseCase.Execute(); err != nil {
		s.storageTransactionDiscardUseCase.Execute()
		log.Fatalf("Failed to open storage transaction: %v\n", err)
	}

	s.logger.Info("Pending signed transaction saved to local storage",
		slog.Any("tx_nonce", stx.GetNonce()))

	return nil
}
