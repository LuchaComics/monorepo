package usecase

import (
	"context"
	"fmt"
	"log/slog"

	pkgkeystore "github.com/LuchaComics/monorepo/cloud/comiccoin-authority/common/blockchain/keystore"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/common/httperror"
	sstring "github.com/LuchaComics/monorepo/cloud/comiccoin-authority/common/security/securestring"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/domain"
	"github.com/ethereum/go-ethereum/common"
)

type AccountEncryptKeyUseCase struct {
	logger   *slog.Logger
	keystore pkgkeystore.KeystoreAdapter
	repo     domain.AccountRepository
}

func NewAccountEncryptKeyUseCase(
	logger *slog.Logger,
	keystore pkgkeystore.KeystoreAdapter,
	repo domain.AccountRepository,
) *AccountEncryptKeyUseCase {
	return &AccountEncryptKeyUseCase{logger, keystore, repo}
}

func (uc *AccountEncryptKeyUseCase) Execute(ctx context.Context, dataDir string, walletPassword *sstring.SecureString) (*common.Address, []byte, error) {
	//
	// STEP 1: Validation.
	//

	e := make(map[string]string)
	if dataDir == "" {
		e["data_dir"] = "missing value"
	}
	if walletPassword == nil {
		e["wallet_password"] = "missing value"
	}
	if len(e) != 0 {
		uc.logger.Warn("Failed reading account key",
			slog.Any("error", e))
		return nil, nil, httperror.NewForBadRequest(&e)
	}

	//
	// STEP 2: Create the encryted physical wallet on file.
	//

	walletAddress, walletKeystoreBytes, err := uc.keystore.CreateWallet(walletPassword)
	if err != nil {
		uc.logger.Error("failed creating new keystore",
			slog.Any("error", err))
		return nil, nil, fmt.Errorf("failed creating new keystore: %s", err)
	}

	return &walletAddress, walletKeystoreBytes, nil
}
