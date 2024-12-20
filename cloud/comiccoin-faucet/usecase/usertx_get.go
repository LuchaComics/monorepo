package usecase

import (
	"context"
	"log/slog"
	"math/big"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/common/httperror"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/config"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserTransactionGetUseCase struct {
	config *config.Configuration
	logger *slog.Logger
	repo   domain.UserTransactionRepository
}

func NewUserTransactionGetUseCase(
	config *config.Configuration,
	logger *slog.Logger,
	repo domain.UserTransactionRepository,
) *UserTransactionGetUseCase {
	return &UserTransactionGetUseCase{config, logger, repo}
}

func (uc *UserTransactionGetUseCase) ExecuteForID(ctx context.Context, id primitive.ObjectID) (*domain.UserTransaction, error) {
	//
	// STEP 1: Validation.
	//

	e := make(map[string]string)
	if id.IsZero() {
		e["id"] = "ID is required"
	}
	if len(e) != 0 {
		uc.logger.Warn("Failed validating",
			slog.Any("error", e))
		return nil, httperror.NewForBadRequest(&e)
	}

	//
	// STEP 2: Insert into database.
	//

	return uc.repo.GetByID(ctx, id)
}

func (uc *UserTransactionGetUseCase) ExecuteForNonce(ctx context.Context, nonce *big.Int) (*domain.UserTransaction, error) {
	//
	// STEP 1: Validation.
	//

	e := make(map[string]string)

	// Special thanks to "Is there another way of testing if a big.Int is 0?" via https://stackoverflow.com/a/64257532
	if len(nonce.Bits()) == 0 {
		e["nonce"] = "Nonce is required"
	}
	if len(e) != 0 {
		uc.logger.Warn("Failed validating",
			slog.Any("error", e))
		return nil, httperror.NewForBadRequest(&e)
	}

	//
	// STEP 2: Get into database.
	//

	return uc.repo.GetByNonce(ctx, nonce)
}
