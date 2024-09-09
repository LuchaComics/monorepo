package controller

import (
	"context"

	"log/slog"

	domain "github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/app/nftasset/datastore"
	user_d "github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/app/user/datastore"
	"github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/config/constants"
	"github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/utils/httperror"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (c *NFTAssetControllerImpl) ListByFilter(ctx context.Context, f *domain.NFTAssetPaginationListFilter) (*domain.NFTAssetPaginationListResult, error) {
	// Extract from our session the following data.
	orgID := ctx.Value(constants.SessionUserTenantID).(primitive.ObjectID)
	userID := ctx.Value(constants.SessionUserID).(primitive.ObjectID)
	userRole := ctx.Value(constants.SessionUserRole).(int8)

	// Apply protection based on ownership and role.
	if userRole != user_d.UserRoleRoot {
		f.TenantID = orgID // Force store tenancy restrictions.
	}

	c.Logger.Debug("fetching nftassets now...", slog.Any("userID", userID))

	aa, err := c.NFTAssetStorer.ListByFilter(ctx, f)
	if err != nil {
		c.Logger.Error("database list by filter error", slog.Any("error", err))
		return nil, err
	}
	c.Logger.Debug("fetched nftassets", slog.Any("aa", aa))

	return aa, err
}

func (c *NFTAssetControllerImpl) ListAsSelectOptionByFilter(ctx context.Context, f *domain.NFTAssetPaginationListFilter) ([]*domain.NFTAssetAsSelectOption, error) {
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

	c.Logger.Debug("fetching nftassets now...", slog.Any("userID", userID))

	m, err := c.NFTAssetStorer.ListAsSelectOptionByFilter(ctx, f)
	if err != nil {
		c.Logger.Error("database list by filter error", slog.Any("error", err))
		return nil, err
	}
	c.Logger.Debug("fetched nftassets", slog.Any("m", m))
	return m, err
}
