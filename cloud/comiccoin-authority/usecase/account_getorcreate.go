package usecase

import (
	"context"
	"log/slog"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/common/httperror"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/config"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/domain"
	"github.com/ethereum/go-ethereum/common"
)

type GetOrCreateAccountUseCase struct {
	config *config.Configuration
	logger *slog.Logger
	repo   domain.AccountRepository
}

func NewGetOrCreateAccountUseCase(config *config.Configuration, logger *slog.Logger, repo domain.AccountRepository) *GetOrCreateAccountUseCase {
	return &GetOrCreateAccountUseCase{config, logger, repo}
}

func (uc *GetOrCreateAccountUseCase) Execute(ctx context.Context, walletAddress *common.Address, balance, nonce uint64) error {
	//
	// STEP 1: Validation.
	//

	e := make(map[string]string)
	if walletAddress == nil {
		e["address"] = "missing value"
	}
	if len(e) != 0 {
		uc.logger.Warn("Validation failed for upsert",
			slog.Any("error", e))
		return httperror.NewForBadRequest(&e)
	}

	//
	// STEP 2: Attempt to get from database.
	//

	// Skip error handling
	getAcc, _ := uc.repo.GetByAddress(ctx, walletAddress)
	if getAcc != nil {
		return nil
	}

	//
	// STEP 2: Create our record and save to database.
	//

	account := &domain.Account{
		Address: walletAddress,
		Nonce:   nonce,
		Balance: balance,
	}

	return uc.repo.Upsert(ctx, account)
}