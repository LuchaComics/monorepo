package service

import (
	"log/slog"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/common/httperror"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/usecase"
)

type ComicSubmissionJudgeService struct {
	logger                        *slog.Logger
	comicSubmissionGetByIDUseCase *usecase.ComicSubmissionGetByIDUseCase
}

func NewComicSubmissionJudgeService(
	logger *slog.Logger,
	uc1 *usecase.ComicSubmissionGetByIDUseCase,
) *ComicSubmissionJudgeService {
	return &ComicSubmissionJudgeService{logger, uc1}
}

func (s *ComicSubmissionJudgeService) Execute(sessCtx mongo.SessionContext, comicSubmissionID primitive.ObjectID) (*ComicSubmissionResponseIDO, error) {
	//
	// STEP 1: Validation.
	//

	e := make(map[string]string)
	if comicSubmissionID.IsZero() {
		e["comic_submission_id"] = "Comic submission identifier is required"
	}
	if len(e) != 0 {
		s.logger.Warn("Failed validating",
			slog.Any("error", e))
		return nil, httperror.NewForBadRequest(&e)
	}

	//
	// STEP 2: Count in database.
	//

	// Lookup the user in our database, else return a `400 Bad Request` error.
	comicSubmission, err := s.comicSubmissionGetByIDUseCase.Execute(sessCtx, comicSubmissionID)
	if err != nil {
		s.logger.Error("database error",
			slog.Any("err", err))
		return nil, err
	}

	// s.logger.Debug("fetched",
	// 	slog.Any("id", id),
	// 	slog.Any("detail", detail))

	return (*ComicSubmissionResponseIDO)(comicSubmission), nil
}
