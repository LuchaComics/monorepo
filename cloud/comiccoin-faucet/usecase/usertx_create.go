package usecase

import (
	"context"
	"log/slog"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/common/httperror"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/config"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/domain"
)

type CreateUserTransactionUseCase struct {
	config *config.Configuration
	logger *slog.Logger
	repo   domain.UserTransactionRepository
}

func NewCreateUserTransactionUseCase(
	config *config.Configuration,
	logger *slog.Logger,
	repo domain.UserTransactionRepository,
) *CreateUserTransactionUseCase {
	return &CreateUserTransactionUseCase{config, logger, repo}
}

func (uc *CreateUserTransactionUseCase) Execute(ctx context.Context, userTransaction *domain.UserTransaction) error {
	//
	// STEP 1: Validation.
	//

	e := make(map[string]string)
	if userTransaction == nil {
		e["userTransaction"] = "User transaction is required"
	} else {

	}
	if len(e) != 0 {
		uc.logger.Warn("Failed validating",
			slog.Any("error", e))
		return httperror.NewForBadRequest(&e)
	}

	//
	// STEP 2: Insert into database.
	//

	return uc.repo.Create(ctx, userTransaction)
}
