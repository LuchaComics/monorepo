package controller

import (
	"context"
	"log/slog"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	user_s "github.com/LuchaComics/monorepo/cloud/cps-pinws-backend/app/user/datastore"
	"github.com/LuchaComics/monorepo/cloud/cps-pinws-backend/config/constants"
	"github.com/LuchaComics/monorepo/cloud/cps-pinws-backend/utils/httperror"
)

func (impl *UserControllerImpl) DeleteByID(ctx context.Context, id primitive.ObjectID) error {
	////
	//// Start the transaction.
	////

	session, err := impl.DbClient.StartSession()
	if err != nil {
		impl.Logger.Error("start session error",
			slog.Any("error", err))
		return err
	}
	defer session.EndSession(ctx)

	// Define a transaction function with a series of operations
	transactionFunc := func(sessCtx mongo.SessionContext) (interface{}, error) {
		// Extract from our session the following data.
		userRole := sessCtx.Value(constants.SessionUserRole).(int8)

		// Apply filtering based on ownership and role.
		if userRole != user_s.UserRoleRoot {
			return nil, httperror.NewForForbiddenWithSingleField("message", "you do not have permission")
		}

		// STEP 1: Lookup the record or error.
		user, err := impl.UserStorer.GetByID(sessCtx, id)
		if err != nil {
			impl.Logger.Error("database get by id error", slog.Any("error", err))
			return nil, err
		}
		if user == nil {
			impl.Logger.Error("database returns nothing from get by id")
			return nil, err
		}

		// Security: Prevent deletion of root user(s).
		if user.Role == user_s.UserRoleRoot {
			impl.Logger.Warn("root user(s) cannot be deleted error")
			return nil, httperror.NewForForbiddenWithSingleField("role", "root user(s) cannot be deleted")
		}

		// STEP 2: Delete from database.
		if err := impl.UserStorer.DeleteByID(sessCtx, id); err != nil {
			impl.Logger.Error("database delete by id error", slog.Any("error", err))
			return nil, err
		}
		return nil, nil
	}

	// Start a transaction

	if _, err := session.WithTransaction(ctx, transactionFunc); err != nil {
		impl.Logger.Error("session failed error",
			slog.Any("error", err))
		return err
	}

	return nil
}
