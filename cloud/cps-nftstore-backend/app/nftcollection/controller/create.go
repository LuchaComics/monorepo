package controller

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	eth "github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/adapter/blockchain/eth"
	s_d "github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/app/nftcollection/datastore"
	u_d "github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/app/user/datastore"
	"github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/config/constants"
	"github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/utils/httperror"
)

type NFTCollectionCreateRequestIDO struct {
	Blockchain     string             `bson:"blockchain" json:"blockchain"`
	NodeURL        string             `bson:"node_url" json:"node_url"`
	SmartContract  string             `bson:"smart_contract" json:"smart_contract"`
	WalletMnemonic string             `bson:"wallet_mnemonic" json:"wallet_mnemonic"`
	TenantID       primitive.ObjectID `bson:"tenant_id" json:"tenant_id"`
	Name           string             `bson:"name" json:"name"`
}

func ValidateCreateRequest(dirtyData *NFTCollectionCreateRequestIDO) error {
	e := make(map[string]string)

	if dirtyData.TenantID.IsZero() {
		e["tenant_id"] = "missing value"
	}
	if dirtyData.NodeURL == "" {
		e["node_url"] = "missing value"
	}
	if dirtyData.Name == "" {
		e["name"] = "missing value"
	}

	if len(e) != 0 {
		return httperror.NewForBadRequest(&e)
	}
	return nil
}

func (impl *NFTCollectionControllerImpl) Create(ctx context.Context, req *NFTCollectionCreateRequestIDO) (*s_d.NFTCollection, error) {
	if err := ValidateCreateRequest(req); err != nil {
		return nil, err
	}

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
		//
		// STEP 1
		// Fetch all the related records.
		//

		tenant, err := impl.TenantStorer.GetByID(sessCtx, req.TenantID)
		if err != nil {
			impl.Logger.Error("failed to get tenant by id", slog.Any("error", err))
			return nil, err
		}
		if tenant != nil {
			//TODO
		}

		//
		// STEP 2
		// Populate collection with default and tenant-specific information
		// along with the user create request data.
		//

		collection := &s_d.NFTCollection{
			ID:                    primitive.NewObjectID(),
			TenantID:              tenantID,
			TenantName:            tenantName,
			TenantTimezone:        tenantTimezone,
			CreatedAt:             time.Now(),
			CreatedFromIPAddress:  ipAddress,
			ModifiedAt:            time.Now(),
			ModifiedFromIPAddress: ipAddress,
			Status:                s_d.StatusActive,
			// Blockchain:            req.Blockchain, //TODO
			// SmartContract:         req.SmartContract, //TODO
			// NodeURL:               req.NodeURL, //TODO
		}

		//
		// STEP 3
		// Generate new wallet for this NFT collection.
		//

		eth := eth.NewAdapter(impl.Logger, req.NodeURL) // https://github.com/miguelmota/go-ethereum-hdwallet/blob/master/example/keys.go | https://goethereumbook.org/client-setup/
		if eth != nil {
			//TODO
		}
		wallet, err := eth.NewWalletFromMnemonic(req.WalletMnemonic)
		if err != nil {
			impl.Logger.Error("failed to generate ethereum wallet", slog.Any("error", err))
			return nil, err
		}
		if wallet != nil {
			//TODO
		}
		// collection.AccountAddress = wallet.AccountAddress
		// collection.EncryptedPrivateKey = wallet.PrivateKey
		// collection.PublicKey = wallet.PublicKey

		//
		// STEP 4
		// Encrypt our wallet private key with the user's password so it is safely stored in our database in encrypted form.
		//

		//TODO

		//
		// STEP 5
		// Generate unique key for this collection to utilize in IPNS.
		//

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

		//
		// STEP 6
		//

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

		//
		// STEP 7
		//

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

		//
		// STEP 8
		//

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
