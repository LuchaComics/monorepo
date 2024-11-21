package usecase

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/ethereum/go-ethereum/accounts/keystore"

	pkgkeystore "github.com/LuchaComics/monorepo/cloud/comiccoin-authority/common/blockchain/keystore"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/common/httperror"
	sstring "github.com/LuchaComics/monorepo/cloud/comiccoin-authority/common/security/securestring"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/domain"
)

type AccountDecryptKeyUseCase struct {
	logger   *slog.Logger
	keystore pkgkeystore.KeystoreAdapter
	repo     domain.AccountRepository
}

func NewAccountDecryptKeyUseCase(
	logger *slog.Logger,
	keystore pkgkeystore.KeystoreAdapter,
	repo domain.AccountRepository,
) *AccountDecryptKeyUseCase {
	return &AccountDecryptKeyUseCase{logger, keystore, repo}
}

func (uc *AccountDecryptKeyUseCase) Execute(ctx context.Context, walletKeystoreBytes []byte, walletPassword *sstring.SecureString) (*keystore.Key, error) {
	//
	// STEP 1: Validation.
	//

	e := make(map[string]string)
	if walletKeystoreBytes == nil {
		e["wallet_keystore_bytes"] = "missing value"
	}
	if walletPassword == nil {
		e["wallet_password"] = "missing value"
	}
	if len(e) != 0 {
		uc.logger.Warn("Failed reading account key",
			slog.Any("error", e))
		return nil, httperror.NewForBadRequest(&e)
	}

	//
	// STEP 2: Decrypt key
	//

	key, err := uc.keystore.OpenWallet(walletKeystoreBytes, walletPassword)
	if err != nil {
		uc.logger.Warn("Failed getting account",
			slog.Any("error", err))
		return nil, httperror.NewForBadRequestWithSingleField("message", fmt.Sprintf("failed getting wallet: %v", err))
	}

	return key, nil
}
