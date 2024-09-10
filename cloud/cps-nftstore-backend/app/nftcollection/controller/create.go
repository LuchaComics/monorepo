package controller

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	s_d "github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/app/nftcollection/datastore"
	u_d "github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/app/user/datastore"
	"github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/config/constants"
	"github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/utils/httperror"
)

func (impl *NFTCollectionControllerImpl) Create(ctx context.Context, collection *s_d.NFTCollection) (*s_d.NFTCollection, error) {
	// Extract user and tenant information from the session context
	userRole, _ := ctx.Value(constants.SessionUserRole).(int8)
	// userID := ctx.Value(constants.SessionUserID).(primitive.ObjectID)
	// userName := ctx.Value(constants.SessionUserName).(string)
	tenantID, _ := ctx.Value(constants.SessionUserTenantID).(primitive.ObjectID)
	tenantName, _ := ctx.Value(constants.SessionUserTenantName).(string)
	tenantTimezone, _ := ctx.Value(constants.SessionUserTenantTimezone).(string)
	ipAddress, _ := ctx.Value(constants.SessionIPAddress).(string)

	// Check if the user has the necessary permissions
	switch userRole {
	case u_d.UserRoleRoot:
		// Access is granted; proceed with the operation
	default:
		// Deny access if the user does not have the 'Root' role
		return nil, httperror.NewForForbiddenWithSingleField("message", "you do not have permission")
	}

	// Start a MongoDB session for transaction management
	session, err := impl.DbClient.StartSession()
	if err != nil {
		impl.Logger.Error("failed to start database session", slog.Any("error", err))
		return nil, err
	}
	defer session.EndSession(ctx)

	// Define the transaction function to perform a series of operations atomically
	transactionFunc := func(sessCtx mongo.SessionContext) (interface{}, error) {
		// Populate collection with default and tenant-specific information
		collection.TenantID = tenantID
		collection.TenantName = tenantName
		collection.TenantTimezone = tenantTimezone
		collection.ID = primitive.NewObjectID()
		collection.CreatedAt = time.Now()
		collection.CreatedFromIPAddress = ipAddress
		collection.ModifiedAt = time.Now()
		collection.ModifiedFromIPAddress = ipAddress
		collection.Status = s_d.StatusActive

		// Generate a unique IPNS key for the collection
		ipnsKeyName := fmt.Sprintf("ipns_key_%s", collection.ID.Hex())
		ipnsName, err := impl.IPFS.GenerateKey(sessCtx, ipnsKeyName)
		if err != nil {
			impl.Logger.Error("failed to generate IPNS key", slog.Any("error", err))
			return nil, err
		}

		// Set a custom directory name for the collection in IPFS
		collection.IPFSDirectoryName = fmt.Sprintf("%v_metadata", collection.ID.Hex())

		// Store IPNS-related data in the collection
		collection.IPNSKeyName = ipnsKeyName
		collection.IPNSName = ipnsName

		// Create a new directory in IPFS with a sample file named "0" (representing the first token)
		collectionDirCID, firstTokenFileCID, err := impl.IPFS.UploadStringToDir(
			context.Background(),
			"Hello world via `Collectibles Protective Services`!", // Sample content for the file
			"0", // First token ID
			collection.IPFSDirectoryName)
		if err != nil {
			return nil, fmt.Errorf("failed to add content to IPFS: %v", err)
		}

		// Update collection data with IPFS directory CID and metadata
		collection.IPFSDirectoryCID = collectionDirCID
		collection.TokensCount = 0
		collection.MetadataFileCIDs = make(map[uint64]string, 0) // Initialize arrays.
		collection.MetadataFileCIDs[0] = firstTokenFileCID

		impl.Logger.Debug("publishing to ipns...",
			slog.String("key_name", ipnsKeyName),
			slog.String("collection_cid", collectionDirCID),
			slog.String("first_token_cid", firstTokenFileCID))

		// Publish the collection directory to IPNS
		resolvedIPNSName, err := impl.IPFS.PublishToIPNS(sessCtx, ipnsKeyName, collectionDirCID)
		if err != nil {
			return nil, fmt.Errorf("failed to publish to IPNS: %v", err)
		}

		impl.Logger.Debug("finished publishing to ipns",
			slog.String("ipns_name", resolvedIPNSName))

		// Ensure the IPNS name matches the resolved name
		if !strings.Contains(ipnsName, resolvedIPNSName) {
			return nil, fmt.Errorf("IPNS name mismatch: expected %s, got %s", ipnsName, resolvedIPNSName)
		}

		// Save the collection data to the database
		if err := impl.NFTCollectionStorer.Create(sessCtx, collection); err != nil {
			impl.Logger.Error("failed to save collection to database", slog.Any("error", err))
			return nil, err
		}

		return collection, nil
	}

	// Execute the transaction function within a MongoDB session
	result, err := session.WithTransaction(ctx, transactionFunc)
	if err != nil {
		impl.Logger.Error("transaction failed", slog.Any("error", err))
		return nil, err
	}

	return result.(*s_d.NFTCollection), nil
}
