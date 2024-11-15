package service

import (
	"context"
	"fmt"
	"log/slog"
	"math/big"
	"strings"
	"time"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/common/httperror"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/domain"
	auth_usecase "github.com/LuchaComics/monorepo/cloud/comiccoin-authority/usecase"
	"github.com/ethereum/go-ethereum/common"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin-cli/usecase"
)

type TokenBurnService struct {
	logger                                                  *slog.Logger
	getAccountUseCase                                       *usecase.GetAccountUseCase
	getWalletUseCase                                        *usecase.GetWalletUseCase
	walletDecryptKeyUseCase                                 *usecase.WalletDecryptKeyUseCase
	getTokenUseCase                                         *usecase.GetTokenUseCase
	submitMempoolTransactionDTOToBlockchainAuthorityUseCase *auth_usecase.SubmitMempoolTransactionDTOToBlockchainAuthorityUseCase
}

func NewTokenBurnService(
	logger *slog.Logger,
	uc1 *usecase.GetAccountUseCase,
	uc2 *usecase.GetWalletUseCase,
	uc3 *usecase.WalletDecryptKeyUseCase,
	uc4 *usecase.GetTokenUseCase,
	uc5 *auth_usecase.SubmitMempoolTransactionDTOToBlockchainAuthorityUseCase,
) *TokenBurnService {
	return &TokenBurnService{logger, uc1, uc2, uc3, uc4, uc5}
}

func (s *TokenBurnService) Execute(
	ctx context.Context,
	chainID uint16,
	fromAccountAddress *common.Address,
	accountWalletPassword string,
	tokenID *big.Int,
) error {
	s.logger.Debug("Validating...",
		slog.Any("chain_id", chainID),
		slog.Any("from_account_address", fromAccountAddress),
		slog.Any("account_wallet_password", accountWalletPassword),
		slog.Any("tokenID", tokenID),
	)

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
	if tokenID == nil {
		e["token_id"] = "missing value"
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
	// Verify the account owns the token
	//

	account, err := s.getAccountUseCase.Execute(ctx, fromAccountAddress)
	if err != nil {
		s.logger.Error("failed getting account",
			slog.Any("from_account_address", fromAccountAddress),
			slog.Any("error", err))
		return fmt.Errorf("failed getting account: %s", err)
	}
	if account == nil {
		return fmt.Errorf("failed getting account: %s", "d.n.e.")
	}
	if strings.ToLower(account.Address.Hex()) != strings.ToLower(tok.Owner.Hex()) {
		s.logger.Warn("you do not own this token",
			slog.Any("account_addr", strings.ToLower(account.Address.Hex())),
			slog.Any("token_addr", strings.ToLower(tok.Owner.Hex())))
		return fmt.Errorf("you do not own the token: %v", strings.ToLower(tok.Owner.Hex()))
	}

	//
	// STEP 4
	// Create our pending transaction and sign it with the accounts private key.
	//

	// Burn an NFT by setting its owner to the burn address
	burnAddress := common.HexToAddress("0x0000000000000000000000000000000000000000")

	tx := &domain.Transaction{
		ChainID:          chainID,
		NonceBytes:       big.NewInt(time.Now().Unix()).Bytes(),
		From:             wallet.Address,
		To:               &burnAddress,
		Value:            0,
		Data:             []byte{},
		Type:             domain.TransactionTypeToken,
		TokenIDBytes:     tok.IDBytes,
		TokenMetadataURI: tok.MetadataURI,
		TokenNonceBytes:  tok.NonceBytes,
	}

	stx, signingErr := tx.Sign(key.PrivateKey)
	if signingErr != nil {
		s.logger.Debug("Failed to sign the transaction",
			slog.Any("error", signingErr))
		return signingErr
	}

	// Defensive Coding.
	if err := stx.Validate(chainID, true); err != nil {
		s.logger.Debug("Failed to validate signature of the signed transaction",
			slog.Any("error", signingErr))
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

	mempoolTx := &domain.MempoolTransaction{
		ID:                primitive.NewObjectID(),
		SignedTransaction: stx,
	}

	// Defensive Coding.
	if err := mempoolTx.Validate(chainID, true); err != nil {
		s.logger.Debug("Failed to validate signature of mempool transaction",
			slog.Any("error", signingErr))
		return signingErr
	}

	s.logger.Debug("Mempool transaction ready for submission",
		slog.Any("Transaction", stx.Transaction),
		slog.Any("tx_sig_v_bytes", stx.VBytes),
		slog.Any("tx_sig_r_bytes", stx.RBytes),
		slog.Any("tx_sig_s_bytes", stx.SBytes))

	//
	// STEP 3
	// Send our pending signed transaction to the authority's mempool to wait
	// in a queue to be processed.
	//

	dto := mempoolTx.ToDTO()

	if err := s.submitMempoolTransactionDTOToBlockchainAuthorityUseCase.Execute(ctx, dto); err != nil {
		s.logger.Error("Failed to broadcast to the blockchain authority",
			slog.Any("error", err))
		return err
	}

	s.logger.Info("Pending signed transaction for token burn submitted to the blockchain authority",
		slog.Any("tx_nonce", stx.GetNonce()))

	return nil
}
