package controller

import (
	"context"

	"log/slog"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	org_d "github.com/LuchaComics/monorepo/cloud/cps-pinws-backend/app/attachment/datastore"
	user_d "github.com/LuchaComics/monorepo/cloud/cps-pinws-backend/app/user/datastore"
	"github.com/LuchaComics/monorepo/cloud/cps-pinws-backend/config/constants"
	"github.com/LuchaComics/monorepo/cloud/cps-pinws-backend/utils/httperror"
)

func (impl *AttachmentControllerImpl) DeleteByID(ctx context.Context, id primitive.ObjectID) error {
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
		userID := sessCtx.Value(constants.SessionUserID).(primitive.ObjectID)
		userRole := sessCtx.Value(constants.SessionUserRole).(int8)

		// Apply protection based on ownership and role.
		if userRole != user_d.UserRoleRoot && userRole != user_d.UserRoleRetailer {
			impl.Logger.Error("authenticated user is not staff role error",
				slog.Any("role", userRole),
				slog.Any("userID", userID))
			return nil, httperror.NewForForbiddenWithSingleField("message", "you role does not grant you access to this")
		}

		// Update the database.
		attachment, err := impl.GetByID(sessCtx, id)
		attachment.Status = org_d.StatusArchived
		if err != nil {
			impl.Logger.Error("database get by id error", slog.Any("error", err))
			return nil, err
		}
		if attachment == nil {
			impl.Logger.Error("database returns nothing from get by id")
			return nil, err
		}
		// // Security: Prevent deletion of root user(s).
		// if attachment.Type == org_d.RootType {
		// 	impl.Logger.Warn("root attachment cannot be deleted error")
		// 	return httperror.NewForForbiddenWithSingleField("role", "root attachment cannot be deleted")
		// }

		// Save to the database the modified attachment.
		if err := impl.AttachmentStorer.UpdateByID(sessCtx, attachment); err != nil {
			impl.Logger.Error("database update by id error", slog.Any("error", err))
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

func (impl *AttachmentControllerImpl) PermanentlyDeleteByID(ctx context.Context, id primitive.ObjectID) error {
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
		userID := sessCtx.Value(constants.SessionUserID).(primitive.ObjectID)
		userRole := sessCtx.Value(constants.SessionUserRole).(int8)

		// Apply protection based on ownership and role.
		if userRole != user_d.UserRoleRoot && userRole != user_d.UserRoleRetailer {
			impl.Logger.Error("authenticated user is not staff role error",
				slog.Any("role", userRole),
				slog.Any("userID", userID))
			return nil, httperror.NewForForbiddenWithSingleField("message", "you role does not grant you access to this")
		}

		// Update the database.
		attachment, err := impl.GetByID(sessCtx, id)
		if err != nil {
			impl.Logger.Error("database get by id error", slog.Any("error", err))
			return nil, err
		}
		if attachment == nil {
			impl.Logger.Error("database returns nothing from get by id")
			return nil, err
		}

		// Proceed to delete the physical files from AWS s3.
		if err := impl.S3.DeleteByKeys(sessCtx, []string{attachment.ObjectKey}); err != nil {
			impl.Logger.Warn("s3 delete by keys error", slog.Any("error", err))
			// Do not return an error, simply continue this function as there might
			// be a case were the file was removed on the s3 bucket by ourselves
			// or some other reason.
		}
		// Proceed to delete the physical files from IPFS.
		if err := impl.IPFS.DeleteContent(sessCtx, attachment.CID); err != nil {
			impl.Logger.Warn("ipfs delete by CID error", slog.Any("error", err))
			// Do not return an error, simply continue this function as there might
			// be a case were the file was removed on the s3 bucket by ourselves
			// or some other reason.
		}

		if err := impl.AttachmentStorer.DeleteByID(sessCtx, attachment.ID); err != nil {
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
