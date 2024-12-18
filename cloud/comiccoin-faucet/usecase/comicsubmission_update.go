package usecase

import (
	"context"
	"log/slog"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/common/httperror"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/config"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/domain"
)

type ComicSubmissionAcceptOperationUseCase struct {
	config *config.Configuration
	logger *slog.Logger
	repo   domain.ComicSubmissionRepository
}

func NewComicSubmissionAcceptOperationUseCase(
	config *config.Configuration,
	logger *slog.Logger,
	repo domain.ComicSubmissionRepository,
) *ComicSubmissionAcceptOperationUseCase {
	return &ComicSubmissionAcceptOperationUseCase{config, logger, repo}
}

func (uc *ComicSubmissionAcceptOperationUseCase) Execute(ctx context.Context, comicSubmission *domain.ComicSubmission) error {
	//
	// STEP 1: Validation.
	//

	e := make(map[string]string)
	if comicSubmission == nil {
		e["comic_submission"] = "Comic submission is required"
	} else {

	}
	if len(e) != 0 {
		uc.logger.Warn("Failed validating",
			slog.Any("error", e))
		return httperror.NewForBadRequest(&e)
	}

	//
	// STEP 2: Update the database.
	//

	return uc.repo.UpdateByID(ctx, comicSubmission)
}
