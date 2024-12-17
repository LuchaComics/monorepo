package usecase

import (
	"context"
	"log/slog"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/common/httperror"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/config"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ComicSubmissionCountByUserUseCase struct {
	config *config.Configuration
	logger *slog.Logger
	repo   domain.ComicSubmissionRepository
}

func NewComicSubmissionCountByUserUseCase(
	config *config.Configuration,
	logger *slog.Logger,
	repo domain.ComicSubmissionRepository,
) *ComicSubmissionCountByUserUseCase {
	return &ComicSubmissionCountByUserUseCase{config, logger, repo}
}

func (uc *ComicSubmissionCountByUserUseCase) Execute(ctx context.Context, userID primitive.ObjectID) (uint64, error) {
	//
	// STEP 1: Validation.
	//

	e := make(map[string]string)
	if userID.IsZero() {
		e["userID"] = "Comic submission is required"
	}
	if len(e) != 0 {
		uc.logger.Warn("Failed validating",
			slog.Any("error", e))
		return 0, httperror.NewForBadRequest(&e)
	}

	//
	// STEP 2: Count in database.
	//

	return uc.repo.CountByUserID(ctx, userID)
}
