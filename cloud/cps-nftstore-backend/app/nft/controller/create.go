package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	eth "github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/adapter/blockchain/eth"
	pinobj_s "github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/app/ipfsgateway/datastore"
	s_d "github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/app/nft/datastore"
	nftasset_s "github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/app/nftasset/datastore"
	u_d "github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/app/user/datastore"
	"github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/config/constants"
	"github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/utils/cryptowrapper"
	"github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/utils/httperror"
)

type NFTCreateRequestIDO struct {
	CollectionID    primitive.ObjectID          `bson:"collection_id" json:"collection_id"`
	Name            string                      `bson:"name" json:"name"`               // Name of the item.
	Description     string                      `bson:"description" json:"description"` // A human-readable description of the item. Markdown is supported.
	ImageID         primitive.ObjectID          `bson:"image_id" json:"image_id"`
	AnimationID     primitive.ObjectID          `bson:"animation_id" json:"animation_id"`
	ExternalURL     string                      `bson:"external_url" json:"external_url"`         // This is the URL that will appear below the asset's image on OpenSea and will allow users to leave OpenSea and view the item on your site.
	Attributes      []*s_d.NFTMetadataAttribute `bson:"attributes" json:"attributes"`             // These are the attributes for the item, which will show up on the OpenSea page for the item. (see below)
	BackgroundColor string                      `bson:"background_color" json:"background_color"` // Background color of the item on OpenSea. Must be a six-character hexadecimal without a pre-pended #.
	YoutubeURL      string                      `bson:"youtube_url" json:"youtube_url"`           // A URL to a YouTube video (only used if animation_url is not provided).
	ToAddress       string                      `bson:"to_address" json:"to_address"`
	WalletPassword  string                      `bson:"wallet_password" json:"wallet_password"`
}

func (impl *NFTControllerImpl) validateCreateRequest(dirtyData *NFTCreateRequestIDO) error {
	e := make(map[string]string)

	if dirtyData.CollectionID.IsZero() {
		e["collection_id"] = "missing value"
	} else {
		doesExist, err := impl.NFTCollectionStorer.CheckIfExistsByID(context.Background(), dirtyData.CollectionID)
		if err != nil {
			e["collection_id"] = fmt.Sprintf("encountered error: %v", err)
		}
		if !doesExist {
			e["collection_id"] = fmt.Sprintf("does not exist for: %v", dirtyData.CollectionID)
		}
	}
	if dirtyData.Name == "" {
		e["name"] = "missing value"
	}
	if dirtyData.Description == "" {
		e["description"] = "missing value"
	}
	if dirtyData.ImageID.IsZero() {
		e["image_id"] = "missing value"
	} else {
		doesExist, err := impl.NFTAssetStorer.CheckIfExistsByID(context.Background(), dirtyData.ImageID)
		if err != nil {
			e["image_id"] = fmt.Sprintf("encountered error: %v", err)
		}
		if !doesExist {
			e["image_id"] = fmt.Sprintf("does not exist for: %v", dirtyData.ImageID)
		}
	}
	if dirtyData.BackgroundColor == "" {
		e["background_color"] = "missing value"
	}
	if dirtyData.ToAddress == "" {
		e["to_address"] = "missing value"
	}
	if dirtyData.WalletPassword == "" {
		e["wallet_password"] = "missing value"
	}

	if len(e) != 0 {
		return httperror.NewForBadRequest(&e)
	}
	return nil
}

