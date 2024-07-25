package controller

import (
	"context"
	"log/slog"

	user_s "github.com/LuchaComics/monorepo/cloud/cps-backend/app/user/datastore"
	"github.com/LuchaComics/monorepo/cloud/cps-backend/config/constants"
	"github.com/LuchaComics/monorepo/cloud/cps-backend/utils/httperror"
)

func (c *UserControllerImpl) ListByFilter(ctx context.Context, f *user_s.UserPaginationListFilter) (*user_s.UserPaginationListResult, error) {
	// Extract from our session the following data.
	userRole := ctx.Value(constants.SessionUserRole).(int8)

	// Apply filtering based on ownership and role.
	if userRole != user_s.UserRoleRoot {
		return nil, httperror.NewForForbiddenWithSingleField("message", "you do not have permission")
	}

	c.Logger.Debug("listing using filter options:",
		slog.Any("StoreID", f.StoreID),
		slog.Any("Cursor", f.Cursor),
		slog.Int64("PageSize", f.PageSize),
		slog.String("SortField", f.SortField),
		slog.Int("SortOrder", int(f.SortOrder)),
		slog.Any("Status", f.Status),
		slog.String("SearchText", f.SearchText),
		slog.Time("CreatedAtGTE", f.CreatedAtGTE))

	// Filtering the database.
	m, err := c.UserStorer.ListByFilter(ctx, f)
	if err != nil {
		c.Logger.Error("database list by filter error", slog.Any("error", err))
		return nil, err
	}
	return m, err
}

func (c *UserControllerImpl) ListAsSelectOptionByFilter(ctx context.Context, f *user_s.UserPaginationListFilter) ([]*user_s.UserAsSelectOption, error) {
	// Extract from our session the following data.
	userRole := ctx.Value(constants.SessionUserRole).(int8)

	// Apply filtering based on ownership and role.
	if userRole != user_s.UserRoleRoot {
		return nil, httperror.NewForForbiddenWithSingleField("message", "you do not have permission")
	}

	c.Logger.Debug("listing using filter options:",
		slog.Any("StoreID", f.StoreID),
		slog.Any("Role", f.Role))

	// Filtering the database.
	m, err := c.UserStorer.ListAsSelectOptionByFilter(ctx, f)
	if err != nil {
		c.Logger.Error("database list by filter error", slog.Any("error", err))
		return nil, err
	}
	return m, err
}
