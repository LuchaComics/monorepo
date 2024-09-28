package service

import (
	"fmt"
	"log/slog"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/account/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/account/usecase"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/pkg/httperror"
	"github.com/ethereum/go-ethereum/accounts/keystore"
)

type GetKeyService struct {
	config                   *config.Config
	logger                   *slog.Logger
	getAccountUseCase        *usecase.GetAccountUseCase
	accountDecryptKeyUseCase *usecase.AccountDecryptKeyUseCase
}

func NewGetKeyService(
	cfg *config.Config,
	logger *slog.Logger,
	uc1 *usecase.GetAccountUseCase,
	uc2 *usecase.AccountDecryptKeyUseCase,
) *GetKeyService {
	return &GetKeyService{cfg, logger, uc1, uc2}
}

func (s *GetKeyService) Execute(id string, walletPassword string) (*keystore.Key, error) {
	//
	// STEP 1: Validation.
	//

	e := make(map[string]string)
	if id == "" {
		e["id"] = "missing value"
	}
	if len(e) != 0 {
		s.logger.Warn("Failed creating new account",
			slog.Any("error", e))
		return nil, httperror.NewForBadRequest(&e)
	}

	//
	// STEP 2: Return the account.
	//

	account, err := s.getAccountUseCase.Execute(id)
	if err != nil {
		s.logger.Error("failed getting from database",
			slog.Any("id", id),
			slog.Any("error", err))
		return nil, fmt.Errorf("failed getting from database: %s", err)
	}
	if account == nil {
		return nil, fmt.Errorf("failed getting from database: %s", "d.n.e.")
	}

	key, err := s.accountDecryptKeyUseCase.Execute(account.WalletFilepath, walletPassword)
	if err != nil {
		s.logger.Error("failed getting key",
			slog.Any("id", id),
			slog.Any("error", err))
		return nil, fmt.Errorf("failed getting key: %s", err)
	}
	if key == nil {
		return nil, fmt.Errorf("failed getting key: %s", "d.n.e.")
	}
	return key, nil
}
