package controller

import (
	"context"

	"log/slog"

	domain "github.com/LuchaComics/monorepo/cloud/cps-pinws-backend/app/tenant/datastore"
	user_d "github.com/LuchaComics/monorepo/cloud/cps-pinws-backend/app/user/datastore"
	"github.com/LuchaComics/monorepo/cloud/cps-pinws-backend/config/constants"
	"github.com/LuchaComics/monorepo/cloud/cps-pinws-backend/utils/httperror"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (c *TenantControllerImpl) ListByFilter(ctx context.Context, f *domain.TenantPaginationListFilter) (*domain.TenantPaginationListResult, error) {
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

	c.Logger.Debug("fetching tenants now...", slog.Any("userID", userID))
	c.Logger.Debug("listing using filter options:",
		slog.Any("TenantID", f.TenantID),
		slog.Any("Cursor", f.Cursor),
		slog.Int64("PageSize", f.PageSize),
		slog.String("SortField", f.SortField),
		slog.Int("SortOrder", int(f.SortOrder)),
		slog.Any("Status", f.Status),
		slog.Time("CreatedAtGTE", f.CreatedAtGTE),
		slog.String("SearchText", f.SearchText),
		slog.Bool("ExcludeArchived", f.ExcludeArchived))

	m, err := c.TenantStorer.ListByFilter(ctx, f)
	if err != nil {
		c.Logger.Error("database list by filter error", slog.Any("error", err))
		return nil, err
	}
	c.Logger.Debug("fetched tenants", slog.Any("m", m))
	return m, err
}

func (c *TenantControllerImpl) ListAsSelectOptionByFilter(ctx context.Context, f *domain.TenantPaginationListFilter) ([]*domain.TenantAsSelectOption, error) {
	// // Extract from our session the following data.
	// userID := ctx.Value(constants.SessionUserID).(primitive.ObjectID)
	// userRole := ctx.Value(constants.SessionUserRole).(int8)
	//
	// // // Apply protection based on ownership and role.
	// // if userRole != user_d.UserRoleRoot {
	// // 	c.Logger.Error("authenticated user is not staff role error",
	// // 		slog.Any("role", userRole),
	// // 		slog.Any("userID", userID))
	// // 	return nil, httperror.NewForForbiddenWithSingleField("message", "you role does not grant you access to this")
	// // }

	c.Logger.Debug("fetching tenants now...", slog.Any("f", f))

	m, err := c.TenantStorer.ListAsSelectOptionByFilter(ctx, f)
	if err != nil {
		c.Logger.Error("database list by filter error", slog.Any("error", err))
		return nil, err
	}
	c.Logger.Debug("fetched tenants", slog.Any("m", m))
	return m, err
}

func (c *TenantControllerImpl) PublicListAsSelectOptionByFilter(ctx context.Context, f *domain.TenantPaginationListFilter) ([]*domain.TenantAsSelectOption, error) {
	// // Extract from our session the following data.
	// userID := ctx.Value(constants.SessionUserID).(primitive.ObjectID)
	// userRole := ctx.Value(constants.SessionUserRole).(int8)
	//
	// // // Apply protection based on ownership and role.
	// // if userRole != user_d.UserRoleRoot {
	// // 	c.Logger.Error("authenticated user is not staff role error",
	// // 		slog.Any("role", userRole),
	// // 		slog.Any("userID", userID))
	// // 	return nil, httperror.NewForForbiddenWithSingleField("message", "you role does not grant you access to this")
	// // }

	c.Logger.Debug("fetching tenants now...", slog.Any("f", f))

	m, err := c.TenantStorer.ListAsSelectOptionByFilter(ctx, f)
	if err != nil {
		c.Logger.Error("database list by filter error", slog.Any("error", err))
		return nil, err
	}
	c.Logger.Debug("fetched tenants", slog.Any("m", m))
	return m, err
}
