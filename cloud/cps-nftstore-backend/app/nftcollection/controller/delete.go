package controller

import (
	"context"
	"log/slog"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	u_d "github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/app/user/datastore"
	"github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/config/constants"
	"github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/utils/httperror"
)

func (impl *NFTCollectionControllerImpl) DeleteByID(ctx context.Context, id primitive.ObjectID) error {
	// Extract user and tenant information from the session context
	userRole, _ := ctx.Value(constants.SessionUserRole).(int8)

	// Check if the user has the necessary permissions
	switch userRole {
	case u_d.UserRoleRoot:
		// Access is granted; proceed with the operation
	default:
		// Deny access if the user does not have the 'Root' role
		return httperror.NewForForbiddenWithSingleField("message", "you do not have permission")
	}

	// Start a MongoDB session for transaction management
	session, startSessErr := impl.DbClient.StartSession()
	if startSessErr != nil {
		impl.Logger.Error("failed to start database session", slog.Any("error", startSessErr))
		return startSessErr
	}
	defer session.EndSession(ctx)

	// Define the transaction function to perform a series of operations atomically
	transactionFunc := func(sessCtx mongo.SessionContext) (interface{}, error) {
		d, err := impl.GetByID(sessCtx, id)
		if err != nil {
			impl.Logger.Error("database get by id error", slog.Any("error", err))
			return nil, err
		}
		if d == nil {
			impl.Logger.Error("database returns nothing from get by id")
			return nil, httperror.NewForBadRequestWithSingleField("id", "collection does not exist")
		}

		//
		// STEP 1
		// List through all the NFT metadata and delete them from IPFS network and
		// from the database.
		//

		metadataRes, listErr := impl.NFTStorer.ListByNFTCollectionID(sessCtx, id)
		if listErr != nil {
			impl.Logger.Error("database list by nft collection error", slog.Any("error", listErr))
			return nil, listErr
		}
		for _, metadata := range metadataRes.Results {
			if err := impl.IPFS.Unpin(sessCtx, metadata.FileCID); err != nil {
				impl.Logger.Error("ipfs failed unpinning file error", slog.Any("error", err))
				return nil, err
			}
			if deleteErr := impl.NFTStorer.DeleteByID(sessCtx, metadata.ID); deleteErr != nil {
				impl.Logger.Error("database delete by id error", slog.Any("error", deleteErr))
				return nil, deleteErr
			}
		}

		//
		// STEP 2
		// List through all the NFT assets and delete them from IPFS network and
		// from the database.
		//

		assetsRes, listErr := impl.NFTAssetStorer.ListByNFTCollectionID(sessCtx, id)
		if listErr != nil {
			impl.Logger.Error("database list by nft asset error", slog.Any("error", listErr))
			return nil, listErr
		}
		for _, asset := range assetsRes.Results {
			if err := impl.IPFS.Unpin(sessCtx, asset.CID); err != nil {
				impl.Logger.Error("ipfs failed unpinning file error", slog.Any("error", err))
				return nil, err
			}
			if deleteErr := impl.NFTAssetStorer.DeleteByID(sessCtx, asset.ID); deleteErr != nil {
				impl.Logger.Error("database delete by id error", slog.Any("error", deleteErr))
				return nil, deleteErr
			}
		}

		//
		// STEP 3
		// Delete the collection folder from IPFS network.
		//

		if err := impl.IPFS.Unpin(sessCtx, d.IPFSDirectoryCID); err != nil {
			impl.Logger.Error("ipfs failed unpinning directory error", slog.Any("error", err))
			return nil, err
		}
		impl.Logger.Debug("deleted all nft metadata files from ipfs for this collection")

		if err := impl.IPFS.RemoveKey(sessCtx, d.IPNSKeyName); err != nil {
			impl.Logger.Error("ipfs removed key error", slog.Any("error", err))
			return nil, err
		}
		impl.Logger.Debug("removed collection ipns key from the ipfs node")

		//
		// STEP 4
		// Deleted the database record.
		//

		if err := impl.NFTCollectionStorer.DeleteByID(sessCtx, id); err != nil {
			impl.Logger.Error("database update by id error", slog.Any("error", err))
			return nil, err
		}
		return nil, nil
	}

	// Execute the transaction function within a MongoDB session
	_, err := session.WithTransaction(ctx, transactionFunc)
	if err != nil {
		impl.Logger.Error("transaction failed", slog.Any("error", err))
		return err
	}
	return nil
}
