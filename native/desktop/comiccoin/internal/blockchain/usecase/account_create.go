package usecase

import (
	"log/slog"

	"github.com/ethereum/go-ethereum/common"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/domain"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/pkg/httperror"
)

type CreateAccountUseCase struct {
	config *config.Config
	logger *slog.Logger
	repo   domain.AccountRepository
}

func NewCreateAccountUseCase(config *config.Config, logger *slog.Logger, repo domain.AccountRepository) *CreateAccountUseCase {
	return &CreateAccountUseCase{config, logger, repo}
}

func (uc *CreateAccountUseCase) Execute(address *common.Address) error {
	//
	// STEP 1: Validation.
	//

	e := make(map[string]string)
	if address == nil {
		e["address"] = "missing value"
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
		Address: address,
		Nonce:   0,
		Balance: 0,
	}

	//
	// STEP 3: Insert into database.
	//

	return uc.repo.Upsert(account)
}
