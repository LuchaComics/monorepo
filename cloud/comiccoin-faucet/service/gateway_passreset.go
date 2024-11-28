package service

import (
	"log/slog"
	"time"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/common/httperror"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/common/kmutexutil"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/common/security/password"
	sstring "github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/common/security/securestring"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/usecase"
)

type GatewayPasswordResetService struct {
	logger                           *slog.Logger
	kmutex                           kmutexutil.KMutexProvider
	passwordProvider                 password.Provider
	userGetByVerificationCodeUseCase *usecase.UserGetByVerificationCodeUseCase
	userUpdateUseCase                *usecase.UserUpdateUseCase
}

func NewGatewayPasswordResetService(
	logger *slog.Logger,
	kmutex kmutexutil.KMutexProvider,
	pp password.Provider,
	uc1 *usecase.UserGetByVerificationCodeUseCase,
	uc2 *usecase.UserUpdateUseCase,
) *GatewayPasswordResetService {
	return &GatewayPasswordResetService{logger, kmutex, pp, uc1, uc2}
}

type GatewayPasswordResetRequestIDO struct {
	Code     string `json:"code"`
	Password string `json:"password"`
}

func (s *GatewayPasswordResetService) Execute(sessCtx mongo.SessionContext, req *GatewayPasswordResetRequestIDO) error {
	s.kmutex.Acquire(req.Code)
	defer func() {
		s.kmutex.Release(req.Code)
	}()

	// // Extract from our session the following data.
	// sessionID := sessCtx.Value(constants.SessionID).(string)

	// Lookup the user in our database, else return a `400 Bad Request` error.
	u, err := s.userGetByVerificationCodeUseCase.Execute(sessCtx, req.Code)
	if err != nil {
		s.logger.Error("database error", slog.Any("err", err))
		return err
	}
	if u == nil {
		s.logger.Warn("user does not exist validation error")
		return httperror.NewForBadRequestWithSingleField("code", "does not exist")
	}

	//TODO: Handle expiry dates.

	securePassword, err := sstring.NewSecureString(req.Password)
	if err != nil {
		s.logger.Error("database error", slog.Any("err", err))
		return err
	}

	passwordHash, err := s.passwordProvider.GenerateHashFromPassword(securePassword)
	if err != nil {
		s.logger.Error("hashing error", slog.Any("error", err))
		return err
	}

	// Verify the user.
	u.PasswordHash = passwordHash
	u.PasswordHashAlgorithm = s.passwordProvider.AlgorithmName()
	u.EmailVerificationCode = "" // Remove email active code so it cannot be used agian.
	u.EmailVerificationExpiry = time.Now()
	u.ModifiedAt = time.Now()
	if err := s.userUpdateUseCase.Execute(sessCtx, u); err != nil {
		s.logger.Error("update error", slog.Any("err", err))
		return err
	}

	return nil
}
