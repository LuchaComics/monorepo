package controller

import (
	"context"

	"log/slog"

	"go.mongodb.org/mongo-driver/bson/primitive"

	domain "github.com/LuchaComics/monorepo/cloud/cps-backend/app/credit/datastore"
	u_s "github.com/LuchaComics/monorepo/cloud/cps-backend/app/user/datastore"
	"github.com/LuchaComics/monorepo/cloud/cps-backend/config/constants"
)

func (c *CreditControllerImpl) ListByFilter(ctx context.Context, f *domain.CreditPaginationListFilter) (*domain.CreditPaginationListResult, error) {
	// // Extract from our session the following data.
	// userID := ctx.Value(constants.SessionUserID).(primitive.ObjectID)
	userOID := ctx.Value(constants.SessionUserStoreID).(primitive.ObjectID)
	userRole := ctx.Value(constants.SessionUserRole).(int8)
	//
	// Apply filtering based on ownership and role.
	if userRole != u_s.UserRoleRoot {
		f.StoreID = userOID
	}

	c.Logger.Debug("listing using filter options:",
		slog.Any("StoreID", f.StoreID),
		slog.Any("UserID", f.UserID),
		slog.Any("Cursor", f.Cursor),
		slog.Int64("PageSize", f.PageSize),
		slog.String("SortField", f.SortField),
		slog.Int("SortOrder", int(f.SortOrder)),
		slog.String("SearchText", f.SearchText))

	m, err := c.CreditStorer.ListByFilter(ctx, f)
	if err != nil {
		c.Logger.Error("database list by filter error", slog.Any("error", err))
		return nil, err
	}
	return m, err
}

func (c *CreditControllerImpl) ListAsSelectOptionByFilter(ctx context.Context, f *domain.CreditPaginationListFilter) ([]*domain.CreditAsSelectOption, error) {
	// Developers note: We want this unrestricted to account.

	c.Logger.Debug("listing using filter options:",
		slog.Any("StoreID", f.StoreID))

	// Filtering the database.
	m, err := c.CreditStorer.ListAsSelectOptionByFilter(ctx, f)
	if err != nil {
		c.Logger.Error("database list by filter error", slog.Any("error", err))
		return nil, err
	}
	return m, err
}
