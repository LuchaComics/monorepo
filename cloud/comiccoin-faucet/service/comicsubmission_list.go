package service

import (
	"log/slog"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/common/httperror"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/domain"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/usecase"
)

type ComicSubmissionListByFilterService struct {
	logger                             *slog.Logger
	cloudStoragePresignedURLUseCase    *usecase.CloudStoragePresignedURLUseCase
	comicSubmissionListByFilterUseCase *usecase.ComicSubmissionListByFilterUseCase
}

func NewComicSubmissionListByFilterService(
	logger *slog.Logger,
	uc1 *usecase.CloudStoragePresignedURLUseCase,
	uc2 *usecase.ComicSubmissionListByFilterUseCase,
) *ComicSubmissionListByFilterService {
	return &ComicSubmissionListByFilterService{logger, uc1, uc2}
}

type ComicSubmissionFilterRequestID domain.ComicSubmissionFilter

type ComicSubmissionFilterResultResponseIDO domain.ComicSubmissionFilterResult

func (s *ComicSubmissionListByFilterService) Execute(sessCtx mongo.SessionContext, filter *ComicSubmissionFilterRequestID) (*ComicSubmissionFilterResultResponseIDO, error) {
	//
	// STEP 1: Validation.
	//

	e := make(map[string]string)
	if filter == nil {
		e["filter"] = "UserID is required"
	}
	if len(e) != 0 {
		s.logger.Warn("Failed validating",
			slog.Any("error", e))
		return nil, httperror.NewForBadRequest(&e)
	}

	//
	// STEP 2: Count in database.
	//

	filter2 := (*domain.ComicSubmissionFilter)(filter)

	// Lookup the user in our database, else return a `400 Bad Request` error.
	listResp, err := s.comicSubmissionListByFilterUseCase.Execute(sessCtx, filter2)
	if err != nil {
		s.logger.Error("database error",
			slog.Any("err", err))
		return nil, err
	}

	for _, submission := range listResp.Submissions {
		// Front cover.
		frontPresignedURL, err := s.cloudStoragePresignedURLUseCase.Execute(sessCtx, submission.FrontCover.ObjectKey)
		if err != nil {
			s.logger.Error("Failed generating presigned url via cloud storage", slog.Any("err", err))
			return nil, err
		}
		submission.FrontCover.ObjectURL = frontPresignedURL

		// Back cover.
		backPresignedURL, err := s.cloudStoragePresignedURLUseCase.Execute(sessCtx, submission.BackCover.ObjectKey)
		if err != nil {
			s.logger.Error("Failed generating presigned url via cloud storage", slog.Any("err", err))
			return nil, err
		}
		submission.BackCover.ObjectURL = backPresignedURL
	}

	// s.logger.Debug("fetched",
	// 	slog.Any("id", id),
	// 	slog.Any("detail", detail))

	return (*ComicSubmissionFilterResultResponseIDO)(listResp), nil
}
