package controller

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/adapter/blockchain/eth"
	s_d "github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/app/nftcollection/datastore"
	u_d "github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/app/user/datastore"
	"github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/config/constants"
	"github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/utils/cryptowrapper"
	"github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/utils/httperror"
)

type NFTCollectionCreateRequestIDO struct {
	Blockchain     string             `bson:"blockchain" json:"blockchain"`
	NodeURL        string             `bson:"node_url" json:"node_url"`
	SmartContract  string             `bson:"smart_contract" json:"smart_contract"`
	WalletMnemonic string             `bson:"wallet_mnemonic" json:"wallet_mnemonic"`
	WalletPassword string             `bson:"wallet_password" json:"wallet_password"`
	TenantID       primitive.ObjectID `bson:"tenant_id" json:"tenant_id"`
	Name           string             `bson:"name" json:"name"`
}

func ValidateCreateRequest(dirtyData *NFTCollectionCreateRequestIDO) error {
	e := make(map[string]string)

	if dirtyData.Blockchain == "" {
		e["blockchain"] = "missing value"
	}
	if dirtyData.NodeURL == "" {
		e["node_url"] = "missing value"
	}
	if dirtyData.SmartContract == "" {
		e["smart_contract"] = "missing value"
	}
	if dirtyData.WalletMnemonic == "" {
		e["wallet_mnemonic"] = "missing value"
	}
	if dirtyData.WalletPassword == "" {
		e["wallet_password"] = "missing value"
	}
	if dirtyData.TenantID.IsZero() {
		e["tenant_id"] = "missing value"
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

		tenant, getTenentErr := impl.TenantStorer.GetByID(sessCtx, req.TenantID)
		if getTenentErr != nil {
			impl.Logger.Error("failed to get tenant by id", slog.Any("error", getTenentErr))
			return nil, getTenentErr
		}
		if tenant == nil {
			impl.Logger.Error("tenant d.n.e.", slog.Any("tenant_id", req.TenantID))
			return nil, httperror.NewForBadRequestWithSingleField("tenant_id", "does not exist")
		}

		//
		// STEP 2
		// Populate collection with default and tenant-specific information
		// along with the user create request data.
		//

		collection := &s_d.NFTCollection{
			ID:                        primitive.NewObjectID(),
			TenantID:                  tenantID,
			TenantName:                tenantName,
			TenantTimezone:            tenantTimezone,
			CreatedAt:                 time.Now(),
			CreatedFromIPAddress:      ipAddress,
			ModifiedAt:                time.Now(),
			ModifiedFromIPAddress:     ipAddress,
			Status:                    s_d.StatusActive,
			Name:                      req.Name,
			Blockchain:                req.Blockchain,
			SmartContract:             req.SmartContract,
			NodeURL:                   req.NodeURL,
			WalletAccountAddress:      "",
			WalletEncryptedPrivateKey: "",
			WalletPublicKey:           "",
			SmartContractStatus:       s_d.SmartContractStatusNotDeployed,
			SmartContractAddress:      "",
		}

		impl.Logger.Debug("creating nft collection...",
			slog.String("collection_id", collection.ID.Hex()))

		//
		// STEP 3
		// Generate new wallet for this NFT collection.
		//

		eth := eth.NewAdapter(impl.Logger)
		if connectErr := eth.ConnectToNodeAtURL(req.NodeURL); connectErr != nil {
			impl.Logger.Error("failed connecting to node", slog.Any("error", connectErr))
			return nil, httperror.NewForBadRequestWithSingleField("node_url", fmt.Sprintf("connection error: %v", connectErr))
		}
		wallet, walletErr := eth.NewWalletFromMnemonic(req.WalletMnemonic)
		if walletErr != nil {
			impl.Logger.Error("failed to generate ethereum wallet", slog.Any("error", walletErr))
			return nil, walletErr
		}
		collection.WalletAccountAddress = wallet.AccountAddress
		collection.WalletPublicKey = wallet.PublicKey

		impl.Logger.Debug("generated ethereum wallet",
			slog.String("collection_id", collection.ID.Hex()),
			slog.String("public_key", collection.WalletPublicKey),
			slog.String("account_address", collection.WalletAccountAddress))

		//
		// STEP 4
		// Encrypt our wallet private key with the user's password so it is safely stored in our database in encrypted form.
		//

		plaintextPrivateKey := wallet.PrivateKey
		encryptedPrivateKey, cryptoErr := cryptowrapper.SymmetricKeyEncryption(plaintextPrivateKey, req.WalletPassword)
		if cryptoErr != nil {
			impl.Logger.Error("failed to encrypt wallet private key", slog.Any("error", cryptoErr))
			return nil, cryptoErr
		}

		collection.WalletEncryptedPrivateKey = encryptedPrivateKey

		impl.Logger.Debug("encrypted ethereum wallet private key",
			slog.String("collection_id", collection.ID.Hex()))

		//
		// STEP 5
		// Generate unique key for this collection to utilize in IPNS.
		//

		// Generate a unique IPNS key for the collection
		ipnsKeyName := fmt.Sprintf("ipns_key_%s", collection.ID.Hex())
		ipnsName, keyGenErr := impl.IPFS.GenerateKey(sessCtx, ipnsKeyName)
		if keyGenErr != nil {
			impl.Logger.Error("failed to generate IPNS key", slog.Any("error", keyGenErr))
			return nil, keyGenErr
		}

		// Set a custom directory name for the collection in IPFS
		collection.IPFSDirectoryName = fmt.Sprintf("%v_metadata", collection.ID.Hex())

		// Store IPNS-related data in the collection
		collection.IPNSKeyName = ipnsKeyName
		collection.IPNSName = ipnsName

		impl.Logger.Debug("generated ipns key",
			slog.String("collection_id", collection.ID.Hex()),
			slog.String("ipns_key_name", ipnsKeyName))

		//
		// STEP 6
		//

		// Create a new directory in IPFS with a sample file named "0" (representing the first token)
		collectionDirCID, firstTokenFileCID, uploadErr := impl.IPFS.UploadStringToDir(
			context.Background(),
			"Hello world via `Collectibles Protective Services`!", // Sample content for the file
			"0", // First token ID
			collection.IPFSDirectoryName)
		if uploadErr != nil {
			return nil, fmt.Errorf("failed to add content to IPFS: %v", uploadErr)
		}

		// Update collection data with IPFS directory CID and metadata
		collection.IPFSDirectoryCID = collectionDirCID
		collection.TokensCount = 0
		collection.MetadataFileCIDs = make(map[uint64]string, 0) // Initialize arrays.
		collection.MetadataFileCIDs[0] = firstTokenFileCID

		impl.Logger.Debug("publishing to ipns...",
			slog.String("collection_id", collection.ID.Hex()),
			slog.String("ipns_key_name", ipnsKeyName),
			slog.String("first_token_cid", firstTokenFileCID))

		//
		// STEP 7
		//

		// Publish the collection directory to IPNS
		resolvedIPNSName, publishIPNSErr := impl.IPFS.PublishToIPNS(sessCtx, ipnsKeyName, collectionDirCID)
		if publishIPNSErr != nil {
			return nil, fmt.Errorf("failed to publish to IPNS: %v", publishIPNSErr)
		}

		impl.Logger.Debug("finished publishing to ipns",
			slog.String("collection_id", collection.ID.Hex()),
			slog.String("ipns_name", resolvedIPNSName))

		// Ensure the IPNS name matches the resolved name
		if !strings.Contains(ipnsName, resolvedIPNSName) {
			return nil, fmt.Errorf("IPNS name mismatch: expected %s, got %s", ipnsName, resolvedIPNSName)
		}

		//
		// STEP 8
		//

		// Save the collection data to the database
		if createErr := impl.NFTCollectionStorer.Create(sessCtx, collection); createErr != nil {
			impl.Logger.Error("failed to save collection to database", slog.Any("error", createErr))
			return nil, createErr
		}

		impl.Logger.Debug("finished creating nft collection",
			slog.String("collection_id", collection.ID.Hex()))

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
