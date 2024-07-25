package controller

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"log/slog"

	org_d "github.com/LuchaComics/monorepo/cloud/cps-backend/app/store/datastore"
	user_d "github.com/LuchaComics/monorepo/cloud/cps-backend/app/user/datastore"
	"github.com/LuchaComics/monorepo/cloud/cps-backend/config/constants"
	"github.com/LuchaComics/monorepo/cloud/cps-backend/utils/httperror"
)

func (impl *StoreControllerImpl) DeleteByID(ctx context.Context, id primitive.ObjectID) error {
	// Extract from our session the following data.
	userID := ctx.Value(constants.SessionUserID).(primitive.ObjectID)
	userRole := ctx.Value(constants.SessionUserRole).(int8)

	// Apply protection based on ownership and role.
	if userRole != user_d.UserRoleRoot {
		impl.Logger.Error("authenticated user is not staff role error",
			slog.Any("role", userRole),
			slog.Any("userID", userID))
		return httperror.NewForForbiddenWithSingleField("message", "you role does not grant you access to this")
	}

	// Update the database.
	store, err := impl.GetByID(ctx, id)
	store.Status = org_d.StoreArchivedStatus
	if err != nil {
		impl.Logger.Error("database get by id error", slog.Any("error", err))
		return err
	}
	if store == nil {
		impl.Logger.Error("database returns nothing from get by id")
		return err
	}
	// Security: Prevent deletion of root user(s).
	if store.Type == org_d.RootType {
		impl.Logger.Warn("root store cannot be deleted error")
		return httperror.NewForForbiddenWithSingleField("role", "root store cannot be deleted")
	}

	// Save to the database the modified store.
	if err := impl.StoreStorer.UpdateByID(ctx, store); err != nil {
		impl.Logger.Error("database update by id error", slog.Any("error", err))
		return err
	}
	return nil
}
