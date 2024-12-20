package usecase

import (
	"context"
	"log/slog"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/common/httperror"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/config"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserTransactionDeleteUseCase struct {
	config *config.Configuration
	logger *slog.Logger
	repo   domain.UserTransactionRepository
}

func NewUserTransactionDeleteUseCase(
	config *config.Configuration,
	logger *slog.Logger,
	repo domain.UserTransactionRepository,
) *UserTransactionDeleteUseCase {
	return &UserTransactionDeleteUseCase{config, logger, repo}
}

func (uc *UserTransactionDeleteUseCase) Execute(ctx context.Context, id primitive.ObjectID) error {
	//
	// STEP 1: Validation.
	//

	e := make(map[string]string)
	if id.IsZero() {
		e["id"] = "UserTransaction ID is required"
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

	return uc.repo.DeleteByID(ctx, id)
}
