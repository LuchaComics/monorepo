package usecase

import (
	"context"
	"log/slog"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/common/httperror"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/config"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/domain"
)

type UserCountByFilterUseCase struct {
	config *config.Configuration
	logger *slog.Logger
	repo   domain.UserRepository
}

func NewUserCountByFilterUseCase(
	config *config.Configuration,
	logger *slog.Logger,
	repo domain.UserRepository,
) *UserCountByFilterUseCase {
	return &UserCountByFilterUseCase{config, logger, repo}
}

func (uc *UserCountByFilterUseCase) Execute(ctx context.Context, filter *domain.UserFilter) (uint64, error) {
	//
	// STEP 1: Validation.
	//

	e := make(map[string]string)
	if filter == nil {
		e["filter"] = "Comic submission is required"
	}
	if len(e) != 0 {
		uc.logger.Warn("Failed validating",
			slog.Any("error", e))
		return 0, httperror.NewForBadRequest(&e)
	}

	//
	// STEP 2: Count in database.
	//

	return uc.repo.CountByFilter(ctx, filter)
}
