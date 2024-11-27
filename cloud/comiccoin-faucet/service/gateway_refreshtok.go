package service

import (
	"encoding/json"
	"errors"
	"log/slog"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/common/security/jwt"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/common/storage/database/mongodbcache"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/domain"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/usecase"
)

type GatewayRefreshTokenService struct {
	logger                *slog.Logger
	cache                 mongodbcache.Cacher
	jwtProvider           jwt.Provider
	userGetByEmailUseCase *usecase.UserGetByEmailUseCase
}

func NewGatewayRefreshTokenService(
	logger *slog.Logger,
	cach mongodbcache.Cacher,
	jwtp jwt.Provider,
	uc1 *usecase.UserGetByEmailUseCase,
) *GatewayRefreshTokenService {
	return &GatewayRefreshTokenService{logger, cach, jwtp, uc1}
}

func (s *GatewayRefreshTokenService) Execute(
	sessCtx mongo.SessionContext,
	value string,
) (*domain.User, string, time.Time, string, time.Time, error) {
	////
	//// Extract the `sessionID` so we can process it.
	////

	sessionID, err := s.jwtProvider.ProcessJWTToken(value)
	if err != nil {
		s.logger.Warn("process jwt refresh token does not exist", slog.String("value", value))
		err := errors.New("jwt refresh token failed")
		return nil, "", time.Now(), "", time.Now(), err
	}

	////
	//// Lookup in our in-memory the user record for the `sessionID` or error.
	////

	uBin, err := s.cache.Get(sessCtx, sessionID)
	if err != nil {
		s.logger.Error("in-memory set error", slog.Any("err", err))
		return nil, "", time.Now(), "", time.Now(), err
	}

	var u *domain.User
	err = json.Unmarshal(uBin, &u)
	if err != nil {
		s.logger.Error("unmarshal error", slog.Any("err", err))
		return nil, "", time.Now(), "", time.Now(), err
	}

	////
	//// Generate new access and refresh tokens and return them.
	////

	// Set expiry duration.
	atExpiry := 24 * time.Hour
	rtExpiry := 14 * 24 * time.Hour

	// Start our session using an access and refresh token.
	newSessionUUID := primitive.NewObjectID().Hex()

	err = s.cache.SetWithExpiry(sessCtx, newSessionUUID, uBin, rtExpiry)
	if err != nil {
		s.logger.Error("cache set with expiry error", slog.Any("err", err))
		return nil, "", time.Now(), "", time.Now(), err
	}

	// Generate our JWT token.
	accessToken, accessTokenExpiry, refreshToken, refreshTokenExpiry, err := s.jwtProvider.GenerateJWTTokenPair(newSessionUUID, atExpiry, rtExpiry)
	if err != nil {
		s.logger.Error("jwt generate pairs error", slog.Any("err", err))
		return nil, "", time.Now(), "", time.Now(), err
	}

	// Return our auth keys.
	return u, accessToken, accessTokenExpiry, refreshToken, refreshTokenExpiry, nil
}
