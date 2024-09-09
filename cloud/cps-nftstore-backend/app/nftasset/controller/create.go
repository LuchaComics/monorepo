package controller

import (
	"context"
	"log/slog"
	"mime/multipart"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	a_d "github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/app/nftasset/datastore"
	"github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/config/constants"
	"github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/utils/httperror"
)

type NFTAssetCreateRequestIDO struct {
	Name string
	Meta map[string]string `bson:"meta" json:"meta"`
	File multipart.File    // Outside of IPFS pinning spec.
}

// NFTAssetCreateResponseIDO represents `PinStatus` spec via https://ipfs.github.io/pinning-services-api-spec/#section/Schemas/Identifiers.
type NFTAssetCreateResponseIDO struct {
	RequestID primitive.ObjectID `bson:"requestid" json:"requestid"`
	Status    string             `bson:"status" json:"status"`
	Created   time.Time          `bson:"created,omitempty" json:"created,omitempty"`
	Delegates []string           `bson:"delegates" json:"delegates"`
	Info      map[string]string  `bson:"info" json:"info"`
}

func ValidateCreateRequest(dirtyData *NFTAssetCreateRequestIDO) error {
	e := make(map[string]string)

	if dirtyData.Meta == nil {
		e["meta"] = "missing value"
	} else {
		if dirtyData.Meta["filename"] == "" {
			e["meta"] = "missing `filename` value"
		}
		if dirtyData.Meta["content_type"] == "" {
			e["meta"] = "missing `content_type` value"
		}
	}
	if dirtyData.File == nil {
		e["file"] = "missing value"
	}
	if len(e) != 0 {
		return httperror.NewForBadRequest(&e)
	}
	return nil
}

func (impl *NFTAssetControllerImpl) Create(ctx context.Context, req *NFTAssetCreateRequestIDO) (*a_d.NFTAsset, error) {
	if err := ValidateCreateRequest(req); err != nil {
		return nil, err
	}

	////
	//// Start the transaction.
	////

	session, err := impl.DbClient.StartSession()
	if err != nil {
		impl.Logger.Error("start session error",
			slog.Any("error", err))
		return nil, err
	}
	defer session.EndSession(ctx)

	// Define a transaction function with a series of operations
	transactionFunc := func(sessCtx mongo.SessionContext) (interface{}, error) {

		// Add our file to the IPFS network (within our specific directory name).
		dirCid, nftAssetFileCID, ipfsUploadErr := impl.IPFS.UploadMultipartToDir(sessCtx, req.File, req.Meta["filename"], "nftassets")
		if ipfsUploadErr != nil {
			impl.Logger.Error("failed uploading NFT asset file",
				slog.Any("error", ipfsUploadErr))
			return nil, err
		}
		impl.Logger.Debug("ipfs storage adapter successfully uploaded NFT asset file",
			slog.String("dir_cid", dirCid),
			slog.String("nft_asset_file_cid", nftAssetFileCID))

		// Extract from our session the following data.
		orgID, _ := sessCtx.Value(constants.SessionUserTenantID).(primitive.ObjectID)
		orgName, _ := sessCtx.Value(constants.SessionUserTenantName).(string)
		orgTimezone, _ := sessCtx.Value(constants.SessionUserTenantTimezone).(string)
		ipAdress, _ := sessCtx.Value(constants.SessionIPAddress).(string)

		// Create our record in the database.
		res := &a_d.NFTAsset{
			Status:                a_d.StatusPending, // Note: Change state when NFT metadata created.
			CID:                   nftAssetFileCID,
			Name:                  req.Name,
			CreatedAt:             time.Now(),
			Filename:              req.Meta["Filename"],
			ContentType:           req.Meta["ContentType"],
			TenantID:              orgID,
			TenantName:            orgName,
			TenantTimezone:        orgTimezone,
			ID:                    primitive.NewObjectID(),
			CreatedFromIPAddress:  ipAdress,
			ModifiedAt:            time.Now(),
			ModifiedFromIPAddress: ipAdress,
			NFTMetadataID:         primitive.NilObjectID,
		}

		// Save to database.
		if err := impl.NFTAssetStorer.Create(sessCtx, res); err != nil {
			impl.Logger.Error("database create error", slog.Any("error", err))
			return nil, err
		}
		return res, nil
	}

	// Start a transaction
	res, err := session.WithTransaction(ctx, transactionFunc)
	if err != nil {
		impl.Logger.Error("session failed error",
			slog.Any("error", err))
		return nil, err
	}

	return res.(*a_d.NFTAsset), nil
}
