package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	domain "github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/app/nft/datastore"
	s_d "github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/app/nft/datastore"
	u_d "github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/app/user/datastore"
	"github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/config/constants"
	"github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/utils/httperror"
)

type NFTUpdateRequestIDO struct {
	ID          primitive.ObjectID `bson:"id" json:"id"`
	Name        string             `bson:"name" json:"name"`               // Name of the item.
	Description string             `bson:"description" json:"description"` // A human-readable description of the item. Markdown is supported.
	//ImageID         primitive.ObjectID          `bson:"image_id" json:"image_id"`
	//AnimationID     primitive.ObjectID          `bson:"animation_id" json:"animation_id"`
	ExternalURL     string                      `bson:"external_url" json:"external_url"`         // This is the URL that will appear below the asset's image on OpenSea and will allow users to leave OpenSea and view the item on your site.
	Attributes      []*s_d.NFTMetadataAttribute `bson:"attributes" json:"attributes"`             // These are the attributes for the item, which will show up on the OpenSea page for the item. (see below)
	BackgroundColor string                      `bson:"background_color" json:"background_color"` // Background color of the item on OpenSea. Must be a six-character hexadecimal without a pre-pended #.
	YoutubeURL      string                      `bson:"youtube_url" json:"youtube_url"`           // A URL to a YouTube video (only used if animation_url is not provided).
}

func (impl *NFTControllerImpl) validateUpdateRequest(dirtyData *NFTUpdateRequestIDO) error {
	e := make(map[string]string)
	if dirtyData.ID.IsZero() {
		e["id"] = "missing value"
	}
	// if dirtyData.CollectionID.IsZero() {
	// 	e["collection_id"] = "missing value"
	// } else {
	// 	doesExist, err := impl.NFTCollectionStorer.CheckIfExistsByID(context.Background(), dirtyData.CollectionID)
	// 	if err != nil {
	// 		e["collection_id"] = fmt.Sprintf("encountered error: %v", err)
	// 	}
	// 	if !doesExist {
	// 		e["collection_id"] = fmt.Sprintf("does not exist for: %v", dirtyData.CollectionID)
	// 	}
	// }
	if dirtyData.Name == "" {
		e["name"] = "missing value"
	}
	if dirtyData.Description == "" {
		e["description"] = "missing value"
	}
	// if dirtyData.ImageID.IsZero() {
	// 	e["image_id"] = "missing value"
	// } else {
	// 	doesExist, err := impl.NFTAssetStorer.CheckIfExistsByID(context.Background(), dirtyData.ImageID)
	// 	if err != nil {
	// 		e["image_id"] = fmt.Sprintf("encountered error: %v", err)
	// 	}
	// 	if !doesExist {
	// 		e["image_id"] = fmt.Sprintf("does not exist for: %v", dirtyData.ImageID)
	// 	}
	// }
	// if dirtyData.AnimationID.IsZero() {
	// 	e["animation_id"] = "missing value"
	// } else {
	// 	doesExist, err := impl.NFTAssetStorer.CheckIfExistsByID(context.Background(), dirtyData.AnimationID)
	// 	if err != nil {
	// 		e["animation_id"] = fmt.Sprintf("encountered error: %v", err)
	// 	}
	// 	if !doesExist {
	// 		e["animation_id"] = fmt.Sprintf("does not exist for: %v", dirtyData.AnimationID)
	// 	}
	// }
	if dirtyData.BackgroundColor == "" {
		e["background_color"] = "missing value"
	}
	if dirtyData.YoutubeURL == "" {
		e["youtube_url"] = "missing value"
	}

	if len(e) != 0 {
		return httperror.NewForBadRequest(&e)
	}
	return nil
}

