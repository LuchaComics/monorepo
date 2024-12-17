package service

import (
	"log/slog"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/common/httperror"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/domain"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/usecase"
)

type ComicSubmissionCountByFilterService struct {
	logger                              *slog.Logger
	comicSubmissionCountByFilterUseCase *usecase.ComicSubmissionCountByFilterUseCase
}

func NewComicSubmissionCountByFilterService(
	logger *slog.Logger,
	uc1 *usecase.ComicSubmissionCountByFilterUseCase,
) *ComicSubmissionCountByFilterService {
	return &ComicSubmissionCountByFilterService{logger, uc1}
}

type ComicSubmissionCountByFilterServiceResponseIDO struct {
	Count uint64 `bson:"count" json:"count"`
}

func (s *ComicSubmissionCountByFilterService) Execute(sessCtx mongo.SessionContext, filter *domain.ComicSubmissionFilter) (*ComicSubmissionCountByFilterServiceResponseIDO, error) {
	//
	// STEP 1: Validation.
	//

	e := make(map[string]string)
	if filter == nil {
		e["filter"] = "Comic submission is required"
	}
	if len(e) != 0 {
		s.logger.Warn("Failed validating",
			slog.Any("error", e))
		return &ComicSubmissionCountByFilterServiceResponseIDO{Count: 0}, httperror.NewForBadRequest(&e)
	}

	//
	// STEP 2: Count in database.
	//

	// s.logger.Debug("Counting by filter",
	// 	slog.Any("filter", filter))

	// Lookup the user in our database, else return a `400 Bad Request` error.
	count, err := s.comicSubmissionCountByFilterUseCase.Execute(sessCtx, filter)
	if err != nil {
		s.logger.Error("database error", slog.Any("err", err))
		return nil, err
	}
	return &ComicSubmissionCountByFilterServiceResponseIDO{Count: count}, nil
}
