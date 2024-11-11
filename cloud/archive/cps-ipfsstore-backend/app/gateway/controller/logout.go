package controller

import (
	"context"

	"log/slog"

	"github.com/LuchaComics/monorepo/cloud/cps-ipfsstore-backend/config/constants"
)

func (impl *GatewayControllerImpl) Logout(ctx context.Context) error {
	// Extract from our session the following data.
	sessionID := ctx.Value(constants.SessionID).(string)

	if err := impl.Cache.Delete(ctx, sessionID); err != nil {
		impl.Logger.Error("cache delete error", slog.Any("err", err))
		return err
	}

	return nil
}
