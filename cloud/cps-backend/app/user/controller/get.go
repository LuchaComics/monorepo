package controller

import (
	"context"
	"log/slog"

	"go.mongodb.org/mongo-driver/bson/primitive"

	domain "github.com/LuchaComics/monorepo/cloud/cps-backend/app/user/datastore"
	user_s "github.com/LuchaComics/monorepo/cloud/cps-backend/app/user/datastore"
	"github.com/LuchaComics/monorepo/cloud/cps-backend/config/constants"
	"github.com/LuchaComics/monorepo/cloud/cps-backend/utils/httperror"
)

func (c *UserControllerImpl) GetByID(ctx context.Context, id primitive.ObjectID) (*user_s.User, error) {
	// Extract from our session the following data.
	userRole, _ := ctx.Value(constants.SessionUserRole).(int8)
	userStoreID, _ := ctx.Value(constants.SessionUserStoreID).(primitive.ObjectID)

	// Retrieve from our database the record for the specific id.
	m, err := c.UserStorer.GetByID(ctx, id)
	if err != nil {
		c.Logger.Error("database get by id error", slog.Any("error", err))
		return nil, err
	}

	// Apply filtering based on ownership and role.
	if userRole != user_s.UserRoleRoot {
		if userStoreID != m.StoreID {
			c.Logger.Error("permission error",
				slog.Any("userStoreID", userStoreID),
				slog.Any("m.StoreID", m.StoreID))
			return nil, httperror.NewForForbiddenWithSingleField("message", "you do not have permission")
		}
	}

	return m, err
}

// CreateInitialRootAdmin function creates the initial root administrator if not previously created.
func (c *UserControllerImpl) GetUserBySessionUUID(ctx context.Context, sessionUUID string) (*domain.User, error) {
	panic("TODO: IMPLEMENT")
	return nil, nil
}
