package usecase

import (
	"fmt"
	"log/slog"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/domain"
	pkgkeystore "github.com/LuchaComics/monorepo/native/desktop/comiccoin/common/blockchain/keystore"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/common/httperror"
	"github.com/ethereum/go-ethereum/common"
)

type AccountEncryptKeyUseCase struct {
	config *config.Config
	logger *slog.Logger
	repo   domain.AccountRepository
}

func NewAccountEncryptKeyUseCase(config *config.Config, logger *slog.Logger, repo domain.AccountRepository) *AccountEncryptKeyUseCase {
	return &AccountEncryptKeyUseCase{config, logger, repo}
}

func (uc *AccountEncryptKeyUseCase) Execute(dataDir string, walletPassword string) (*common.Address, string, error) {
	//
	// STEP 1: Validation.
	//

	e := make(map[string]string)
	if dataDir == "" {
		e["data_dir"] = "missing value"
	}
	if walletPassword == "" {
		e["wallet_password"] = "missing value"
	}
	if len(e) != 0 {
		uc.logger.Warn("Failed reading account key",
			slog.Any("error", e))
		return nil, "", httperror.NewForBadRequest(&e)
	}

	//
	// STEP 2: Create the encryted physical wallet on file.
	//

	walletAddress, walletFilepath, err := pkgkeystore.NewKeystore(dataDir, walletPassword)
	if err != nil {
		uc.logger.Error("failed creating new keystore",
			slog.Any("error", err))
		return nil, "", fmt.Errorf("failed creating new keystore: %s", err)
	}

	return &walletAddress, walletFilepath, nil
}
