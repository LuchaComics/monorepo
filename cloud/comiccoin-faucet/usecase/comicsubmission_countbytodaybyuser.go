package usecase

import (
	"context"
	"log/slog"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/common/httperror"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/config"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ComicSubmissionCountByStatusAndByUserUseCase struct {
	config *config.Configuration
	logger *slog.Logger
	repo   domain.ComicSubmissionRepository
}

func NewComicSubmissionCountByStatusAndByUserUseCase(
	config *config.Configuration,
	logger *slog.Logger,
	repo domain.ComicSubmissionRepository,
) *ComicSubmissionCountByStatusAndByUserUseCase {
	return &ComicSubmissionCountByStatusAndByUserUseCase{config, logger, repo}
}

func (uc *ComicSubmissionCountByStatusAndByUserUseCase) Execute(ctx context.Context, status int8, userID primitive.ObjectID) (uint64, error) {
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

	return uc.repo.CountByStatusAndByUserID(ctx, status, userID)
}
