package usecase

import (
	"context"
	"log/slog"
	"math/big"

	"github.com/ethereum/go-ethereum/common"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/common/httperror"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/domain"
)

type CreateAccountUseCase struct {
	logger *slog.Logger
	repo   domain.AccountRepository
}

func NewCreateAccountUseCase(logger *slog.Logger, repo domain.AccountRepository) *CreateAccountUseCase {
	return &CreateAccountUseCase{logger, repo}
}

func (uc *CreateAccountUseCase) Execute(ctx context.Context, address *common.Address) error {
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
		Address:    address,
		NonceBytes: big.NewInt(0).Bytes(),
		Balance:    0,
	}

	//
	// STEP 3: Insert into database.
	//

	return uc.repo.Upsert(ctx, account)
}
