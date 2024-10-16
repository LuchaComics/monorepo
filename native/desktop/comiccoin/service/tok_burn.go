package service

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/ethereum/go-ethereum/common"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/domain"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/usecase"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/common/httperror"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/common/kmutexutil"
)

// Service represents token owners burning the token
type BurnTokenService struct {
	config                                *config.Config
	logger                                *slog.Logger
	kmutex                                kmutexutil.KMutexProvider
	getWalletUseCase                      *usecase.GetWalletUseCase
	walletDecryptKeyUseCase               *usecase.WalletDecryptKeyUseCase
	getTokenUseCase                       *usecase.GetTokenUseCase
	broadcastMempoolTransactionDTOUseCase *usecase.BroadcastMempoolTransactionDTOUseCase
}

func NewBurnTokenService(
	cfg *config.Config,
	logger *slog.Logger,
	kmutex kmutexutil.KMutexProvider,
	uc1 *usecase.GetWalletUseCase,
	uc2 *usecase.WalletDecryptKeyUseCase,
	uc3 *usecase.GetTokenUseCase,
	uc4 *usecase.BroadcastMempoolTransactionDTOUseCase,
) *BurnTokenService {
	return &BurnTokenService{cfg, logger, kmutex, uc1, uc2, uc3, uc4}
}

func (s *BurnTokenService) Execute(
	ctx context.Context,
	tokenOwnerAddr *common.Address,
	tokenOwnerPassword string,
	tokenID uint64,
) error {
	// Lock the mining service until it has completed executing (or errored).
	s.kmutex.Acquire("token-burning")
	defer s.kmutex.Release("token-burning")

	//
	// STEP 1:
	// Validation.
	//

	e := make(map[string]string)
	if tokenOwnerAddr == nil {
		e["token_owner_address"] = "missing value"
	}
	if tokenOwnerPassword == "" {
		e["token_owner_password"] = "missing value"
	}
	if len(e) != 0 {
		s.logger.Warn("Failed validating token burn parameters",
			slog.Any("error", e))
		return httperror.NewForBadRequest(&e)
	}

	//
	// STEP 2:
	// Get the account and extract the wallet private/public key.
	//

	wallet, err := s.getWalletUseCase.Execute(tokenOwnerAddr)
	if err != nil {
		s.logger.Error("failed getting from database",
			slog.Any("error", err))
		return fmt.Errorf("failed getting from database: %s", err)
	}
	if wallet == nil {
		s.logger.Error("failed getting from database",
			slog.Any("error", "d.n.e."))
		return fmt.Errorf("failed getting from database: %s", "wallet d.n.e.")
	}

	key, err := s.walletDecryptKeyUseCase.Execute(wallet.Filepath, tokenOwnerPassword)
	if err != nil {
		s.logger.Error("failed getting key",
			slog.Any("error", err))
		return fmt.Errorf("failed getting key: %s", err)
	}
	if key == nil {
		return fmt.Errorf("failed getting key: %s", "d.n.e.")
	}

	//
	// STEP 3:
	// Get the token for the particular token ID.
	//

	token, err := s.getTokenUseCase.Execute(tokenID)
	if err != nil {
		s.logger.Error("failed getting token",
			slog.Any("error", err))
		return fmt.Errorf("failed getting token: %s", err)
	}

	// Defensive code.
	if token == nil {
		s.logger.Warn("failed getting token",
			slog.Any("token_id", tokenID),
			slog.Any("error", "token does not exist"))
		return fmt.Errorf("failed getting token: does not exist for ID: %v", tokenID)
	}

	//
	// STEP 3:
	// Verify the account owns the token
	//

	if tokenOwnerAddr.Hex() != token.Owner.Hex() {
		s.logger.Warn("permission failed",
			slog.Any("token_id", tokenID))
		return fmt.Errorf("permission denied: token address is %v but your address is %v", token.Owner.Hex(), tokenOwnerAddr.Hex())
	}

	//
	// STEP 4:
	// Increment token `nonce` - this is very important as it tells the
	// blockchain that we are commiting a transaction and hence the miner will
	// execute the burn. If we do not increment the nonce then no
	// transaction happens!
	//

	token.Nonce += 1

	//
	// STEP 5
	// Create our pending transaction and sign it with the accounts private key.
	//

	// Burn an NFT by setting its owner to the burn address
	burnAddress := common.HexToAddress("0x0000000000000000000000000000000000000000")

	tx := &domain.Transaction{
		ChainID:          s.config.Blockchain.ChainID,
		Nonce:            uint64(time.Now().Unix()),
		From:             tokenOwnerAddr,
		To:               &burnAddress,
		Value:            0, // Token have no value!
		Tip:              0,
		Data:             make([]byte, 0),
		Type:             domain.TransactionTypeToken,
		TokenID:          token.ID,
		TokenMetadataURI: token.MetadataURI,
		TokenNonce:       token.Nonce,
	}

	stx, signingErr := tx.Sign(key.PrivateKey)
	if signingErr != nil {
		s.logger.Debug("Failed to sign",
			slog.Any("error", signingErr))
		return signingErr
	}

	s.logger.Debug("Pending token burn transaction signed successfully",
		slog.Uint64("tx_token_id", stx.TokenID))

	//
	// STEP 6
	// Send our pending signed transaction to our distributed mempool nodes
	// in the blochcian network.
	//

	mempoolTx := &domain.MempoolTransaction{
		Transaction: stx.Transaction,
		V:           stx.V,
		R:           stx.R,
		S:           stx.S,
	}

	if err := s.broadcastMempoolTransactionDTOUseCase.Execute(ctx, mempoolTx); err != nil {
		s.logger.Error("Failed to broadcast to the blockchain network",
			slog.Any("error", err))
		return err
	}

	s.logger.Info("Pending signed token transaction submitted to blockchain",
		slog.Uint64("tx_token_id", stx.TokenID))

	return nil
}
