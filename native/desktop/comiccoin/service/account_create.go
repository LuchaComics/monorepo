package service

import (
	"fmt"
	"log/slog"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/domain"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/usecase"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/common/blockchain/signature"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/common/httperror"
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

	walletKey, err := s.walletDecryptKeyUseCase.Execute(walletFilepath, walletPassword)
	if err != nil {
		s.logger.Error("failed getting wallet key",
			slog.Any("data_dir", dataDir),
			slog.Any("error", err))
		return nil, fmt.Errorf("failed getting wallet key: %s", err)
	}

	val := "ComicCoin Blockchain"

	// Break the signature into the 3 parts: R, S, and V.
	v1, r1, s1, err := signature.Sign(val, walletKey.PrivateKey)
	if err != nil {
		return nil, err
	}
	// Recombine and get our address from the signature.
	addressFromSig, err := signature.FromAddress(val, v1, r1, s1)
	if err != nil {
		return nil, err
	}

	// Defensive Code: Do a check to ensure our signer to be working correctly.
	if walletAddress.Hex() != addressFromSig {
		s.logger.Error("address from signature does not match the wallet address",
			slog.Any("addressFromSig", addressFromSig),
			slog.Any("walletAddress", walletAddress.Hex()),
			slog.Any("data_dir", dataDir))
		return nil, fmt.Errorf("address from signature at %v does not match the wallet address of %v", addressFromSig, walletAddress.Hex())
	}

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
