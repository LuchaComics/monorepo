package controller

import (
	"context"
	"log/slog"

	"go.mongodb.org/mongo-driver/bson/primitive"

	submission_s "github.com/LuchaComics/monorepo/cloud/cps-backend/app/comicsub/datastore"
	user_d "github.com/LuchaComics/monorepo/cloud/cps-backend/app/user/datastore"
	"github.com/LuchaComics/monorepo/cloud/cps-backend/config/constants"
)

func (c *ComicSubmissionControllerImpl) ListByFilter(ctx context.Context, f *submission_s.ComicSubmissionPaginationListFilter) (*submission_s.ComicSubmissionPaginationListResult, error) {
	// Extract from our session the following data.
	storeID := ctx.Value(constants.SessionUserStoreID).(primitive.ObjectID)
	userID := ctx.Value(constants.SessionUserID).(primitive.ObjectID)
	userRole := ctx.Value(constants.SessionUserRole).(int8)

	// Apply filtering based on tenancy if the user is not a system administrator.
	if userRole != user_d.UserRoleRoot {
		f.StoreID = storeID
		c.Logger.Debug("applying security policy to filters",
			slog.Any("store_id", storeID),
			slog.Any("user_id", userID),
			slog.Any("user_role", userRole))
	}

	c.Logger.Debug("listing using filter options:",
		slog.Any("StoreID", f.StoreID),
		slog.Any("Cursor", f.Cursor),
		slog.Int64("PageSize", f.PageSize),
		slog.String("SortField", f.SortField),
		slog.Int("SortOrder", int(f.SortOrder)),
		slog.Any("Status", f.Status),
		slog.Time("CreatedAtGTE", f.CreatedAtGTE),
		slog.String("SearchText", f.SearchText),
		slog.Bool("ExcludeArchived", f.ExcludeArchived))

	m, err := c.ComicSubmissionStorer.ListByFilter(ctx, f)
	if err != nil {
		c.Logger.Error("database list by filter error", slog.Any("error", err))
		return nil, err
	}
	return m, err
}
