package usecase

import (
	"fmt"
	"log/slog"

	"github.com/ethereum/go-ethereum/accounts/keystore"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/domain"
	pkgkeystore "github.com/LuchaComics/monorepo/native/desktop/comiccoin/pkg/blockchain/keystore"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/pkg/httperror"
)

type AccountDecryptKeyUseCase struct {
	config *config.Config
	logger *slog.Logger
	repo   domain.AccountRepository
}

func NewAccountDecryptKeyUseCase(config *config.Config, logger *slog.Logger, repo domain.AccountRepository) *AccountDecryptKeyUseCase {
	return &AccountDecryptKeyUseCase{config, logger, repo}
}

func (uc *AccountDecryptKeyUseCase) Execute(walletFilepath string, walletPassword string) (*keystore.Key, error) {
	//
	// STEP 1: Validation.
	//

	e := make(map[string]string)
	if walletFilepath == "" {
		e["wallet_filepath"] = "missing value"
	}
	if walletPassword == "" {
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

	key, err := pkgkeystore.GetKeyAfterDecryptingWalletAtFilepath(walletFilepath, walletPassword)
	if err != nil {
		uc.logger.Warn("Failed getting account",
			slog.Any("error", err))
		return nil, httperror.NewForBadRequestWithSingleField("message", fmt.Sprintf("failed getting wallet: %v", err))
	}

	return key, nil
}
