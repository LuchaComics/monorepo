package service

import (
	"log/slog"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/common/httperror"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/usecase"
)

type ComicSubmissionCountTotalCreatedTodayByUserService struct {
	logger                                             *slog.Logger
	comicSubmissionCountTotalCreatedTodayByUserUseCase *usecase.ComicSubmissionCountTotalCreatedTodayByUserUseCase
}

func NewComicSubmissionCountTotalCreatedTodayByUserService(
	logger *slog.Logger,
	uc1 *usecase.ComicSubmissionCountTotalCreatedTodayByUserUseCase,
) *ComicSubmissionCountTotalCreatedTodayByUserService {
	return &ComicSubmissionCountTotalCreatedTodayByUserService{logger, uc1}
}

type ComicSubmissionCountTotalCreatedTodayByUserServiceResponseIDO struct {
	Count uint64 `bson:"count" json:"count"`
}

func (s *ComicSubmissionCountTotalCreatedTodayByUserService) Execute(sessCtx mongo.SessionContext, userID primitive.ObjectID) (*ComicSubmissionCountTotalCreatedTodayByUserServiceResponseIDO, error) {
	//
	// STEP 1: Validation.
	//

	e := make(map[string]string)
	if userID.IsZero() {
		e["user_id"] = "UserID is required"
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
	count, err := s.comicSubmissionCountTotalCreatedTodayByUserUseCase.Execute(sessCtx, userID)
	if err != nil {
		s.logger.Error("database error", slog.Any("err", err))
		return nil, err
	}
	return &ComicSubmissionCountTotalCreatedTodayByUserServiceResponseIDO{Count: count}, nil
}
