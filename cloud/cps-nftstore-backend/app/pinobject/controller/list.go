package controller

import (
	"context"
	"time"

	domain "github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/app/pinobject/datastore"
	user_d "github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/app/user/datastore"
	"github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/config/constants"
	"github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/utils/httperror"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log/slog"
)

func (c *PinObjectControllerImpl) ListByFilter(ctx context.Context, f *domain.PinObjectPaginationListFilter) (*domain.PinObjectPaginationListResult, error) {
	// Extract from our session the following data.
	orgID := ctx.Value(constants.SessionUserTenantID).(primitive.ObjectID)
	userID := ctx.Value(constants.SessionUserID).(primitive.ObjectID)
	userRole := ctx.Value(constants.SessionUserRole).(int8)

	// Apply protection based on ownership and role.
	if userRole != user_d.UserRoleRoot {
		f.TenantID = orgID // Force store tenancy restrictions.
	}

	c.Logger.Debug("fetching pinobjects now...", slog.Any("userID", userID))

	aa, err := c.PinObjectStorer.ListByFilter(ctx, f)
	if err != nil {
		c.Logger.Error("database list by filter error", slog.Any("error", err))
		return nil, err
	}
	c.Logger.Debug("fetched pinobjects", slog.Any("aa", aa))

	for _, a := range aa.Results {
		// Generate the URL.
		fileURL, err := c.S3.GetPresignedURL(ctx, a.ObjectKey, 5*time.Minute)
		if err != nil {
			c.Logger.Error("s3 failed get presigned url error", slog.Any("error", err))
			return nil, err
		}
		a.ObjectURL = fileURL
	}
	return aa, err
}

func (c *PinObjectControllerImpl) ListAsSelectOptionByFilter(ctx context.Context, f *domain.PinObjectPaginationListFilter) ([]*domain.PinObjectAsSelectOption, error) {
	// Extract from our session the following data.
	userID := ctx.Value(constants.SessionUserID).(primitive.ObjectID)
	userRole := ctx.Value(constants.SessionUserRole).(int8)

	// Apply protection based on ownership and role.
	if userRole != user_d.UserRoleRoot {
		c.Logger.Error("authenticated user is not staff role error",
			slog.Any("role", userRole),
			slog.Any("userID", userID))
		return nil, httperror.NewForForbiddenWithSingleField("message", "you role does not grant you access to this")
	}

	c.Logger.Debug("fetching pinobjects now...", slog.Any("userID", userID))

	m, err := c.PinObjectStorer.ListAsSelectOptionByFilter(ctx, f)
	if err != nil {
		c.Logger.Error("database list by filter error", slog.Any("error", err))
		return nil, err
	}
	c.Logger.Debug("fetched pinobjects", slog.Any("m", m))
	return m, err
}
