package usecase

import (
	"fmt"
	"log/slog"

	"github.com/ethereum/go-ethereum/common"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/domain"
	pkgkeystore "github.com/LuchaComics/monorepo/native/desktop/comiccoin/common/blockchain/keystore"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/common/httperror"
)

type WalletEncryptKeyUseCase struct {
	config *config.Config
	logger *slog.Logger
	repo   domain.WalletRepository
}

func NewWalletEncryptKeyUseCase(config *config.Config, logger *slog.Logger, repo domain.WalletRepository) *WalletEncryptKeyUseCase {
	return &WalletEncryptKeyUseCase{config, logger, repo}
}

func (uc *WalletEncryptKeyUseCase) Execute(dataDir string, password string) (*common.Address, string, error) {
	//
	// STEP 1: Validation.
	//

	e := make(map[string]string)
	if dataDir == "" {
		e["data_dir"] = "missing value"
	}
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

	walletAddress, walletFilepath, err := pkgkeystore.NewKeystore(dataDir, password)
	if err != nil {
		uc.logger.Error("failed creating new keystore",
			slog.Any("error", err))
		return nil, "", fmt.Errorf("failed creating new keystore: %s", err)
	}

	return &walletAddress, walletFilepath, nil
}
