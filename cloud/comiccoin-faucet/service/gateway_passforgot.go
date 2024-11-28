package service

import (
	"log/slog"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/common/httperror"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/common/templatedemailer"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/usecase"
)

type GatewayForgotPasswordService struct {
	logger                *slog.Logger
	templatedEmailer      templatedemailer.TemplatedEmailer
	userGetByEmailUseCase *usecase.UserGetByEmailUseCase
	userUpdateUseCase     *usecase.UserUpdateUseCase
}

func NewGatewayForgotPasswordService(
	logger *slog.Logger,
	templatedEmailer templatedemailer.TemplatedEmailer,
	uc1 *usecase.UserGetByEmailUseCase,
	uc2 *usecase.UserUpdateUseCase,
) *GatewayForgotPasswordService {
	return &GatewayForgotPasswordService{logger, templatedEmailer, uc1, uc2}
}

type GatewayForgotPasswordRequestIDO struct {
	Email string `json:"email"`
}

func (s *GatewayForgotPasswordService) Execute(sessCtx mongo.SessionContext, req *GatewayForgotPasswordRequestIDO) error {
	// // Extract from our session the following data.
	// userID := sessCtx.Value(constants.SessionUserID).(primitive.ObjectID)

	e := make(map[string]string)
	if req.Email == "" {
		e["email"] = "missing value"
	}
	if len(e) != 0 {
		return httperror.NewForBadRequest(&e)
	}

	// Lookup the user in our database, else return a `400 Bad Request` error.
	u, err := s.userGetByEmailUseCase.Execute(sessCtx, req.Email)
	if err != nil {
		s.logger.Error("database error", slog.Any("err", err))
		return err
	}
	if u == nil {
		s.logger.Warn("user does not exist validation error")
		return httperror.NewForBadRequestWithSingleField("id", "does not exist")
	}

	// Generate unique token and save it to the user record.
	u.EmailVerificationCode = primitive.NewObjectID().Hex()
	if err := s.userUpdateUseCase.Execute(sessCtx, u); err != nil {
		s.logger.Warn("user update by id failed", slog.Any("error", err))
		return err
	}

	// Send password reset email.
	return s.templatedEmailer.SendForgotPasswordEmail(req.Email, u.EmailVerificationCode, u.FirstName)
}
