package service

import (
	"errors"
	"fmt"
	"log/slog"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/common/httperror"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/config/constants"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/domain"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/usecase"
)

type UserProfileVerificationJudgeOperationService struct {
	logger             *slog.Logger
	userGetByIDUseCase *usecase.UserGetByIDUseCase
	userUpdateUseCase  *usecase.UserUpdateUseCase
}

func NewUserProfileVerificationJudgeOperationService(
	logger *slog.Logger,
	uc1 *usecase.UserGetByIDUseCase,
	uc2 *usecase.UserUpdateUseCase,
) *UserProfileVerificationJudgeOperationService {
	return &UserProfileVerificationJudgeOperationService{logger, uc1, uc2}
}

type UserProfileVerificationJudgeOperationRequestIDO struct {
	UserID                    primitive.ObjectID `bson:"user_id" json:"user_id,omitempty"`
	ProfileVerificationStatus int8               `bson:"profile_verification_status" json:"profile_verification_status,omitempty"`
}

func (s *UserProfileVerificationJudgeOperationService) Execute(sessCtx mongo.SessionContext, req *UserProfileVerificationJudgeOperationRequestIDO) (*domain.User, error) {
	//
	// STEP 1: Validation
	//

	e := make(map[string]string)

	if req == nil {
		err := errors.New("No request data inputted")
		s.logger.Error("validation error", slog.Any("err", err))
		return nil, err
	}

	if req.UserID.IsZero() {
		e["user_id"] = "User ID is required"
	}
	if req.ProfileVerificationStatus == 0 {
		e["profile_verification_status"] = "Profile verification status is required"
	}
	if len(e) != 0 {
		s.logger.Warn("validation failure",
			slog.Any("e", e))
		return nil, httperror.NewForBadRequest(&e)
	}

	// Extract from our session the following data.
	userID := sessCtx.Value(constants.SessionUserID).(primitive.ObjectID)
	ipAddress, _ := sessCtx.Value(constants.SessionIPAddress).(string)

	// Lookup the user in our database, else return a `400 Bad Request` error.
	ou, err := s.userGetByIDUseCase.Execute(sessCtx, req.UserID)
	if err != nil {
		s.logger.Error("database error", slog.Any("err", err))
		return nil, err
	}
	if ou == nil {
		s.logger.Warn("user does not exist validation error")
		return nil, httperror.NewForBadRequestWithSingleField("id", "does not exist")
	}
	adminUser, err := s.userGetByIDUseCase.Execute(sessCtx, userID)
	if err != nil {
		s.logger.Error("database error", slog.Any("err", err))
		return nil, err
	}
	if adminUser == nil {
		s.logger.Warn("Admin user does not exist validation error")
		return nil, httperror.NewForBadRequestWithSingleField("id", "does not exist")
	}

	// Verify profile status.
	if ou.ProfileVerificationStatus != domain.UserProfileVerificationStatusSubmittedForReview {
		s.logger.Warn("Profile verification status must by submitted for review")
		return nil, httperror.NewForBadRequestWithSingleField("message", "You cannot submit verification judgement")
	}

	ou.ProfileVerificationStatus = req.ProfileVerificationStatus
	ou.ModifiedByUserID = userID
	ou.ModifiedAt = time.Now()
	ou.ModifiedByName = fmt.Sprintf("%s %s", adminUser.FirstName, adminUser.LastName)
	ou.ModifiedFromIPAddress = ipAddress

	if err := s.userUpdateUseCase.Execute(sessCtx, ou); err != nil {
		s.logger.Error("user update by id error", slog.Any("error", err))
		return nil, err
	}

	return ou, nil
}
