package service

import (
	"log/slog"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/domain"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/usecase"
	"github.com/bartmika/timekit"
)

type UserCountJoinedThisWeekService struct {
	logger                              *slog.Logger
	comicSubmissionCountByFilterUseCase *usecase.UserCountByFilterUseCase
}

func NewUserCountJoinedThisWeekService(
	logger *slog.Logger,
	uc1 *usecase.UserCountByFilterUseCase,
) *UserCountJoinedThisWeekService {
	return &UserCountJoinedThisWeekService{logger, uc1}
}

type UserCountJoinedThisWeekResponseIDO struct {
	Count uint64 `bson:"count" json:"count"`
}

func (s *UserCountJoinedThisWeekService) Execute(sessCtx mongo.SessionContext, tenantID primitive.ObjectID, timezone string) (*UserCountJoinedThisWeekResponseIDO, error) {
	//
	// STEP 1: Generate range.
	//

	loc, err := time.LoadLocation(timezone)
	if err != nil {
		s.logger.Warn("Failed validating",
			slog.Any("error", err))
		return nil, err
	}
	now := time.Now()
	userTime := now.In(loc)

	thisWeekStart := timekit.FirstDayOfThisISOWeek(func() time.Time {
		return userTime
	})
	thisWeekEnd := timekit.LastDayOfThisISOWeek(func() time.Time {
		return userTime
	})

	//
	// STEP 2: Count in database.
	//

	filter := &domain.UserFilter{
		TenantID:       tenantID,
		CreatedAtStart: &thisWeekStart,
		CreatedAtEnd:   &thisWeekEnd,
	}

	// s.logger.Debug("Counting by filter",
	// 	slog.Any("filter", filter))

	// Lookup the user in our database, else return a `400 Bad Request` error.
	count, err := s.comicSubmissionCountByFilterUseCase.Execute(sessCtx, filter)
	if err != nil {
		s.logger.Error("database error", slog.Any("err", err))
		return nil, err
	}
	return &UserCountJoinedThisWeekResponseIDO{Count: count}, nil
}
