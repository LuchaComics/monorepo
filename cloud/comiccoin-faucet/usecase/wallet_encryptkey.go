package usecase

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/ethereum/go-ethereum/common"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/common/blockchain/keystore"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/common/httperror"
	sstring "github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/common/security/securestring"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/config"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/domain"
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

func (uc *WalletEncryptKeyUseCase) Execute(ctx context.Context, password *sstring.SecureString) (*common.Address, []byte, error) {
	//
	// STEP 1: Validation.
	//

	e := make(map[string]string)
	if password == nil {
		e["password"] = "missing value"
	}
	if len(e) != 0 {
		uc.logger.Warn("Failed reading wallet key",
			slog.Any("error", e))
		return nil, nil, httperror.NewForBadRequest(&e)
	}

	//
	// STEP 2: Create the encryted physical wallet on file.
	//

	uc.logger.Debug("Creating wallet...")

	walletAddress, walletKeystoreBytes, err := uc.keystore.CreateWallet(password)
	if err != nil {
		uc.logger.Error("failed creating new keystore",
			slog.Any("error", err))
		return nil, nil, fmt.Errorf("failed creating new keystore: %s", err)
	}

	uc.logger.Debug("New wallet created", slog.Any("address", walletAddress))

	return &walletAddress, walletKeystoreBytes, nil
}
