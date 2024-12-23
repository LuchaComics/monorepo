package service

import (
	"log/slog"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/common/httperror"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/domain"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/usecase"
)

type UserListByFilterService struct {
	logger                             *slog.Logger
	cloudStoragePresignedURLUseCase    *usecase.CloudStoragePresignedURLUseCase
	comicSubmissionListByFilterUseCase *usecase.UserListByFilterUseCase
}

func NewUserListByFilterService(
	logger *slog.Logger,
	uc1 *usecase.CloudStoragePresignedURLUseCase,
	uc2 *usecase.UserListByFilterUseCase,
) *UserListByFilterService {
	return &UserListByFilterService{logger, uc1, uc2}
}

type UserFilterRequestID domain.UserFilter

type UserFilterResultResponseIDO domain.UserFilterResult

func (s *UserListByFilterService) Execute(sessCtx mongo.SessionContext, filter *UserFilterRequestID) (*UserFilterResultResponseIDO, error) {
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

	filter2 := (*domain.UserFilter)(filter)

	// Lookup the user in our database, else return a `400 Bad Request` error.
	listResp, err := s.comicSubmissionListByFilterUseCase.Execute(sessCtx, filter2)
	if err != nil {
		s.logger.Error("database error",
			slog.Any("err", err))
		return nil, err
	}

	// s.logger.Debug("fetched",
	// 	slog.Any("id", id),
	// 	slog.Any("detail", detail))

	return (*UserFilterResultResponseIDO)(listResp), nil
}
