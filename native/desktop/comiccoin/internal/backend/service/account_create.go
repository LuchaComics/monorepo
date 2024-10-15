package service

import (
	"fmt"
	"log/slog"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/backend/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/backend/domain"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/backend/usecase"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/pkg/httperror"
)

type CreateAccountService struct {
	config                  *config.Config
	logger                  *slog.Logger
	walletEncryptKeyUseCase *usecase.WalletEncryptKeyUseCase
	walletDecryptKeyUseCase *usecase.WalletDecryptKeyUseCase
	createWalletUseCase     *usecase.CreateWalletUseCase
	createAccountUseCase    *usecase.CreateAccountUseCase
	getAccountUseCase       *usecase.GetAccountUseCase
}

func NewCreateAccountService(
	cfg *config.Config,
	logger *slog.Logger,
	uc1 *usecase.WalletEncryptKeyUseCase,
	uc2 *usecase.WalletDecryptKeyUseCase,
	uc3 *usecase.CreateWalletUseCase,
	uc4 *usecase.CreateAccountUseCase,
	uc5 *usecase.GetAccountUseCase,
) *CreateAccountService {
	return &CreateAccountService{cfg, logger, uc1, uc2, uc3, uc4, uc5}
}

func (s *CreateAccountService) Execute(dataDir, walletPassword string) (*domain.Account, error) {
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
		s.logger.Warn("Failed creating new account",
			slog.Any("error", e))
		return nil, httperror.NewForBadRequest(&e)
	}

	//
	// STEP 2:
	// Create the encryted physical wallet on file.
	//

	walletAddress, walletFilepath, err := s.walletEncryptKeyUseCase.Execute(dataDir, walletPassword)
	if err != nil {
		s.logger.Error("failed creating new keystore",
			slog.Any("data_dir", dataDir),
			slog.Any("error", err))
		return nil, fmt.Errorf("failed creating new keystore: %s", err)
	}

	//
	// STEP 3:
	// Decrypt the wallet so we can extract data from it.
	//

	// walletKey, err := s.walletDecryptKeyUseCase.Execute(walletFilepath, walletPassword)
	// if err != nil {
	// 	s.logger.Error("failed getting wallet key",
	// 		slog.Any("data_dir", dataDir),
	// 		slog.Any("error", err))
	// 	return nil, fmt.Errorf("failed getting wallet key: %s", err)
	// }

	//
	// STEP 4:
	// Converts the wallet's public key to an account value.
	//

	// // DEVELOPERS NOTE:
	// // AccountID represents an account id that is used to sign transactions and is
	// // associated with transactions on the blockchain. This will be the last 20
	// // bytes of the public key.
	// privateKey := walletKey.PrivateKey
	// publicKey := privateKey.PublicKey
	// accountID := crypto.PubkeyToAddress(publicKey).String()

	//
	// STEP 3:
	// Save wallet to our database.
	//

	if err := s.createWalletUseCase.Execute(walletAddress, walletFilepath); err != nil {
		s.logger.Error("failed saving to database",
			slog.Any("data_dir", dataDir),
			slog.Any("error", err))
		return nil, fmt.Errorf("failed saving to database: %s", err)
	}

	//
	// STEP 4: Create the account.
	//

	if err := s.createAccountUseCase.Execute(walletAddress); err != nil {
		s.logger.Error("failed saving to database",
			slog.Any("data_dir", dataDir),
			slog.Any("error", err))
		return nil, fmt.Errorf("failed saving to database: %s", err)
	}

	//
	// STEP 5: Return the saved account.
	//

	return s.getAccountUseCase.Execute(walletAddress)
}
