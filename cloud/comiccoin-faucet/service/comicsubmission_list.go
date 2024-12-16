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
	comicSubmissionListByFilterUseCase *usecase.ComicSubmissionListByFilterUseCase
}

func NewComicSubmissionListByFilterService(
	logger *slog.Logger,
	uc1 *usecase.ComicSubmissionListByFilterUseCase,
) *ComicSubmissionListByFilterService {
	return &ComicSubmissionListByFilterService{logger, uc1}
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
	detail, err := s.comicSubmissionListByFilterUseCase.Execute(sessCtx, filter2)
	if err != nil {
		s.logger.Error("database error",
			slog.Any("err", err))
		return nil, err
	}

	// s.logger.Debug("fetched",
	// 	slog.Any("id", id),
	// 	slog.Any("detail", detail))

	return (*ComicSubmissionFilterResultResponseIDO)(detail), nil
}
