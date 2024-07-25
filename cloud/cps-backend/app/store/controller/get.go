package controller

import (
	"context"

	domain "github.com/LuchaComics/monorepo/cloud/cps-backend/app/store/datastore"
	user_d "github.com/LuchaComics/monorepo/cloud/cps-backend/app/user/datastore"
	"github.com/LuchaComics/monorepo/cloud/cps-backend/config/constants"
	"github.com/LuchaComics/monorepo/cloud/cps-backend/utils/httperror"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log/slog"
)

func (c *StoreControllerImpl) GetByID(ctx context.Context, id primitive.ObjectID) (*domain.Store, error) {
	// Extract from our session the following data.
	userStoreID := ctx.Value(constants.SessionUserStoreID).(primitive.ObjectID)
	userRole := ctx.Value(constants.SessionUserRole).(int8)

	// If user is not administrator nor belongs to the store then error.
	if userRole != user_d.UserRoleRoot && id != userStoreID {
		c.Logger.Error("authenticated user is not staff role nor belongs to the store error",
			slog.Any("userRole", userRole),
			slog.Any("userStoreID", userStoreID))
		return nil, httperror.NewForForbiddenWithSingleField("message", "you do not belong to this store")
	}

	// Retrieve from our database the record for the specific id.
	m, err := c.StoreStorer.GetByID(ctx, id)
	if err != nil {
		c.Logger.Error("database get by id error", slog.Any("error", err))
		return nil, err
	}
	return m, err
}
