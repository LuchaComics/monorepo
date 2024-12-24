package service

import (
	"log/slog"
	"time"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/common/httperror"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/common/kmutexutil"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/config/constants"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/domain"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/usecase"
)

type GatewayVerifyEmailService struct {
	logger                           *slog.Logger
	kmutex                           kmutexutil.KMutexProvider
	userGetByVerificationCodeUseCase *usecase.UserGetByVerificationCodeUseCase
	userUpdateUseCase                *usecase.UserUpdateUseCase
}

func NewGatewayVerifyEmailService(
	logger *slog.Logger,
	kmutex kmutexutil.KMutexProvider,
	uc1 *usecase.UserGetByVerificationCodeUseCase,
	uc2 *usecase.UserUpdateUseCase,
) *GatewayVerifyEmailService {
	return &GatewayVerifyEmailService{logger, kmutex, uc1, uc2}
}

type GatewayVerifyRequestIDO struct {
	Code string `json:"code"`
}

type GatwayVerifyResponseIDO struct {
	Message  string `json:"message"`
	UserRole int8   `bson:"user_role" json:"user_role"`
}

func (s *GatewayVerifyEmailService) Execute(sessCtx mongo.SessionContext, req *GatewayVerifyRequestIDO) (*GatwayVerifyResponseIDO, error) {
	s.kmutex.Acquire(req.Code)
	defer func() {
		s.kmutex.Release(req.Code)
	}()

	// // Extract from our session the following data.
	// sessionID := sessCtx.Value(constants.SessionID).(string)

	res := &GatwayVerifyResponseIDO{}

	// Lookup the user in our database, else return a `400 Bad Request` error.
	u, err := s.userGetByVerificationCodeUseCase.Execute(sessCtx, req.Code)
	if err != nil {
		s.logger.Error("database error", slog.Any("err", err))
		return nil, err
	}
	if u == nil {
		s.logger.Warn("user does not exist validation error")
		return nil, httperror.NewForBadRequestWithSingleField("code", "does not exist")
	}

	//TODO: Handle expiry dates.

	// Extract from our session the following data.
	// userID := sessCtx.Value(constants.SessionUserID).(primitive.ObjectID)
	ipAddress, _ := sessCtx.Value(constants.SessionIPAddress).(string)

	// Verify the user.
	u.WasEmailVerified = true
	// ou.ModifiedByUserID = userID
	u.ModifiedAt = time.Now()
	// ou.ModifiedByName = fmt.Sprintf("%s %s", ou.FirstName, ou.LastName)
	u.ModifiedFromIPAddress = ipAddress
	if err := s.userUpdateUseCase.Execute(sessCtx, u); err != nil {
		s.logger.Error("update error", slog.Any("err", err))
		return nil, err
	}

	//
	// Send notification based on user role
	//

	switch u.Role {
	case domain.UserRoleCustomer:
		{
			res.Message = "Thank you for verifying. You may log in now to get started!"
			s.logger.Debug("customer user verified")
			break
		}
	default:
		{
			res.Message = "Thank you for verifying. You may log in now to get started!"
			s.logger.Debug("unknown user verified")
			break
		}
	}
	res.UserRole = u.Role

	return res, nil
}
