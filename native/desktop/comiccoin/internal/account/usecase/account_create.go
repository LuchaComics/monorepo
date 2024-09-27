package usecase

import (
	"log/slog"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/account/domain"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/pkg/httperror"
	"github.com/ethereum/go-ethereum/common"
)

type CreateAccountUseCase struct {
	config *config.Config
	logger *slog.Logger
	repo   domain.AccountRepository
}

func NewCreateAccountUseCase(config *config.Config, logger *slog.Logger, repo domain.AccountRepository) *CreateAccountUseCase {
	return &CreateAccountUseCase{config, logger, repo}
}

func (uc *CreateAccountUseCase) Execute(id string, walletAddress common.Address, walletFilepath string) error {
	//
	// STEP 1: Validation.
	//

	e := make(map[string]string)
	if id == "" {
		e["id"] = "missing value"
	}
	if walletAddress.String() == "" {
		e["wallet_address"] = "missing value"
	}
	if walletFilepath == "" {
		e["wallet_filepath"] = "missing value"
	}
	if len(e) != 0 {
		uc.logger.Warn("Failed creating new account",
			slog.Any("error", e))
		return httperror.NewForBadRequest(&e)
	}

	//
	// STEP 2: Create our strucutre.
	//

	account := &domain.Account{
		ID:             id,
		WalletAddress:  walletAddress,
		WalletFilepath: walletFilepath,
	}

	//
	// STEP 3: Insert into database.
	//

	return uc.repo.Upsert(account)
}
