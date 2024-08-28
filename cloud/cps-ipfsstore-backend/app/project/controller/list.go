package controller

import (
	"context"
	"log/slog"

	"go.mongodb.org/mongo-driver/bson/primitive"

	u_s "github.com/LuchaComics/monorepo/cloud/cps-ipfsstore-backend/app/user/datastore"
	domain "github.com/LuchaComics/monorepo/cloud/cps-ipfsstore-backend/app/project/datastore"
	"github.com/LuchaComics/monorepo/cloud/cps-ipfsstore-backend/config/constants"
)

func (c *ProjectControllerImpl) ListByFilter(ctx context.Context, f *domain.ProjectPaginationListFilter) (*domain.ProjectPaginationListResult, error) {
	// // Extract from our session the following data.
	// userID := ctx.Value(constants.SessionUserID).(primitive.ObjectID)
	userOID := ctx.Value(constants.SessionUserTenantID).(primitive.ObjectID)
	userRole := ctx.Value(constants.SessionUserRole).(int8)
	//
	// Apply filtering based on ownership and role.
	if userRole != u_s.UserRoleRoot {
		f.TenantID = userOID
	}

	c.Logger.Debug("listing using filter options:",
		slog.Any("TenantID", f.TenantID),
		slog.Any("Cursor", f.Cursor),
		slog.Int64("PageSize", f.PageSize),
		slog.String("SortField", f.SortField),
		slog.Int("SortOrder", int(f.SortOrder)),
		slog.String("SearchText", f.SearchText),
		slog.Bool("ExcludeArchived", f.ExcludeArchived))

	m, err := c.ProjectStorer.ListByFilter(ctx, f)
	if err != nil {
		c.Logger.Error("database list by filter error", slog.Any("error", err))
		return nil, err
	}
	return m, err
}

func (c *ProjectControllerImpl) ListAsSelectOptionByFilter(ctx context.Context, f *domain.ProjectPaginationListFilter) ([]*domain.ProjectAsSelectOption, error) {
	// Developers note: We want this unrestricted to account.

	c.Logger.Debug("listing using filter options:",
		slog.Any("TenantID", f.TenantID))

	// Filtering the database.
	m, err := c.ProjectStorer.ListAsSelectOptionByFilter(ctx, f)
	if err != nil {
		c.Logger.Error("database list by filter error", slog.Any("error", err))
		return nil, err
	}
	return m, err
}
