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

type WalletDecryptKeyUseCase struct {
	logger   *slog.Logger
	keystore pkgkeystore.KeystoreAdapter
	repo     domain.WalletRepository
}

func NewWalletDecryptKeyUseCase(
	logger *slog.Logger,
	keystore pkgkeystore.KeystoreAdapter,
	repo domain.WalletRepository,
) *WalletDecryptKeyUseCase {
	return &WalletDecryptKeyUseCase{logger, keystore, repo}
}

func (uc *WalletDecryptKeyUseCase) Execute(ctx context.Context, keystoreBytes []byte, password *sstring.SecureString) (*keystore.Key, error) {
	//
	// STEP 1: Validation.
	//

	e := make(map[string]string)
	if keystoreBytes == nil {
		e["keystore_bytes"] = "missing value"
	}
	if password == nil {
		e["password"] = "missing value"
	}
	if len(e) != 0 {
		uc.logger.Warn("Failed reading wallet key",
			slog.Any("error", e))
		return nil, httperror.NewForBadRequest(&e)
	}

	//
	// STEP 2: Decrypt key
	//

	key, err := uc.keystore.OpenWallet(keystoreBytes, password)
	if err != nil {
		uc.logger.Warn("Failed getting wallet key",
			slog.Any("error", err))
		return nil, httperror.NewForBadRequestWithSingleField("message", fmt.Sprintf("failed getting wallet key: %v", err))
	}

	return key, nil
}
