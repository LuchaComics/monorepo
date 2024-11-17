package usecase

import (
	"context"
	"log/slog"
	"math/big"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/common/httperror"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/domain"
	"github.com/ethereum/go-ethereum/common"
)

type UpsertAccountUseCase struct {
	logger *slog.Logger
	repo   domain.AccountRepository
}

func NewUpsertAccountUseCase(logger *slog.Logger, repo domain.AccountRepository) *UpsertAccountUseCase {
	return &UpsertAccountUseCase{logger, repo}
}

func (uc *UpsertAccountUseCase) Execute(ctx context.Context, address *common.Address, balance uint64, nonce *big.Int) error {
	//
	// STEP 1: Validation.
	//

	e := make(map[string]string)
	if address == nil {
		e["address"] = "missing value"
	}
	if len(e) != 0 {
		uc.logger.Warn("Validation failed for upsert",
			slog.Any("error", e))
		return httperror.NewForBadRequest(&e)
	}

	//
	// STEP 2: Upsert our strucutre.
	//

	account := &domain.Account{
		Address:    address,
		NonceBytes: nonce.Bytes(),
		Balance:    balance,
	}

	//
	// STEP 3: Insert into database.
	//

	return uc.repo.Upsert(ctx, account)
}
