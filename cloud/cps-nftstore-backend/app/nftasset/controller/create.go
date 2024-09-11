package controller

import (
	"context"
	"log/slog"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	pinobj_s "github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/app/ipfsgateway/datastore"
	a_d "github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/app/nftasset/datastore"
	"github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/config/constants"
	"github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/utils/httperror"
)

type NFTAssetCreateRequestIDO struct {
	Name        string
	Filename    string `bson:"filename" json:"filename"`
	ContentType string `bson:"content_type" json:"content_type"`
	Data        []byte `bson:"data" json:"data"`
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

	if dirtyData.Filename == "" {
		e["filename"] = "missing value"
	}
	if dirtyData.ContentType == "" {
		e["content_type"] = "missing value"
	}
	if dirtyData.Data == nil {
		e["data"] = "missing value"
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

		//
		// STEP 1
		//

		// Add our file to the IPFS network (within our specific directory name).
		_, nftAssetFileCID, ipfsUploadErr := impl.IPFS.UploadBytesToDir(sessCtx, req.Data, req.Filename, "nftassets")
		if ipfsUploadErr != nil {
			impl.Logger.Error("failed uploading NFT asset file",
				slog.Any("error", ipfsUploadErr))
			return nil, err
		}
		impl.Logger.Debug("successfully uploaded to ipfs network",
			slog.String("cid", nftAssetFileCID))

		// Pin our upload file, which means IPFS will not delete this file when
		// it does its periodic garbage colleciton. If we didn't do this then
		// we will have this file deleted in the future.
		if ipfsPinErr := impl.IPFS.Pin(sessCtx, nftAssetFileCID); ipfsPinErr != nil {
			impl.Logger.Error("failed pinning to ipfs network",
				slog.Any("error", ipfsPinErr))
			return nil, err
		}

		impl.Logger.Debug("successfully pinned in ipfs network",
			slog.String("cid", nftAssetFileCID))

		//
		// STEP 2
		//

		// Extract from our session the following data.
		orgID, _ := sessCtx.Value(constants.SessionUserTenantID).(primitive.ObjectID)
		orgName, _ := sessCtx.Value(constants.SessionUserTenantName).(string)
		orgTimezone, _ := sessCtx.Value(constants.SessionUserTenantTimezone).(string)
		ipAddress, _ := sessCtx.Value(constants.SessionIPAddress).(string)

		// Create our record in the database.
		res := &a_d.NFTAsset{
			Status:                a_d.StatusPending, // Note: Change state when NFT metadata created.
			CID:                   nftAssetFileCID,
			Name:                  req.Name,
			Filename:              req.Filename,
			ContentType:           req.ContentType,
			TenantID:              orgID,
			TenantName:            orgName,
			TenantTimezone:        orgTimezone,
			ID:                    primitive.NewObjectID(),
			CreatedAt:             time.Now(),
			CreatedFromIPAddress:  ipAddress,
			ModifiedAt:            time.Now(),
			ModifiedFromIPAddress: ipAddress,
			NFTMetadataID:         primitive.NilObjectID,
		}

		// Save to database.
		if err := impl.NFTAssetStorer.Create(sessCtx, res); err != nil {
			impl.Logger.Error("database create error", slog.Any("error", err))
			return nil, err
		}

		impl.Logger.Debug("successfully created nft asset file",
			slog.String("filename", req.Filename),
			slog.String("content_type", req.ContentType),
			slog.String("cid", nftAssetFileCID),
			slog.String("status", res.Status),
			slog.String("id", res.ID.Hex()))

		//
		// STEP 3
		// Keep a record of our pinned object for the IPFS gateway.
		//

		pinObject := &pinobj_s.PinObject{
			ID:          primitive.NewObjectID(),
			IPNSPath:    "", // Set to empty b/c this pin is not mounted to IPNS.
			CID:         res.CID,
			Content:     nil,
			Filename:    res.Filename,
			ContentType: res.ContentType,
			CreatedAt:   res.CreatedAt,
			ModifiedAt:  res.ModifiedAt,
		}
		if createdErr := impl.PinObjectStorer.Create(sessCtx, pinObject); createdErr != nil {
			impl.Logger.Error("database create error", slog.Any("error", createdErr))
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
