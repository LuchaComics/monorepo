package controller

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	nftasset_s "github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/app/nftasset/datastore"
	s_d "github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/app/nftmetadata/datastore"
	u_d "github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/app/user/datastore"
	"github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/config/constants"
	"github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/utils/httperror"
)

type NFTMetadataCreateRequestIDO struct {
	CollectionID    primitive.ObjectID          `bson:"collection_id" json:"collection_id"`
	Name            string                      `bson:"name" json:"name"` // Name of the item.
	ImageID         primitive.ObjectID          `bson:"image_id" json:"image_id"`
	ExternalURL     string                      `bson:"external_url" json:"external_url"`         // This is the URL that will appear below the asset's image on OpenSea and will allow users to leave OpenSea and view the item on your site.
	Description     string                      `bson:"description" json:"description"`           // A human-readable description of the item. Markdown is supported.
	Attributes      []*s_d.NFTMetadataAttribute `bson:"attributes" json:"attributes"`             // These are the attributes for the item, which will show up on the OpenSea page for the item. (see below)
	BackgroundColor string                      `bson:"background_color" json:"background_color"` // Background color of the item on OpenSea. Must be a six-character hexadecimal without a pre-pended #.
	AnimationID     primitive.ObjectID          `bson:"animation_id" json:"animation_id"`
	YoutubeURL      string                      `bson:"youtube_url" json:"youtube_url"` // A URL to a YouTube video (only used if animation_url is not provided).
}

func (impl *NFTMetadataControllerImpl) validateCreateRequest(dirtyData *NFTMetadataCreateRequestIDO) error {
	e := make(map[string]string)

	if dirtyData.CollectionID.IsZero() {
		e["collection_id"] = "missing value"
	} else {
		doesExist, err := impl.CollectionStorer.CheckIfExistsByID(context.Background(), dirtyData.CollectionID)
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
	if dirtyData.Description == "" {
		e["description"] = "missing value"
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
		// Begin populating our NFT meta with the user add request data.
		nftmetadata := &s_d.NFTMetadata{
			CollectionID:    req.CollectionID,
			Name:            req.Name,
			ImageID:         req.ImageID,
			AnimationID:     req.AnimationID,
			ExternalURL:     req.ExternalURL,
			Description:     req.Description,
			Attributes:      req.Attributes,
			BackgroundColor: req.BackgroundColor,
			YoutubeURL:      req.YoutubeURL,
		}

		// Populate nftmetadata with default and tenant-specific information
		nftmetadata.TenantID = tenantID
		nftmetadata.TenantName = tenantName
		nftmetadata.TenantTimezone = tenantTimezone
		nftmetadata.ID = primitive.NewObjectID()
		nftmetadata.CreatedAt = time.Now()
		nftmetadata.ModifiedAt = time.Now()
		nftmetadata.Status = s_d.StatusActive

		impl.Logger.Debug("nft metadata creation beginning...",
			slog.String("id", nftmetadata.ID.Hex()))

		// DEVELOPER NOTE: We don't need to check if `collection` exists b/c
		// we already validated it before entering this transaction.
		// Lookup our image and attach related.
		collection, err := impl.CollectionStorer.GetByID(sessCtx, req.CollectionID)
		if err != nil {
			impl.Logger.Error("failed get collection by id", slog.Any("error", err))
			return nil, err
		}

		// Update our metadata record.
		nftmetadata.TokenID = collection.TokensCount
		nftmetadata.CollectionName = collection.Name

		// Update collection to have our new token id.
		collection.TokensCount = collection.TokensCount + 1
		collection.ModifiedAt = time.Now()
		if err := impl.CollectionStorer.UpdateByID(sessCtx, collection); err != nil {
			impl.Logger.Error("failed get collection by id", slog.Any("error", err))
			return nil, err
		}

		impl.Logger.Debug("nft collection generated new token id",
			slog.Uint64("token_id", collection.TokensCount))

		// DEVELOPER NOTE: We don't need to check if `imageAsset` exists b/c
		// we already validated it before entering this transaction.
		// Lookup our image and attach related.
		imageAsset, err := impl.NFTAssetStorer.GetByID(sessCtx, req.ImageID)
		if err != nil {
			impl.Logger.Error("failed get NFT asset by id", slog.Any("error", err))
			return nil, err
		}
		imageAsset.Status = nftasset_s.StatusPinning
		if err := impl.NFTAssetStorer.UpdateByID(sessCtx, imageAsset); err != nil {
			impl.Logger.Error("failed update NFT asset by id", slog.Any("error", err))
			return nil, err
		}

		// Update our metadata record.
		nftmetadata.Image = fmt.Sprintf("/ipfs/%v", imageAsset.CID)

		impl.Logger.Debug("image asset status changed")

		// Save the nftmetadata data to the database
		if err := impl.NFTMetadataStorer.Create(sessCtx, nftmetadata); err != nil {
			impl.Logger.Error("failed to save nftmetadata to database", slog.Any("error", err))
			return nil, err
		}

		impl.Logger.Debug("nft metadata created",
			slog.String("id", nftmetadata.ID.Hex()))

		//TODO: WE NEED TO GENERATE THE METADATA FOLDER IN OUR COLLECITON DIR.
		//TODO: UPDATE `MetadataFileCID`
		//TODO: UPDATE `IPNSPath`
		//TODO: WE NEED TO UPDATE IPNS WITH THE NEW COLLECTION DIR CONTENTS.

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
