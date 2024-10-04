package service

import (
	"fmt"
	"log/slog"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/usecase"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/pkg/httperror"
	"github.com/ethereum/go-ethereum/accounts/keystore"
)

type GetKeyService struct {
	config                  *config.Config
	logger                  *slog.Logger
	getAccountUseCase       *usecase.GetAccountUseCase
	walletDecryptKeyUseCase *usecase.WalletDecryptKeyUseCase
}

func NewGetKeyService(
	cfg *config.Config,
	logger *slog.Logger,
	uc1 *usecase.GetAccountUseCase,
	uc2 *usecase.WalletDecryptKeyUseCase,
) *GetKeyService {
	return &GetKeyService{cfg, logger, uc1, uc2}
}

func (s *GetKeyService) Execute(accountID string, walletPassword string) (*keystore.Key, error) {
	//
	// STEP 1: Validation.
	//

	e := make(map[string]string)
	if accountID == "" {
		e["account_id"] = "missing value"
	}
	if len(e) != 0 {
		s.logger.Warn("Failed creating new account",
			slog.Any("error", e))
		return nil, httperror.NewForBadRequest(&e)
	}

	//
	// STEP 2: Return the account.
	//

	account, err := s.getAccountUseCase.Execute(accountID)
	if err != nil {
		s.logger.Error("failed getting from database",
			slog.Any("account_id", accountID),
			slog.Any("error", err))
		return nil, fmt.Errorf("failed getting from database: %s", err)
	}
	if account == nil {
		return nil, fmt.Errorf("failed getting from database: %s", "d.n.e.")
	}

	key, err := s.walletDecryptKeyUseCase.Execute(account.WalletFilepath, walletPassword)
	if err != nil {
		s.logger.Error("failed getting key",
			slog.Any("account_id", accountID),
			slog.Any("error", err))
		return nil, fmt.Errorf("failed getting key: %s", err)
	}
	if key == nil {
		return nil, fmt.Errorf("failed getting key: %s", "d.n.e.")
	}
	return key, nil
}
