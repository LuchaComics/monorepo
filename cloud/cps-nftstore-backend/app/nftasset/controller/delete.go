package controller

import (
	"context"

	"log/slog"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	user_d "github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/app/user/datastore"
	"github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/config/constants"
	"github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/utils/httperror"
)

func (impl *NFTAssetControllerImpl) DeleteByID(ctx context.Context, id primitive.ObjectID) error {
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
				slog.Any("id", id),
				slog.Any("role", userRole),
				slog.Any("userID", userID))
			return nil, httperror.NewForForbiddenWithSingleField("message", "you role does not grant you access to this")
		}

		// Update the database.
		nftasset, err := impl.GetByID(sessCtx, id)
		if err != nil {
			impl.Logger.Error("database get by id error",
				slog.Any("id", id),
				slog.Any("error", err))
			return nil, err
		}
		if nftasset == nil {
			impl.Logger.Error("database returns nothing from get by id")
			return nil, err
		}

		// Proceed to delete the physical files from IPFS.
		if err := impl.IPFS.Unpin(sessCtx, nftasset.CID); err != nil {
			impl.Logger.Warn("ipfs delete by cid error",
				slog.Any("id", id),
				slog.String("cid", nftasset.CID),
				slog.Any("error", err))
			// Do not return an error, simply continue this function as there might
			// be a case were the file was removed on the s3 bucket by ourselves
			// or some other reason.
		} else {
			impl.Logger.Debug("nft asset deleted from ipfs")
		}

		if err := impl.NFTAssetStorer.DeleteByID(sessCtx, nftasset.ID); err != nil {
			impl.Logger.Error("database delete by id error",
				slog.Any("id", id),
				slog.Any("error", err))
			return nil, err
		}
		return nil, nil
	}

	// Start a transaction
	if _, err := session.WithTransaction(ctx, transactionFunc); err != nil {
		impl.Logger.Error("session failed error",
			slog.Any("id", id),
			slog.Any("error", err))
		return err
	}

	return nil
}
