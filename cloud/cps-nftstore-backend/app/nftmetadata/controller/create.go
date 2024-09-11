package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	pinobj_s "github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/app/ipfsgateway/datastore"
	nftasset_s "github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/app/nftasset/datastore"
	s_d "github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/app/nftmetadata/datastore"
	u_d "github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/app/user/datastore"
	"github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/config/constants"
	"github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/utils/httperror"
)

type NFTMetadataCreateRequestIDO struct {
	CollectionID    primitive.ObjectID          `bson:"collection_id" json:"collection_id"`
	Name            string                      `bson:"name" json:"name"`               // Name of the item.
	Description     string                      `bson:"description" json:"description"` // A human-readable description of the item. Markdown is supported.
	ImageID         primitive.ObjectID          `bson:"image_id" json:"image_id"`
	AnimationID     primitive.ObjectID          `bson:"animation_id" json:"animation_id"`
	ExternalURL     string                      `bson:"external_url" json:"external_url"`         // This is the URL that will appear below the asset's image on OpenSea and will allow users to leave OpenSea and view the item on your site.
	Attributes      []*s_d.NFTMetadataAttribute `bson:"attributes" json:"attributes"`             // These are the attributes for the item, which will show up on the OpenSea page for the item. (see below)
	BackgroundColor string                      `bson:"background_color" json:"background_color"` // Background color of the item on OpenSea. Must be a six-character hexadecimal without a pre-pended #.
	YoutubeURL      string                      `bson:"youtube_url" json:"youtube_url"`           // A URL to a YouTube video (only used if animation_url is not provided).
}

func (impl *NFTMetadataControllerImpl) validateCreateRequest(dirtyData *NFTMetadataCreateRequestIDO) error {
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

	if len(e) != 0 {
		return httperror.NewForBadRequest(&e)
	}
	return nil
}

