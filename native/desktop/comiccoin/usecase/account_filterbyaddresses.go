package usecase

import (
	"context"
	"log/slog"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/common/httperror"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/domain"
	"github.com/ethereum/go-ethereum/common"
)

type AccountsFilterByAddressesUseCase struct {
	logger *slog.Logger
	repo   domain.AccountRepository
}

func NewAccountsFilterByAddressesUseCase(logger *slog.Logger, repo domain.AccountRepository) *AccountsFilterByAddressesUseCase {
	return &AccountsFilterByAddressesUseCase{logger, repo}
}

func (uc *AccountsFilterByAddressesUseCase) Execute(ctx context.Context, addresses []*common.Address) ([]*domain.Account, error) {
	//
	// STEP 1: Validation.
	//

	e := make(map[string]string)
	if addresses == nil {
		e["addresses"] = "missing value"
	}
	if len(addresses) == 0 {
		e["addresses"] = "empty array"
	}
	if len(e) != 0 {
		uc.logger.Warn("Validation failed",
			slog.Any("error", e))
		return nil, httperror.NewForBadRequest(&e)
	}

	//
	// STEP 2: Insert into database.
	//

	return uc.repo.ListWithFilterByAddresses(ctx, addresses)
}
