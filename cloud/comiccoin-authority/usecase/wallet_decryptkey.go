package usecase

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/ethereum/go-ethereum/accounts/keystore"

	pkgkeystore "github.com/LuchaComics/monorepo/cloud/comiccoin-authority/common/blockchain/keystore"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/common/httperror"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/config"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/domain"
)

type WalletDecryptKeyUseCase struct {
	config   *config.Configuration
	logger   *slog.Logger
	keystore pkgkeystore.KeystoreAdapter
	repo     domain.WalletRepository
}

func NewWalletDecryptKeyUseCase(
	config *config.Configuration,
	logger *slog.Logger,
	keystore pkgkeystore.KeystoreAdapter,
	repo domain.WalletRepository,
) *WalletDecryptKeyUseCase {
	return &WalletDecryptKeyUseCase{config, logger, keystore, repo}
}

func (uc *WalletDecryptKeyUseCase) Execute(ctx context.Context, filepath string, password string) (*keystore.Key, error) {
	//
	// STEP 1: Validation.
	//

	e := make(map[string]string)
	if filepath == "" {
		e["filepath"] = "missing value"
	}
	if password == "" {
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

	key, err := uc.keystore.OpenWallet(filepath, password)
	if err != nil {
		uc.logger.Warn("Failed getting wallet key",
			slog.Any("error", err))
		return nil, httperror.NewForBadRequestWithSingleField("message", fmt.Sprintf("failed getting wallet key: %v", err))
	}

	return key, nil
}