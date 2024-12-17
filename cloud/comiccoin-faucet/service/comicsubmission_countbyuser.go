package service

import (
	"log/slog"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/common/httperror"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/usecase"
)

type ComicSubmissionCountByUserService struct {
	logger                            *slog.Logger
	comicSubmissionCountByUserUseCase *usecase.ComicSubmissionCountByUserUseCase
}

func NewComicSubmissionCountByUserService(
	logger *slog.Logger,
	uc1 *usecase.ComicSubmissionCountByUserUseCase,
) *ComicSubmissionCountByUserService {
	return &ComicSubmissionCountByUserService{logger, uc1}
}

type ComicSubmissionCountByUserServiceResponseIDO struct {
	Count uint64 `bson:"count" json:"count"`
}

func (s *ComicSubmissionCountByUserService) Execute(sessCtx mongo.SessionContext, userID primitive.ObjectID) (*ComicSubmissionCountByUserServiceResponseIDO, error) {
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
	count, err := s.comicSubmissionCountByUserUseCase.Execute(sessCtx, userID)
	if err != nil {
		s.logger.Error("database error", slog.Any("err", err))
		return nil, err
	}
	return &ComicSubmissionCountByUserServiceResponseIDO{Count: count}, nil
}
