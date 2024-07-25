package controller

import (
	"context"

	"log/slog"

	"go.mongodb.org/mongo-driver/bson/primitive"

	user_s "github.com/LuchaComics/monorepo/cloud/cps-backend/app/user/datastore"
	"github.com/LuchaComics/monorepo/cloud/cps-backend/config/constants"
)

func (c *CustomerControllerImpl) ListByFilter(ctx context.Context, f *user_s.UserPaginationListFilter) (*user_s.UserPaginationListResult, error) {
	// // Extract from our session the following data.
	storeID := ctx.Value(constants.SessionUserStoreID).(primitive.ObjectID)
	userRole := ctx.Value(constants.SessionUserRole).(int8)

	// Apply filtering based on ownership and role.
	if userRole == user_s.UserRoleRetailer {
		f.StoreID = storeID
	}

	f.Role = user_s.UserRoleCustomer // Manditory

	m, err := c.UserStorer.ListByFilter(ctx, f)
	if err != nil {
		c.Logger.Error("database list by filter error", slog.Any("error", err))
		return nil, err
	}
	return m, err
}
