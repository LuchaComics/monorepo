package service

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/common/blockchain/signature"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/common/httperror"
	sstring "github.com/LuchaComics/monorepo/cloud/comiccoin-authority/common/security/securestring"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/domain"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/usecase"
)

type CreateAccountService struct {
	logger                  *slog.Logger
	walletEncryptKeyUseCase *usecase.WalletEncryptKeyUseCase
	walletDecryptKeyUseCase *usecase.WalletDecryptKeyUseCase
	createWalletUseCase     *usecase.CreateWalletUseCase
	createAccountUseCase    *usecase.CreateAccountUseCase
	getAccountUseCase       *usecase.GetAccountUseCase
}

func NewCreateAccountService(
	logger *slog.Logger,
	uc1 *usecase.WalletEncryptKeyUseCase,
	uc2 *usecase.WalletDecryptKeyUseCase,
	uc3 *usecase.CreateWalletUseCase,
	uc4 *usecase.CreateAccountUseCase,
	uc5 *usecase.GetAccountUseCase,
) *CreateAccountService {
	return &CreateAccountService{logger, uc1, uc2, uc3, uc4, uc5}
}

func (s *CreateAccountService) Execute(ctx context.Context, walletPassword *sstring.SecureString, walletPasswordRepeated *sstring.SecureString, walletLabel string) (*domain.Account, error) {
	//
	// STEP 1: Validation.
	//

	e := make(map[string]string)
	if walletPassword == nil {
		e["wallet_password"] = "missing value"
	}
	if walletPasswordRepeated == nil {
		e["wallet_password_repeated"] = "missing value"
	}
	if walletPassword != nil && walletPasswordRepeated != nil {
		if walletPassword.String() != walletPasswordRepeated.String() {
			s.logger.Error("passwords do not match",
				slog.String("wallet_label", walletLabel),
			)
			e["wallet_password"] = "do not match"
			e["wallet_password_repeated"] = "do not match"
		}
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

	walletAddress, walletFilepath, err := s.walletEncryptKeyUseCase.Execute(ctx, walletPassword)
	if err != nil {
		s.logger.Error("failed creating new keystore",
			slog.Any("error", err))
		return nil, fmt.Errorf("failed creating new keystore: %s", err)
	}

	s.logger.Debug("Created new wallet for account",
		slog.Any("wallet_address", walletAddress),
		slog.Any("wallet_filepath", walletFilepath))

	//
	// STEP 3:
	// Decrypt the wallet so we can extract data from it.
	//

	walletKey, err := s.walletDecryptKeyUseCase.Execute(ctx, walletFilepath, walletPassword)
	if err != nil {
		s.logger.Error("failed getting wallet key",
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
			slog.Any("walletAddress", walletAddress.Hex()))
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

	if err := s.createWalletUseCase.Execute(ctx, walletAddress, walletFilepath, walletLabel); err != nil {
		s.logger.Error("failed saving to database",
			slog.Any("error", err))
		return nil, fmt.Errorf("failed saving to database: %s", err)
	}

	//
	// STEP 4: Create the account.
	//

	if err := s.createAccountUseCase.Execute(ctx, walletAddress); err != nil {
		s.logger.Error("failed saving to database",
			slog.Any("error", err))
		return nil, fmt.Errorf("failed saving to database: %s", err)
	}

	//
	// STEP 5: Return the saved account.
	//

	return s.getAccountUseCase.Execute(ctx, walletAddress)
}