func (impl *NFTControllerImpl) UpdateByID(ctx context.Context, req *NFTUpdateRequestIDO) (*domain.NFT, error) {
	// Extract from our session the following data.
	urole, _ := ctx.Value(constants.SessionUserRole).(int8)
	// uid := ctx.Value(constants.SessionUserID).(primitive.ObjectID)
	// uname := ctx.Value(constants.SessionUserName).(string)
	oid, _ := ctx.Value(constants.SessionUserTenantID).(primitive.ObjectID)
	oname, _ := ctx.Value(constants.SessionUserTenantName).(string)
	otz, _ := ctx.Value(constants.SessionUserTenantTimezone).(string)
	ipAddress, _ := ctx.Value(constants.SessionIPAddress).(string)

	switch urole { // Security.
	case u_d.UserRoleRoot:
		impl.Logger.Debug("access granted")
	default:
		return nil, httperror.NewForForbiddenWithSingleField("message", "you do not have permission")
	}

	// Perform our validation and return validation error on any issues detected.
	if err := impl.validateUpdateRequest(req); err != nil {
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
		// Fetch all the related records.
		//

		// Fetch the original tenant.
		os, err := impl.NFTStorer.GetByID(sessCtx, req.ID)
		if err != nil {
			impl.Logger.Error("database get by id error", slog.Any("error", err))
			return nil, err
		}
		if os == nil {
			return nil, httperror.NewForBadRequestWithSingleField("id", "nft type does not exist")
		}
		collection, err := impl.NFTCollectionStorer.GetByID(sessCtx, os.CollectionID)
		if err != nil {
			impl.Logger.Error("failed get collection by id", slog.Any("error", err))
			return nil, err
		}
		if collection == nil {
			return nil, httperror.NewForBadRequestWithSingleField("id", "nftcollection type does not exist")
		}

		//
		// STEP 2
		// Modify our record.
		//

		// Modify our record.
		os.ModifiedAt = time.Now()
		os.ModifiedFromIPAddress = ipAddress
		os.TenantID = oid
		os.TenantName = oname
		os.TenantTimezone = otz
		// os.Status = ns.Status
		os.Name = req.Name
		os.Description = req.Description
		os.ExternalURL = req.ExternalURL
		os.Attributes = req.Attributes
		os.BackgroundColor = req.BackgroundColor
		os.YoutubeURL = req.YoutubeURL

		//
		// STEP 3
		// Unpin our previous metadata file currently in the network/
		//

		if ipfsUnpinErr := impl.IPFS.Unpin(sessCtx, os.FileCID); ipfsUnpinErr != nil {
			impl.Logger.Error("failed unpinning NFT metadata file",
				slog.Any("error", ipfsUnpinErr))
			// Skip any errors that occur during the IPFS unpinning process.
		}

		impl.Logger.Debug("unpinned old nft metadata file in ipfs network",
			slog.Uint64("token_id", os.TokenID),
			slog.String("cid", os.FileCID))

		//
		// STEP 4
		// Regenerate our metadata file in our collection directory and add to
		// IPFS network.
		//

		impl.Logger.Debug("recreating metadata file and adding to ipfs network...",
			slog.Uint64("token_id", os.TokenID))

		metadataFile := &s_d.NFTMetadataFile{
			Image:           os.Image,
			ExternalURL:     os.ExternalURL,
			Description:     os.Description,
			Name:            os.Name,
			Attributes:      os.Attributes,
			BackgroundColor: os.BackgroundColor,
			AnimationURL:    os.AnimationURL,
			YoutubeURL:      os.YoutubeURL,
		}
		metadataFileBin, err := json.Marshal(metadataFile)
		if err != nil {
			impl.Logger.Error("failed marshalling metadata file", slog.Any("error", err))
			return nil, err
		}

		dirCid, metadataFileCID, ipfsUploadErr := impl.IPFS.UploadBytes(sessCtx, metadataFileBin, fmt.Sprintf("%v", os.TokenID), collection.IPFSDirectoryName)
		if ipfsUploadErr != nil {
			impl.Logger.Error("failed uploading NFT metadata file",
				slog.Any("error", ipfsUploadErr))
			return nil, err
		}

		impl.Logger.Debug("nft metadata file re-added to ipfs network",
			slog.Uint64("token_id", os.TokenID),
			slog.String("cid", metadataFileCID))

		//
		// STEP 5
		// Publish to IPNS.
		//

		impl.Logger.Debug("nft collection being republished to ipns...",
			slog.Uint64("token_id", os.TokenID))

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
		os.FileCID = metadataFileCID
		os.FileIPNSPath = fmt.Sprintf("/ipns/%v/%v", resolvedIPNSName, os.TokenID)

		impl.Logger.Debug("nft successuflly republished in ipns",
			slog.Uint64("token_id", os.TokenID),
			slog.String("ipfs_cid", os.FileCID),
			slog.String("ipns_path", os.FileIPNSPath))

		//
		// STEP 6
		// Submit our records for saving in our database.
		//

		// Save to the database the modified tenant.
		if err := impl.NFTStorer.UpdateByID(sessCtx, os); err != nil {
			impl.Logger.Error("database update by id error", slog.Any("error", err))
			return nil, err
		}

		if err := impl.NFTCollectionStorer.UpdateByID(sessCtx, collection); err != nil {
			impl.Logger.Error("failed updating collection by id",
				slog.Any("error", err))
			return nil, err
		}

		impl.Logger.Debug("nft metadata updated",
			slog.Uint64("token_id", os.TokenID))

		return os, nil
	}

	// Execute the transaction function within a MongoDB session
	result, err := session.WithTransaction(ctx, transactionFunc)
	if err != nil {
		impl.Logger.Error("transaction failed", slog.Any("error", err))
		return nil, err
	}

	return result.(*s_d.NFT), nil
}
