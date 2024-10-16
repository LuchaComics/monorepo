package service

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/domain"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/usecase"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/pkg/httperror"
	"github.com/ethereum/go-ethereum/common"
)

type GetAccountService struct {
	config               *config.Config
	logger               *slog.Logger
	getAccountUseCase    *usecase.GetAccountUseCase
	getWalletUseCase     *usecase.GetWalletUseCase
	createAccountUseCase *usecase.CreateAccountUseCase
}

func NewGetAccountService(
	cfg *config.Config,
	logger *slog.Logger,
	uc1 *usecase.GetAccountUseCase,
	uc2 *usecase.GetWalletUseCase,
	uc3 *usecase.CreateAccountUseCase,
) *GetAccountService {
	return &GetAccountService{cfg, logger, uc1, uc2, uc3}
}

func (s *GetAccountService) Execute(address *common.Address) (*domain.Account, error) {
	//
	// STEP 1: Validation.
	//

	e := make(map[string]string)
	if address == nil {
		e["address"] = "missing value"
	}
	if len(e) != 0 {
		s.logger.Warn("Validation failed for getting account",
			slog.Any("error", e))
		return nil, httperror.NewForBadRequest(&e)
	}

	//
	// STEP 2: Get our account from our in-memory database if it exists.
	//

	account, err := s.getAccountUseCase.Execute(address)
	if err != nil {
		if !strings.Contains(err.Error(), "does not exist") {
			s.logger.Error("failed getting account",
				slog.Any("address", address),
				slog.Any("error", err))
			return nil, err
		}
	}
	if account != nil {
		return account, nil
	}

	//
	// STEP 3: (Optional) If the account doesn't exist, check the wallets db
	// and see if we have an empty wallet hanging around for this account.
	//

	wallet, err := s.getWalletUseCase.Execute(address)
	if err != nil {
		s.logger.Error("failed getting wallet",
			slog.Any("address", address),
			slog.Any("error", err))
		return nil, err
	}
	if wallet == nil {
		return nil, nil
	}

	//
	// STEP 4: (Optional) Populate our in-memory database.
	//

	if err := s.createAccountUseCase.Execute(address); err != nil {
		s.logger.Error("failed saving to accounts database",
			slog.Any("address", address),
			slog.Any("error", err))
		return nil, fmt.Errorf("failed saving to database: %s", err)
	}

	//
	// STEP 5: (Optional) Fetch again.
	//

	account, err = s.getAccountUseCase.Execute(address)
	if err != nil {
		s.logger.Error("failed getting account",
			slog.Any("address", address),
			slog.Any("error", err))
		return nil, err
	}
	return account, nil
}
