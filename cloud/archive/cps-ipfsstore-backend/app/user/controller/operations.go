package controller

import (
	"context"
	"log/slog"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	user_s "github.com/LuchaComics/monorepo/cloud/cps-ipfsstore-backend/app/user/datastore"
	"github.com/LuchaComics/monorepo/cloud/cps-ipfsstore-backend/config/constants"
	"github.com/LuchaComics/monorepo/cloud/cps-ipfsstore-backend/utils/httperror"
)

func (impl *UserControllerImpl) CreateComment(ctx context.Context, customerID primitive.ObjectID, content string) (*user_s.User, error) {
	////
	//// Start the transaction.
	////

	session, err := impl.DbClient.StartSession()
	if err != nil {
		impl.Logger.Error("start session error",
			slog.Any("error", err))
		return nil, err
	}
	defer session.EndSession(ctx)

	// Define a transaction function with a series of operations
	transactionFunc := func(sessCtx mongo.SessionContext) (interface{}, error) {
		// Fetch the original customer.
		s, err := impl.UserStorer.GetByID(sessCtx, customerID)
		if err != nil {
			impl.Logger.Error("database get by id error", slog.Any("error", err))
			return nil, err
		}
		if s == nil {
			impl.Logger.Error("user does not exist error",
				slog.Any("userid", customerID))
			return nil, httperror.NewForBadRequestWithSingleField("message", "user does not exist")
		}

		// Create our comment.
		comment := &user_s.UserComment{
			ID:               primitive.NewObjectID(),
			Content:          content,
			TenantID:          sessCtx.Value(constants.SessionUserTenantID).(primitive.ObjectID),
			CreatedByUserID:  sessCtx.Value(constants.SessionUserID).(primitive.ObjectID),
			CreatedByName:    sessCtx.Value(constants.SessionUserName).(string),
			CreatedAt:        time.Now(),
			ModifiedByUserID: sessCtx.Value(constants.SessionUserID).(primitive.ObjectID),
			ModifiedByName:   sessCtx.Value(constants.SessionUserName).(string),
			ModifiedAt:       time.Now(),
		}

		// Add our comment to the comments.
		s.ModifiedByUserID = sessCtx.Value(constants.SessionUserID).(primitive.ObjectID)
		s.ModifiedAt = time.Now()
		s.Comments = append(s.Comments, comment)

		// Save to the database the modified customer.
		if err := impl.UserStorer.UpdateByID(sessCtx, s); err != nil {
			impl.Logger.Error("database update by id error", slog.Any("error", err))
			return nil, err
		}

		return s, nil
	}

	// Start a transaction
	res, err := session.WithTransaction(ctx, transactionFunc)
	if err != nil {
		impl.Logger.Error("session failed error",
			slog.Any("error", err))
		return nil, err
	}

	return res.(*user_s.User), nil
}

// Star function will either set this user's `is_starred` to `true` or `false`
// depending if it was previously starred.
func (impl *UserControllerImpl) Star(ctx context.Context, id primitive.ObjectID) (*user_s.User, error) {
	////
	//// Start the transaction.
	////

	session, err := impl.DbClient.StartSession()
	if err != nil {
		impl.Logger.Error("start session error",
			slog.Any("error", err))
		return nil, err
	}
	defer session.EndSession(ctx)

	// Define a transaction function with a series of operations
	transactionFunc := func(sessCtx mongo.SessionContext) (interface{}, error) {
		// Extract from our session the following data.
		userID, _ := sessCtx.Value(constants.SessionUserID).(primitive.ObjectID)
		userName, _ := sessCtx.Value(constants.SessionUserName).(string)

		// Extract from our session the following data.
		userRole := sessCtx.Value(constants.SessionUserRole).(int8)

		// Apply filtering based on ownership and role.
		if userRole != user_s.UserRoleRoot {
			return nil, httperror.NewForForbiddenWithSingleField("message", "you do not have permission")
		}

		// Lookup the user in our database, else return a `400 Bad Request` error.
		ou, err := impl.UserStorer.GetByID(sessCtx, id)
		if err != nil {
			impl.Logger.Error("database error", slog.Any("err", err))
			return nil, err
		}
		if ou == nil {
			impl.Logger.Warn("user does not exist validation error")
			return nil, httperror.NewForBadRequestWithSingleField("id", "does not exist")
		}

		ou.IsStarred = !ou.IsStarred
		ou.ModifiedByUserID = userID
		ou.ModifiedAt = time.Now()
		ou.ModifiedByName = userName

		if err := impl.UserStorer.UpdateByID(sessCtx, ou); err != nil {
			impl.Logger.Error("user update by id error", slog.Any("error", err))
			return nil, err
		}
		return ou, nil
	}

	// Start a transaction
	res, err := session.WithTransaction(ctx, transactionFunc)
	if err != nil {
		impl.Logger.Error("session failed error",
			slog.Any("error", err))
		return nil, err
	}

	return res.(*user_s.User), nil
}

func (impl *UserControllerImpl) ArchiveByID(ctx context.Context, id primitive.ObjectID) (*user_s.User, error) {
	////
	//// Start the transaction.
	////

	session, err := impl.DbClient.StartSession()
	if err != nil {
		impl.Logger.Error("start session error",
			slog.Any("error", err))
		return nil, err
	}
	defer session.EndSession(ctx)

	// Define a transaction function with a series of operations
	transactionFunc := func(sessCtx mongo.SessionContext) (interface{}, error) {
		// Extract from our session the following data.
		userID, _ := sessCtx.Value(constants.SessionUserID).(primitive.ObjectID)
		userName, _ := sessCtx.Value(constants.SessionUserName).(string)

		// Extract from our session the following data.
		userRole := sessCtx.Value(constants.SessionUserRole).(int8)

		// Apply filtering based on ownership and role.
		if userRole != user_s.UserRoleRoot {
			return nil, httperror.NewForForbiddenWithSingleField("message", "you do not have permission")
		}

		// Lookup the user in our database, else return a `400 Bad Request` error.
		ou, err := impl.UserStorer.GetByID(sessCtx, id)
		if err != nil {
			impl.Logger.Error("database error", slog.Any("err", err))
			return nil, err
		}
		if ou == nil {
			impl.Logger.Warn("user does not exist validation error")
			return nil, httperror.NewForBadRequestWithSingleField("id", "does not exist")
		}

		if ou.Status == user_s.UserStatusActive {
			ou.Status = user_s.UserStatusArchived
			impl.Logger.Debug("user was archived", slog.String("user_id", ou.ID.Hex()))
		} else {
			ou.Status = user_s.UserStatusActive
			impl.Logger.Debug("user was unarchived", slog.String("user_id", ou.ID.Hex()))
		}
		ou.ModifiedByUserID = userID
		ou.ModifiedAt = time.Now()
		ou.ModifiedByName = userName

		if err := impl.UserStorer.UpdateByID(sessCtx, ou); err != nil {
			impl.Logger.Error("user update by id error", slog.Any("error", err))
			return nil, err
		}
		return ou, nil
	}

	// Start a transaction
	res, err := session.WithTransaction(ctx, transactionFunc)
	if err != nil {
		impl.Logger.Error("session failed error",
			slog.Any("error", err))
		return nil, err
	}

	return res.(*user_s.User), nil
}
