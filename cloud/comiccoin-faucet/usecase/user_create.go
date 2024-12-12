package usecase

import (
	"context"
	"log/slog"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/common/httperror"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/config"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/domain"
)

type UserCreateUseCase struct {
	config *config.Configuration
	logger *slog.Logger
	repo   domain.UserRepository
}

func NewUserCreateUseCase(config *config.Configuration, logger *slog.Logger, repo domain.UserRepository) *UserCreateUseCase {
	return &UserCreateUseCase{config, logger, repo}
}

func (uc *UserCreateUseCase) Execute(ctx context.Context, user *domain.User) error {
	//
	// STEP 1: Validation.
	//

	e := make(map[string]string)
	if user == nil {
		e["user"] = "missing value"
	} else {
		//TODO: IMPL.
	}
	if len(e) != 0 {
		uc.logger.Warn("Validation failed for upsert",
			slog.Any("error", e))
		return httperror.NewForBadRequest(&e)
	}
	//
	// STEP 2: Insert into database.
	//

	return uc.repo.Create(ctx, user)
}
