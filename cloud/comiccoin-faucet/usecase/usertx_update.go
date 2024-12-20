package usecase

import (
	"context"
	"log/slog"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/common/httperror"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/config"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/domain"
)

type UserTransactionUpdateUseCase struct {
	config *config.Configuration
	logger *slog.Logger
	repo   domain.UserTransactionRepository
}

func NewUserTransactionUpdateUseCase(
	config *config.Configuration,
	logger *slog.Logger,
	repo domain.UserTransactionRepository,
) *UserTransactionUpdateUseCase {
	return &UserTransactionUpdateUseCase{config, logger, repo}
}

func (uc *UserTransactionUpdateUseCase) Execute(ctx context.Context, userTransaction *domain.UserTransaction) error {
	//
	// STEP 1: Validation.
	//

	e := make(map[string]string)
	if userTransaction == nil {
		e["user_transaction"] = "User transaction is required"
	} else {

	}
	if len(e) != 0 {
		uc.logger.Warn("Failed validating",
			slog.Any("error", e))
		return httperror.NewForBadRequest(&e)
	}

	//
	// STEP 2: Update in our database.
	//

	return uc.repo.UpdateByID(ctx, userTransaction)
}
