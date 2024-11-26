package service

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/common/httperror"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/common/security/jwt"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/common/security/password"
	sstring "github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/common/security/securestring"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/common/storage/database/mongodbcache"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/domain"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/usecase"
)

type GatewayLoginService struct {
	logger                *slog.Logger
	passwordProvider      password.Provider
	cache                 mongodbcache.Cacher
	jwtProvider           jwt.Provider
	tenantGetByIDUseCase  *usecase.TenantGetByIDUseCase
	userGetByEmailUseCase *usecase.UserGetByEmailUseCase
	userUpdateUseCase     *usecase.UserUpdateUseCase
}

func NewGatewayLoginService(
	logger *slog.Logger,
	pp password.Provider,
	cach mongodbcache.Cacher,
	jwtp jwt.Provider,
	uc1 *usecase.TenantGetByIDUseCase,
	uc2 *usecase.UserGetByEmailUseCase,
	uc3 *usecase.UserUpdateUseCase,
) *GatewayLoginService {
	return &GatewayLoginService{logger, pp, cach, jwtp, uc1, uc2, uc3}
}

func (s *GatewayLoginService) Execute(sessCtx mongo.SessionContext, email string, password *sstring.SecureString) (*LoginResponseIDO, error) {
	//
	// STEP 1: Sanization of input.
	//

	// Defensive Code: For security purposes we need to perform some sanitization on the inputs.
	email = strings.ToLower(email)
	email = strings.ReplaceAll(email, " ", "")
	email = strings.ReplaceAll(email, "\t", "")
	email = strings.TrimSpace(email)
	unsecurePassword := password.String()
	unsecurePassword = strings.ReplaceAll(unsecurePassword, " ", "")
	unsecurePassword = strings.ReplaceAll(unsecurePassword, "\t", "")
	unsecurePassword = strings.TrimSpace(unsecurePassword)
	password, err := sstring.NewSecureString(unsecurePassword)
	if err != nil {
		s.logger.Error("secure string error", slog.Any("err", err))
		return nil, err
	}

	//
	// STEP 2: Validation of input.
	//

	e := make(map[string]string)
	if email == "" {
		e["email"] = "missing value"
	}
	if password == nil {
		e["password"] = "missing value"
	}

	if len(e) != 0 {
		s.logger.Warn("Failed validation login",
			slog.Any("error", e))
		return nil, httperror.NewForBadRequest(&e)
	}

	//
	// STEP 3:
	//

	// Lookup the user in our database, else return a `400 Bad Request` error.
	u, err := s.userGetByEmailUseCase.Execute(sessCtx, email)
	if err != nil {
		s.logger.Error("database error", slog.Any("err", err))
		return nil, err
	}
	if u == nil {
		s.logger.Warn("user does not exist validation error")
		return nil, httperror.NewForBadRequestWithSingleField("email", "does not exist")
	}

	// Lookup the store and check to see if it's active or not, if not active then return the specific requests.
	t, err := s.tenantGetByIDUseCase.Execute(sessCtx, u.TenantID)
	if err != nil {
		s.logger.Error("database error", slog.Any("err", err))
		return nil, err
	}
	if t == nil {
		err := fmt.Errorf("Tenant does not exist for ID: %v", u.TenantID.Hex())
		s.logger.Error("database error", slog.Any("err", err))
		return nil, err
	}

	// Verify the inputted password and hashed password match.
	passwordMatch, _ := s.passwordProvider.ComparePasswordAndHash(password, u.PasswordHash)
	if passwordMatch == false {
		s.logger.Warn("password check validation error")
		return nil, httperror.NewForBadRequestWithSingleField("password", "password do not match with record")
	}

	// Enforce the verification code of the email.
	if u.WasEmailVerified == false {
		s.logger.Warn("email verification validation error", slog.Any("u", u))
		return nil, httperror.NewForBadRequestWithSingleField("email", "was not verified")
	}

	// // Enforce 2FA if enabled.
	if u.OTPEnabled {
		// We need to reset the `otp_validated` status to be false to force
		// the user to use their `totp authenticator` application.
		u.OTPValidated = false
		u.ModifiedAt = time.Now()
		if err := s.userUpdateUseCase.Execute(sessCtx, u); err != nil {
			s.logger.Error("failed updating user during login",
				slog.Any("err", err))
			return nil, err
		}
	}

	return s.loginWithUser(sessCtx, u)
}

type LoginResponseIDO struct {
	User                   *domain.User `json:"user"`
	AccessToken            string       `json:"access_token"`
	AccessTokenExpiryTime  time.Time    `json:"access_token_expiry_time"`
	RefreshToken           string       `json:"refresh_token"`
	RefreshTokenExpiryTime time.Time    `json:"refresh_token_expiry_time"`
}

func (s *GatewayLoginService) loginWithUser(sessCtx mongo.SessionContext, u *domain.User) (*LoginResponseIDO, error) {
	uBin, err := json.Marshal(u)
	if err != nil {
		s.logger.Error("marshalling error", slog.Any("err", err))
		return nil, err
	}

	// Set expiry duration.
	atExpiry := 24 * time.Hour
	rtExpiry := 14 * 24 * time.Hour

	// Start our session using an access and refresh token.
	sessionUUID := primitive.NewObjectID().Hex()

	err = s.cache.SetWithExpiry(sessCtx, sessionUUID, uBin, rtExpiry)
	if err != nil {
		s.logger.Error("cache set with expiry error", slog.Any("err", err))
		return nil, err
	}

	// Generate our JWT token.
	accessToken, accessTokenExpiry, refreshToken, refreshTokenExpiry, err := s.jwtProvider.GenerateJWTTokenPair(sessionUUID, atExpiry, rtExpiry)
	if err != nil {
		s.logger.Error("jwt generate pairs error", slog.Any("err", err))
		return nil, err
	}

	// Return our auth keys.
	return &LoginResponseIDO{
		User:                   u,
		AccessToken:            accessToken,
		AccessTokenExpiryTime:  accessTokenExpiry,
		RefreshToken:           refreshToken,
		RefreshTokenExpiryTime: refreshTokenExpiry,
	}, nil
}
