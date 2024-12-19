package service

import (
	"log/slog"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/usecase"
)

type ComicSubmissionTotalCoinsAwardedService struct {
	logger                                  *slog.Logger
	comicSubmissionTotalCoinsAwardedUseCase *usecase.ComicSubmissionTotalCoinsAwardedUseCase
}

func NewComicSubmissionTotalCoinsAwardedService(
	logger *slog.Logger,
	uc1 *usecase.ComicSubmissionTotalCoinsAwardedUseCase,
) *ComicSubmissionTotalCoinsAwardedService {
	return &ComicSubmissionTotalCoinsAwardedService{logger, uc1}
}

type ComicSubmissionTotalCoinsAwardedResponseIDO struct {
	Count uint64 `bson:"count" json:"count"`
}

func (s *ComicSubmissionTotalCoinsAwardedService) Execute(sessCtx mongo.SessionContext) (*ComicSubmissionTotalCoinsAwardedResponseIDO, error) {
	count, err := s.comicSubmissionTotalCoinsAwardedUseCase.Execute(sessCtx)
	if err != nil {
		s.logger.Error("database error",
			slog.Any("err", err))
		return nil, err
	}
	return &ComicSubmissionTotalCoinsAwardedResponseIDO{Count: count}, nil
}
