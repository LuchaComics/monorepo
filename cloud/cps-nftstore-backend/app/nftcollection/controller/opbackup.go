package controller

import (
	"context"
	"log/slog"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	nft_s "github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/app/nft/datastore"
	nftasset_s "github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/app/nftasset/datastore"
	collection_s "github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/app/nftcollection/datastore"
	u_d "github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/app/user/datastore"
	"github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/config/constants"
	"github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/utils/cryptowrapper"
	"github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/utils/httperror"
)

type NFTCollectionBackupOperationRequestIDO struct {
	NFTCollectionID primitive.ObjectID `bson:"nft_collection_id" json:"nft_collection_id"`
	WalletPassword  string             `bson:"wallet_password" json:"wallet_password"`
}

func ValidateBackupRequest(dirtyData *NFTCollectionBackupOperationRequestIDO) error {
	e := make(map[string]string)

	if dirtyData.NFTCollectionID.IsZero() {
		e["nft_collection_id"] = "missing value"
	}
	if dirtyData.WalletPassword == "" {
		e["wallet_password"] = "missing value"
	}

	if len(e) != 0 {
		return httperror.NewForBadRequest(&e)
	}
	return nil
}

type NFTCollectionBackupOperationResponseIDO struct {
	NFTCollection *collection_s.NFTCollection `bson:"nft_collection" json:"nft_collection"`
	NFTs          []*nft_s.NFT                `bson:"nfts" json:"nfts"`
	NFTAssets     []*nftasset_s.NFTAsset      `bson:"nft_assets" json:"nft_assets"`
}

func (impl *NFTCollectionControllerImpl) OperationBackup(ctx context.Context, req *NFTCollectionBackupOperationRequestIDO) (*NFTCollectionBackupOperationResponseIDO, error) {
	if err := ValidateBackupRequest(req); err != nil {
		return nil, err
	}

	// Extract user and tenant information from the session context
	userRole, _ := ctx.Value(constants.SessionUserRole).(int8)
	// userID := ctx.Value(constants.SessionUserID).(primitive.ObjectID)
	// userName := ctx.Value(constants.SessionUserName).(string)
	// tenantID, _ := ctx.Value(constants.SessionUserTenantID).(primitive.ObjectID)
	// tenantName, _ := ctx.Value(constants.SessionUserTenantName).(string)
	// tenantTimezone, _ := ctx.Value(constants.SessionUserTenantTimezone).(string)
	// ipAddress, _ := ctx.Value(constants.SessionIPAddress).(string)

	// Check if the user has the necessary permissions
	switch userRole {
	case u_d.UserRoleRoot:
		// Access is granted; proceed with the operation
	default:
		// Deny access if the user does not have the 'Root' role
		return nil, httperror.NewForForbiddenWithSingleField("message", "you do not have permission")
	}

	// Start a MongoDB session for transaction management
	session, startSessErr := impl.DbClient.StartSession()
	if startSessErr != nil {
		impl.Logger.Error("failed to start database session", slog.Any("error", startSessErr))
		return nil, startSessErr
	}
	defer session.EndSession(ctx)

	// Define the transaction function to perform a series of operations atomically
	transactionFunc := func(sessCtx mongo.SessionContext) (interface{}, error) {
		//
		// STEP 1
		// Fetch all the related records.
		//

		collection, getCollectionErr := impl.NFTCollectionStorer.GetByID(sessCtx, req.NFTCollectionID)
		if getCollectionErr != nil {
			impl.Logger.Error("failed to get collection by id", slog.Any("error", getCollectionErr))
			return nil, getCollectionErr
		}
		if collection == nil {
			impl.Logger.Error("collection d.n.e.", slog.Any("nft_collection_id", req.NFTCollectionID))
			return nil, httperror.NewForBadRequestWithSingleField("nft_collection_id", "does not exist")
		}

		nfts, getNFTsErr := impl.NFTStorer.ListByNFTCollectionID(sessCtx, req.NFTCollectionID)
		if getNFTsErr != nil {
			impl.Logger.Error("failed to get nft assets by colleciton id", slog.Any("error", getNFTsErr))
			return nil, getNFTsErr
		}
		nftAssets, getNFTAssetsErr := impl.NFTAssetStorer.ListByNFTCollectionID(sessCtx, req.NFTCollectionID)
		if getNFTAssetsErr != nil {
			impl.Logger.Error("failed to get nft assets by colleciton id", slog.Any("error", getNFTAssetsErr))
			return nil, getNFTAssetsErr
		}

		//
		// STEP 2
		// Decrypt the wallet private key (which is saved in our database in
		// encrypted form) so we can use the plaintext private key for our
		// ethereum deploy smart contract operation.
		//

		_, cryptoErr := cryptowrapper.SymmetricKeyDecryption(collection.WalletEncryptedPrivateKey, req.WalletPassword)
		if cryptoErr != nil {
			impl.Logger.Error("failed to decrypt wallet private key", slog.Any("error", cryptoErr))
			return nil, httperror.NewForBadRequestWithSingleField("wallet_password", "incorrect password used")
		}

		impl.Logger.Debug("decrypted ethereum wallet private key",
			slog.String("collection_id", collection.ID.Hex()))

		res := &NFTCollectionBackupOperationResponseIDO{
			NFTCollection: collection,
			NFTs:          nfts.Results,
			NFTAssets:     nftAssets.Results,
		}

		return res, nil
	}

	// Execute the transaction function within a MongoDB session
	result, err := session.WithTransaction(ctx, transactionFunc)
	if err != nil {
		impl.Logger.Error("transaction failed", slog.Any("error", err))
		return nil, err
	}

	return result.(*NFTCollectionBackupOperationResponseIDO), nil
}
