package service

import (
	"log/slog"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/common/httperror"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/common/kmutexutil"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/common/security/password"
	sstring "github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/common/security/securestring"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/config/constants"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/usecase"
)

type GatewayChangePasswordService struct {
	logger             *slog.Logger
	kmutex             kmutexutil.KMutexProvider
	passwordProvider   password.Provider
	userGetByIDUseCase *usecase.UserGetByIDUseCase
	userUpdateUseCase  *usecase.UserUpdateUseCase
}

func NewGatewayChangePasswordService(
	logger *slog.Logger,
	kmutex kmutexutil.KMutexProvider,
	passwordProvider password.Provider,
	uc1 *usecase.UserGetByIDUseCase,
	uc2 *usecase.UserUpdateUseCase,
) *GatewayChangePasswordService {
	return &GatewayChangePasswordService{logger, kmutex, passwordProvider, uc1, uc2}
}

type GatewayChangePasswordRequestIDO struct {
	OldPassword         string `json:"old_password"`
	NewPassword         string `json:"new_password"`
	NewPasswordRepeated string `json:"new_password_repeated"`
}

func (s *GatewayChangePasswordService) Execute(sessCtx mongo.SessionContext, req *GatewayChangePasswordRequestIDO) error {
	// Extract from our session the following data.
	userID := sessCtx.Value(constants.SessionUserID).(primitive.ObjectID)

	// Lookup the user in our database, else return a `400 Bad Request` error.
	u, err := s.userGetByIDUseCase.Execute(sessCtx, userID)
	if err != nil {
		s.logger.Error("database error", slog.Any("err", err))
		return err
	}
	if u == nil {
		s.logger.Warn("user does not exist validation error")
		return httperror.NewForBadRequestWithSingleField("id", "does not exist")
	}
	if err := ValidateProfileChangePassworRequest(req); err != nil {
		s.logger.Warn("user validation failed", slog.Any("err", err))
		return err
	}

	oldPass, err := sstring.NewSecureString(req.OldPassword)
	if err != nil {
		s.logger.Error("Failed to secure old password string",
			slog.Any("err", err))
		return err
	}
	newPass, err := sstring.NewSecureString(req.NewPassword)
	if err != nil {
		s.logger.Error("Failed to secure new password string",
			slog.Any("err", err))
		return err
	}

	// Verify the inputted password and hashed password match.
	if passwordMatch, _ := s.passwordProvider.ComparePasswordAndHash(oldPass, u.PasswordHash); passwordMatch == false {
		s.logger.Warn("password check validation error")
		return httperror.NewForBadRequestWithSingleField("old_password", "old password do not match with record of existing password")
	}

	passwordHash, err := s.passwordProvider.GenerateHashFromPassword(newPass)
	if err != nil {
		s.logger.Error("hashing error", slog.Any("error", err))
		return err
	}
	u.PasswordHash = passwordHash
	u.PasswordHashAlgorithm = s.passwordProvider.AlgorithmName()
	if err := s.userUpdateUseCase.Execute(sessCtx, u); err != nil {
		s.logger.Error("user update by id error", slog.Any("error", err))
		return err
	}

	return nil
}

func ValidateProfileChangePassworRequest(dirtyData *GatewayChangePasswordRequestIDO) error {
	e := make(map[string]string)

	if dirtyData.OldPassword == "" {
		e["old_password"] = "missing value"
	}
	if dirtyData.NewPassword == "" {
		e["new_password"] = "missing value"
	}
	if dirtyData.NewPasswordRepeated == "" {
		e["new_password_repeated"] = "missing value"
	}
	if dirtyData.NewPasswordRepeated != dirtyData.NewPassword {
		e["new_password"] = "does not match"
		e["new_password_repeated"] = "does not match"
	}

	if len(e) != 0 {
		return httperror.NewForBadRequest(&e)
	}
	return nil
}
