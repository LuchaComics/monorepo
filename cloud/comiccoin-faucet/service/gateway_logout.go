package service

import (
	"log/slog"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/common/storage/database/mongodbcache"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/config/constants"
)

type GatewayLogoutService struct {
	logger *slog.Logger
	cache  mongodbcache.Cacher
}

func NewGatewayLogoutService(
	logger *slog.Logger,
	cach mongodbcache.Cacher,
) *GatewayLogoutService {
	return &GatewayLogoutService{logger, cach}
}

func (s *GatewayLogoutService) Execute(sessCtx mongo.SessionContext) error {
	// Extract from our session the following data.
	sessionID := sessCtx.Value(constants.SessionID).(string)

	if err := s.cache.Delete(sessCtx, sessionID); err != nil {
		s.logger.Error("cache delete error", slog.Any("err", err))
		return err
	}
	return nil
}
