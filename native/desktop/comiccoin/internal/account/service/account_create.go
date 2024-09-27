package service

import (
	"fmt"
	"log/slog"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/account/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/account/domain"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/account/usecase"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/pkg/blockchain/keystore"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/pkg/httperror"
)

type CreateAccountService struct {
	config               *config.Config
	logger               *slog.Logger
	createAccountUseCase *usecase.CreateAccountUseCase
	getAccountUseCase    *usecase.GetAccountUseCase
}

func NewCreateAccountService(
	cfg *config.Config,
	logger *slog.Logger,
	uc1 *usecase.CreateAccountUseCase,
	uc2 *usecase.GetAccountUseCase,
) *CreateAccountService {
	return &CreateAccountService{cfg, logger, uc1, uc2}
}

func (s *CreateAccountService) Execute(dataDir, id, walletPassword string) (*domain.Account, error) {
	//
	// STEP 1: Validation.
	//

	e := make(map[string]string)
	if dataDir == "" {
		e["data_dir"] = "missing value"
	}
	if id == "" {
		e["id"] = "missing value"
	}
	if walletPassword == "" {
		e["wallet_password"] = "missing value"
	}
	if len(e) != 0 {
		s.logger.Warn("Failed creating new account",
			slog.Any("error", e))
		return nil, httperror.NewForBadRequest(&e)
	}

	//
	// STEP 2: Create the encryted physical wallet on file.
	//

	walletAddress, walletFilepath, err := keystore.NewKeystore(dataDir, walletPassword)
	if err != nil {
		s.logger.Error("failed creating new keystore",
			slog.Any("id", id),
			slog.Any("error", err))
		return nil, fmt.Errorf("failed creating new keystore: %s", err)
	}

	//
	// STEP 3: Save to our database.
	//

	if err := s.createAccountUseCase.Execute(id, walletAddress, walletFilepath); err != nil {
		s.logger.Error("failed saving to database",
			slog.Any("id", id),
			slog.Any("error", err))
		return nil, fmt.Errorf("failed saving to database: %s", err)
	}

	//
	// STEP 4: Return the account.
	//

	return s.getAccountUseCase.Execute(id)
}