func (impl *NFTMetadataControllerImpl) Create(ctx context.Context, req *NFTMetadataCreateRequestIDO) (*s_d.NFTMetadata, error) {
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

		nftmetadata := &s_d.NFTMetadata{
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
		// Fetch related records.
		//

		// DEVELOPER NOTE: We don't need to check if the fetched records exists
		// b/c we already validated it before entering this transaction. Also
		// animation can be null so we will handle it below.
		collection, err := impl.NFTCollectionStorer.GetByID(sessCtx, req.CollectionID)
		if err != nil {
			impl.Logger.Error("failed get collection by id", slog.Any("error", err))
			return nil, err
		}
		imageAsset, err := impl.NFTAssetStorer.GetByID(sessCtx, req.ImageID)
		if err != nil {
			impl.Logger.Error("failed get NFT asset by id", slog.Any("error", err))
			return nil, err
		}
		animationAsset, err := impl.NFTAssetStorer.GetByID(sessCtx, req.AnimationID)
		if err != nil {
			impl.Logger.Error("failed get NFT asset by id", slog.Any("error", err))
			return nil, err
		}

		//
		// STEP 3
		// Assign to NFT metadata next available `token_id` and then generate
		// new next available token.
		//

		// Update our metadata record.
		nftmetadata.TokenID = collection.TokensCount // We are doing this because of our smart-contract.
		nftmetadata.CollectionName = collection.Name

		// Update collection to have our new token id.
		collection.TokensCount = collection.TokensCount + 1
		collection.ModifiedAt = time.Now()

		impl.Logger.Debug("nft collection generated new token id",
			slog.Uint64("metadata_token_id", nftmetadata.TokenID),
			slog.Uint64("new_token_id", collection.TokensCount))

		//
		// STEP 3
		// Update our `image` asset to reference our new NFT metadata record and
		// therefore our system will not garbage collect this asset. What is
		// our system doing with garbage collection? Basically if our NFT
		// assets are not pointing to a NFT metadata record within 1 day of
		// creation then our system will garbage collect (i.e. delete) the
		// NFT asset - the following change settings will make sure our
		// system doesn't garbage collect this particular NFT asset.
		//

		imageAsset.Status = nftasset_s.StatusPinning
		imageAsset.NFTMetadataID = nftmetadata.ID
		imageAsset.NFTCollectionID = collection.ID
		imageAsset.ModifiedAt = time.Now()
		imageAsset.ModifiedFromIPAddress = ipAddress
		nftmetadata.Image = fmt.Sprintf("ipfs://%v", imageAsset.CID)
		nftmetadata.ImageFilename = imageAsset.Filename
		nftmetadata.ImageCID = imageAsset.CID

		impl.Logger.Debug("image asset set")

		//
		// STEP 4
		// Same as above step 3 but do this for `animation_url` if the user
		// uploaded an animation with this metadata.
		//

		if animationAsset != nil {
			animationAsset.Status = nftasset_s.StatusPinning
			animationAsset.NFTMetadataID = nftmetadata.ID
			animationAsset.NFTCollectionID = collection.ID
			animationAsset.ModifiedAt = time.Now()
			animationAsset.ModifiedFromIPAddress = ipAddress
			nftmetadata.AnimationURL = fmt.Sprintf("ipfs://%v", animationAsset.CID)
			nftmetadata.AnimationFilename = animationAsset.Filename
			nftmetadata.AnimationCID = animationAsset.CID

			impl.Logger.Debug("animation asset set")
		} else {
			impl.Logger.Debug("animation asset ignored")
		}

		//
		// STEP 5
		// Generate our metadata file in our collection directory and add to
		// IPFS network. Afterwords keep a record of it.
		//

		impl.Logger.Debug("creating metadata file and adding to ipfs network...",
			slog.Uint64("token_id", nftmetadata.TokenID))

		metadataFile := &s_d.NFTMetadataFile{
			Image:           nftmetadata.Image,
			ExternalURL:     nftmetadata.ExternalURL,
			Description:     nftmetadata.Description,
			Name:            nftmetadata.Name,
			Attributes:      nftmetadata.Attributes,
			BackgroundColor: nftmetadata.BackgroundColor,
			AnimationURL:    nftmetadata.AnimationURL,
			YoutubeURL:      nftmetadata.YoutubeURL,
		}
		metadataFileBin, err := json.Marshal(metadataFile)
		if err != nil {
			impl.Logger.Error("failed marshalling metadata file", slog.Any("error", err))
			return nil, err
		}

		dirCid, metadataFileCID, ipfsUploadErr := impl.IPFS.UploadBytesToDir(sessCtx, metadataFileBin, fmt.Sprintf("%v", nftmetadata.TokenID), collection.IPFSDirectoryName)
		if ipfsUploadErr != nil {
			impl.Logger.Error("failed uploading NFT metadata file",
				slog.Any("error", ipfsUploadErr))
			return nil, err
		}

		impl.Logger.Debug("nft metadata file added to ipfs network",
			slog.Uint64("token_id", nftmetadata.TokenID),
			slog.String("cid", metadataFileCID))

		//
		// STEP 6
		// Publish to IPNS.
		//

		impl.Logger.Debug("nft collection being republished to ipns...",
			slog.Uint64("token_id", nftmetadata.TokenID))

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
		nftmetadata.FileCID = metadataFileCID
		nftmetadata.FileIPNSPath = fmt.Sprintf("/ipns/%v/%v", resolvedIPNSName, nftmetadata.TokenID)

		impl.Logger.Debug("nft successuflly republished in ipns",
			slog.Uint64("token_id", nftmetadata.TokenID),
			slog.String("ipfs_cid", nftmetadata.FileCID),
			slog.String("ipns_path", nftmetadata.FileIPNSPath))

		//
		// STEP 7
		// Submit our records for saving in our database.
		//

		// Save the nftmetadata data to the database
		if err := impl.NFTMetadataStorer.Create(sessCtx, nftmetadata); err != nil {
			impl.Logger.Error("failed to create nftmetadata to database",
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
			slog.Uint64("token_id", nftmetadata.TokenID))

		//
		// STEP 8
		// Keep a record of our pinned object for IPFS gateway.
		//

		pinObject := &pinobj_s.PinObject{
			ID:          primitive.NewObjectID(),
			IPNSPath:    nftmetadata.FileIPNSPath,
			CID:         nftmetadata.FileCID,
			Content:     nil,
			Filename:    fmt.Sprintf("%v", nftmetadata.TokenID), // We set it to this way b/c it is required by our `Smart Contract` to write the names like this - This is not an error!
			ContentType: "application/json",
			CreatedAt:   nftmetadata.CreatedAt,
			ModifiedAt:  nftmetadata.ModifiedAt,
		}
		if createdErr := impl.PinObjectStorer.Create(sessCtx, pinObject); createdErr != nil {
			impl.Logger.Error("database create error", slog.Any("error", createdErr))
			return nil, err
		}

		return nftmetadata, nil
	}

	// Execute the transaction function within a MongoDB session
	result, err := session.WithTransaction(ctx, transactionFunc)
	if err != nil {
		impl.Logger.Error("transaction failed", slog.Any("error", err))
		return nil, err
	}

	return result.(*s_d.NFTMetadata), nil
}
