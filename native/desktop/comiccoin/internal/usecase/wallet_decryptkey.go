package usecase

import (
	"fmt"
	"log/slog"

	"github.com/ethereum/go-ethereum/accounts/keystore"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/domain"
	pkgkeystore "github.com/LuchaComics/monorepo/native/desktop/comiccoin/pkg/blockchain/keystore"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/pkg/httperror"
)

type WalletDecryptKeyUseCase struct {
	config *config.Config
	logger *slog.Logger
	repo   domain.WalletRepository
}

func NewWalletDecryptKeyUseCase(config *config.Config, logger *slog.Logger, repo domain.WalletRepository) *WalletDecryptKeyUseCase {
	return &WalletDecryptKeyUseCase{config, logger, repo}
}

func (uc *WalletDecryptKeyUseCase) Execute(filepath string, password string) (*keystore.Key, error) {
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

	key, err := pkgkeystore.GetKeyAfterDecryptingWalletAtFilepath(filepath, password)
	if err != nil {
		uc.logger.Warn("Failed getting wallet key",
			slog.Any("error", err))
		return nil, httperror.NewForBadRequestWithSingleField("message", fmt.Sprintf("failed getting wallet key: %v", err))
	}

	return key, nil
}