func (impl *NFTControllerImpl) Create(ctx context.Context, req *NFTCreateRequestIDO) (*s_d.NFT, error) {
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

	// Perform our validation and return validation error on any issues detected.
	if err := impl.validateCreateRequest(req); err != nil {
		impl.Logger.Warn("validate nft create request",
			slog.Any("error", err))
		return nil, err
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
		// Begin populating our NFT meta with the user add request data, in
		// addition add tenant-specific information.
		//

		nft := &s_d.NFT{
			CollectionID:          req.CollectionID,
			Name:                  req.Name,
			ImageID:               req.ImageID,
			AnimationID:           req.AnimationID,
			ExternalURL:           req.ExternalURL,
			Description:           req.Description,
			Attributes:            req.Attributes,
			BackgroundColor:       req.BackgroundColor,
			YoutubeURL:            req.YoutubeURL,
			TenantID:              tenantID,
			TenantName:            tenantName,
			TenantTimezone:        tenantTimezone,
			ID:                    primitive.NewObjectID(),
			CreatedAt:             time.Now(),
			CreatedFromIPAddress:  ipAddress,
			ModifiedAt:            time.Now(),
			ModifiedFromIPAddress: ipAddress,
			Status:                s_d.StatusActive,
		}

		//
		// STEP 2
		// Fetch related records from our database.
		//

		// DEVELOPER NOTE: We don't need to check if the fetched records exists
		// b/c we already validated it before entering this transaction. Also
		// animation can be null so we will handle it below.
		collection, err := impl.NFTCollectionStorer.GetByID(sessCtx, req.CollectionID)
		if err != nil {
			impl.Logger.Error("failed get collection by id",
				slog.Any("error", err))
			return nil, err
		}
		imageAsset, err := impl.NFTAssetStorer.GetByID(sessCtx, req.ImageID)
		if err != nil {
			impl.Logger.Error("failed get NFT asset by id",
				slog.Any("error", err))
			return nil, err
		}
		animationAsset, err := impl.NFTAssetStorer.GetByID(sessCtx, req.AnimationID)
		if err != nil {
			impl.Logger.Error("failed get NFT asset by id",
				slog.Any("error", err))
			return nil, err
		}

		//
		// STEP 3
		// Assign to NFT metadata next available `token_id` and then generate
		// new next available token.
		//

		// Update our NFT record.
		nft.TokenID = collection.TokensCount // We are doing this because of our smart-contract.
		nft.CollectionName = collection.Name

		// Update collection to have our new token id.
		collection.TokensCount = collection.TokensCount + 1
		collection.ModifiedAt = time.Now()

		impl.Logger.Debug("nft collection generated new token id",
			slog.String("collection_id", nft.CollectionID.String()),
			slog.Uint64("curr_token_id", nft.TokenID),
			slog.Uint64("new_token_id", collection.TokensCount))

		//
		// STEP 4
		// Decrypt the wallet private key (which is saved in our database in
		// encrypted form) so we can use the plaintext private key for our
		// ethereum deploy smart contract operation.
		//

		plaintextPrivateKey, cryptoErr := cryptowrapper.SymmetricKeyDecryption(collection.WalletEncryptedPrivateKey, req.WalletPassword)
		if cryptoErr != nil {
			impl.Logger.Error("failed to decrypt wallet private key",
				slog.Any("error", cryptoErr))
			return nil, httperror.NewForBadRequestWithSingleField("wallet_password", "incorrect password used")
		}

		impl.Logger.Debug("decrypted ethereum wallet private key",
			slog.String("collection_id", collection.ID.Hex()),
			slog.Uint64("token_id", nft.TokenID))

		//
		// STEP 5
		// Connect to ethereum blockchain network via our node.
		//

		eth := eth.NewAdapter(impl.Logger)
		if connectErr := eth.ConnectToNodeAtURL(collection.NodeURL); connectErr != nil {
			impl.Logger.Error("failed connecting to node",
				slog.Any("error", connectErr))
			return nil, httperror.NewForBadRequestWithSingleField("node_url", fmt.Sprintf("connection error: %v", connectErr))
		}

		//
		// STEP 6
		// Execute the `mint` funciton to our smart contract in the ethereum
		// blockchain network. Afterwords update nft record.
		//

		if mintErr := eth.Mint(collection.SmartContract, plaintextPrivateKey, collection.SmartContractAddress, req.ToAddress); mintErr != nil {
			impl.Logger.Error("failed minting",
				slog.Any("error", mintErr))
			return nil, mintErr
		}

		impl.Logger.Debug("successfully minted",
			slog.String("collection_id", collection.ID.Hex()),
			slog.Uint64("token_id", nft.TokenID),
			slog.String("smart_contract_address", collection.SmartContractAddress))

		// Update our database record.
		nft.MintedToAddress = req.ToAddress
		nft.ModifiedAt = time.Now()
		nft.ModifiedFromIPAddress = ipAddress

		//
		// STEP 7
		// Update our `image` asset to reference our new NFT metadata record and
		// therefore our system will not garbage collect this asset. What is
		// our system doing with garbage collection? Basically if our NFT
		// assets are not pointing to a NFT metadata record within 1 day of
		// creation then our system will garbage collect (i.e. delete) the
		// NFT asset - the following change settings will make sure our
		// system doesn't garbage collect this particular NFT asset.
		//

		imageAsset.Status = nftasset_s.StatusPinning
		imageAsset.NFTMetadataID = nft.ID
		imageAsset.NFTCollectionID = collection.ID
		imageAsset.ModifiedAt = time.Now()
		imageAsset.ModifiedFromIPAddress = ipAddress
		nft.Image = fmt.Sprintf("ipfs://%v", imageAsset.CID)
		nft.ImageFilename = imageAsset.Filename
		nft.ImageCID = imageAsset.CID

		impl.Logger.Debug("image asset set",
			slog.Uint64("token_id", nft.TokenID))

		//
		// STEP 8
		// Same as above step 3 but do this for `animation_url` if the user
		// uploaded an animation with this metadata.
		//

		if animationAsset != nil {
			animationAsset.Status = nftasset_s.StatusPinning
			animationAsset.NFTMetadataID = nft.ID
			animationAsset.NFTCollectionID = collection.ID
			animationAsset.ModifiedAt = time.Now()
			animationAsset.ModifiedFromIPAddress = ipAddress
			nft.AnimationURL = fmt.Sprintf("ipfs://%v", animationAsset.CID)
			nft.AnimationFilename = animationAsset.Filename
			nft.AnimationCID = animationAsset.CID

			impl.Logger.Debug("animation asset set",
				slog.Uint64("token_id", nft.TokenID))
		} else {
			impl.Logger.Debug("animation asset ignored",
				slog.Uint64("token_id", nft.TokenID))
		}

		//
		// STEP 9
		// Generate our metadata file in our collection directory and add to
		// IPFS network. Afterwords keep a record of it.
		//

		impl.Logger.Debug("creating metadata file and adding to ipfs network...",
			slog.Uint64("token_id", nft.TokenID))

		metadataFile := &s_d.NFTMetadataFile{
			Image:           nft.Image,
			ExternalURL:     nft.ExternalURL,
			Description:     nft.Description,
			Name:            nft.Name,
			Attributes:      nft.Attributes,
			BackgroundColor: nft.BackgroundColor,
			AnimationURL:    nft.AnimationURL,
			YoutubeURL:      nft.YoutubeURL,
		}
		metadataFileBin, err := json.Marshal(metadataFile)
		if err != nil {
			impl.Logger.Error("failed marshalling metadata file", slog.Any("error", err))
			return nil, err
		}

		dirCid, metadataFileCID, ipfsUploadErr := impl.IPFS.UploadBytesToDir(sessCtx, metadataFileBin, fmt.Sprintf("%v", nft.TokenID), collection.IPFSDirectoryName)
		if ipfsUploadErr != nil {
			impl.Logger.Error("failed uploading NFT metadata file",
				slog.Any("error", ipfsUploadErr))
			return nil, err
		}

		impl.Logger.Debug("nft metadata file added to ipfs network",
			slog.Uint64("token_id", nft.TokenID),
			slog.String("cid", metadataFileCID))

		//
		// STEP 10
		// Publish to IPNS.
		//

		impl.Logger.Debug("nft collection being republished to ipns...",
			slog.Uint64("token_id", nft.TokenID))

		// Publish the collection directory to IPNS
		resolvedIPNSName, err := impl.IPFS.PublishToIPNS(sessCtx, collection.IPNSKeyName, dirCid)
		if err != nil {
			return nil, fmt.Errorf("failed to publish to IPNS: %v", err)
		}

		// Update the collection record.
		collection.MetadataFileCIDs[0] = metadataFileCID
		collection.IPFSDirectoryCID = dirCid
		collection.IPNSName = resolvedIPNSName
		collection.ModifiedAt = time.Now()
		collection.ModifiedFromIPAddress = ipAddress

		// Update the record.
		nft.FileCID = metadataFileCID
		nft.FileIPNSPath = fmt.Sprintf("/ipns/%v/%v", resolvedIPNSName, nft.TokenID)

		impl.Logger.Debug("nft successuflly republished in ipns",
			slog.Uint64("token_id", nft.TokenID),
			slog.String("ipfs_cid", nft.FileCID),
			slog.String("ipns_path", nft.FileIPNSPath))

		//
		// STEP 11
		// Submit our records for saving in our database.
		//

		// Save the nft data to the database
		if err := impl.NFTStorer.Create(sessCtx, nft); err != nil {
			impl.Logger.Error("failed to create nft to database",
				slog.Any("error", err))
			return nil, err
		}

		if err := impl.NFTCollectionStorer.UpdateByID(sessCtx, collection); err != nil {
			impl.Logger.Error("failed updating collection by id",
				slog.Any("error", err))
			return nil, err
		}

		if err := impl.NFTAssetStorer.UpdateByID(sessCtx, imageAsset); err != nil {
			impl.Logger.Error("failed update NFT asset by id",
				slog.Any("error", err))
			return nil, err
		}

		if animationAsset != nil {
			if err := impl.NFTAssetStorer.UpdateByID(sessCtx, animationAsset); err != nil {
				impl.Logger.Error("failed update NFT asset by id",
					slog.Any("error", err))
				return nil, err
			}
		}

		impl.Logger.Debug("nft metadata created",
			slog.Uint64("token_id", nft.TokenID))

		//
		// STEP 12
		// Keep a record of our pinned object for IPFS gateway.
		//

		pinObject := &pinobj_s.PinObject{
			ID:          primitive.NewObjectID(),
			IPNSPath:    nft.FileIPNSPath,
			CID:         nft.FileCID,
			Content:     nil,
			Filename:    fmt.Sprintf("%v", nft.TokenID), // We set it to this way b/c it is required by our `Smart Contract` to write the names like this - This is not an error!
			ContentType: "application/json",
			CreatedAt:   nft.CreatedAt,
			ModifiedAt:  nft.ModifiedAt,
		}
		if createdErr := impl.PinObjectStorer.Create(sessCtx, pinObject); createdErr != nil {
			impl.Logger.Error("database create error", slog.Any("error", createdErr))
			return nil, err
		}

		return nft, nil
	}

	// Execute the transaction function within a MongoDB session
	result, err := session.WithTransaction(ctx, transactionFunc)
	if err != nil {
		impl.Logger.Error("transaction failed", slog.Any("error", err))
		return nil, err
	}

	return result.(*s_d.NFT), nil
}
