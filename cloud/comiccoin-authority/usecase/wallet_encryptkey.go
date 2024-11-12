package usecase

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/ethereum/go-ethereum/common"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/common/blockchain/keystore"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/common/httperror"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/config"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/domain"
)

type WalletEncryptKeyUseCase struct {
	config   *config.Configuration
	logger   *slog.Logger
	keystore keystore.KeystoreAdapter
	repo     domain.WalletRepository
}

func NewWalletEncryptKeyUseCase(
	config *config.Configuration,
	logger *slog.Logger,
	keystore keystore.KeystoreAdapter,
	repo domain.WalletRepository,
) *WalletEncryptKeyUseCase {
	return &WalletEncryptKeyUseCase{config, logger, keystore, repo}
}

func (uc *WalletEncryptKeyUseCase) Execute(ctx context.Context, password string) (*common.Address, string, error) {
	//
	// STEP 1: Validation.
	//

	e := make(map[string]string)
	if password == "" {
		e["password"] = "missing value"
	}
	if len(e) != 0 {
		uc.logger.Warn("Failed reading wallet key",
			slog.Any("error", e))
		return nil, "", httperror.NewForBadRequest(&e)
	}

	//
	// STEP 2: Create the encryted physical wallet on file.
	//

	walletAddress, walletFilepath, err := uc.keystore.CreateWallet(password)
	if err != nil {
		uc.logger.Error("failed creating new keystore",
			slog.Any("error", err))
		return nil, "", fmt.Errorf("failed creating new keystore: %s", err)
	}

	return &walletAddress, walletFilepath, nil
}