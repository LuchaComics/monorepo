package usecase

import (
	"context"
	"log/slog"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/common/httperror"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/config"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/domain"
	"github.com/ethereum/go-ethereum/common"
)

type GetAccountUseCase struct {
	config *config.Configuration
	logger *slog.Logger
	repo   domain.AccountRepository
}

func NewGetAccountUseCase(config *config.Configuration, logger *slog.Logger, repo domain.AccountRepository) *GetAccountUseCase {
	return &GetAccountUseCase{config, logger, repo}
}

func (uc *GetAccountUseCase) Execute(ctx context.Context, address *common.Address) (*domain.Account, error) {
	//
	// STEP 1: Validation.
	//

	e := make(map[string]string)
	if address == nil {
		e["address"] = "missing value"
	}
	if len(e) != 0 {
		uc.logger.Warn("Failed getting account",
			slog.Any("error", e))
		return nil, httperror.NewForBadRequest(&e)
	}

	//
	// STEP 2: Insert into database.
	//

	// uc.logger.Debug("Getting account...",
	// 	slog.Any("address", address))

	return uc.repo.GetByAddress(ctx, address)
}
