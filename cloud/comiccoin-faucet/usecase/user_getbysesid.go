package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/common/httperror"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/common/storage/database/mongodbcache"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/config"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/domain"
)

type UserGetBySessionIDUseCase struct {
	config *config.Configuration
	logger *slog.Logger
	cache  mongodbcache.Cacher
}

func NewUserGetBySessionIDUseCase(config *config.Configuration, logger *slog.Logger, ca mongodbcache.Cacher) *UserGetBySessionIDUseCase {
	return &UserGetBySessionIDUseCase{config, logger, ca}
}

func (uc *UserGetBySessionIDUseCase) Execute(ctx context.Context, sessionID string) (*domain.User, error) {
	//
	// STEP 1: Validation.
	//

	e := make(map[string]string)
	if sessionID == "" {
		e["session_id"] = "missing value"
	} else {
		//TODO: IMPL.
	}
	if len(e) != 0 {
		uc.logger.Warn("Validation failed for upsert",
			slog.Any("error", e))
		return nil, httperror.NewForBadRequest(&e)
	}

	//
	// STEP 2:
	//

	uc.logger.Debug("gateway controller initialization started...")

	userBytes, err := uc.cache.Get(ctx, sessionID)
	if err != nil {
		return nil, err
	}
	if userBytes == nil {
		uc.logger.Warn("record not found")
		return nil, errors.New("record not found")
	}
	var user domain.User
	err = json.Unmarshal(userBytes, &user)
	if err != nil {
		uc.logger.Error("unmarshalling failed", slog.Any("err", err))
		return nil, err
	}

	uc.logger.Debug("gateway controller initialized")

	return &user, nil
}
