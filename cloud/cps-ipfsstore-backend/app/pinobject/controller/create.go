package controller

import (
	"context"
	"log/slog"
	"mime/multipart"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	a_d "github.com/LuchaComics/monorepo/cloud/cps-ipfsstore-backend/app/pinobject/datastore"
	"github.com/LuchaComics/monorepo/cloud/cps-ipfsstore-backend/config/constants"
	"github.com/LuchaComics/monorepo/cloud/cps-ipfsstore-backend/utils/httperror"
)

type PinObjectCreateRequestIDO struct {
	Name      string
	Origins   []string           `bson:"origins" json:"origins"`
	Meta      map[string]string  `bson:"meta" json:"meta"`
	ProjectID primitive.ObjectID // Outside of IPFS pinning spec.
	File      multipart.File     // Outside of IPFS pinning spec.
}

// PinObjectCreateResponseIDO represents `PinStatus` spec via https://ipfs.github.io/pinning-services-api-spec/#section/Schemas/Identifiers.
type PinObjectCreateResponseIDO struct {
	RequestID primitive.ObjectID `bson:"requestid" json:"requestid"`
	Status    string             `bson:"status" json:"status"`
	Created   time.Time          `bson:"created,omitempty" json:"created,omitempty"`
	Delegates []string           `bson:"delegates" json:"delegates"`
	Info      map[string]string  `bson:"info" json:"info"`
}

func ValidateCreateRequest(dirtyData *PinObjectCreateRequestIDO) error {
	e := make(map[string]string)

	if dirtyData.ProjectID.IsZero() {
		e["project_id"] = "missing value"
	}
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

func (impl *PinObjectControllerImpl) Create(ctx context.Context, req *PinObjectCreateRequestIDO) (*a_d.PinObject, error) {
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

		// Upload to IPFS network.
		_, cid, err := impl.IPFS.UploadMultipart(ctx, req.File, req.Meta["filename"], "uploads")
		if err != nil {
			impl.Logger.Error("failed uploading to IPFS", slog.Any("error", err))
			return nil, err
		}

		// Pin the file so it won't get deleted by IPFS garbage collection.
		if err := impl.IPFS.Pin(ctx, cid); err != nil {
			impl.Logger.Error("failed pinning to IPFS", slog.Any("error", err))
			return nil, err
		}

		// Extract from our session the following data.
		orgID, _ := sessCtx.Value(constants.SessionUserTenantID).(primitive.ObjectID)
		orgName, _ := sessCtx.Value(constants.SessionUserTenantName).(string)
		orgTimezone, _ := sessCtx.Value(constants.SessionUserTenantTimezone).(string)
		ipAdress, _ := sessCtx.Value(constants.SessionIPAddress).(string)

		// Create our meta record in the database.
		res := &a_d.PinObject{
			// Core fields required for a `pin` in IPFS.
			Status:    a_d.StatusPinned,
			CID:       cid,
			RequestID: primitive.NewObjectID(),
			Name:      req.Name,
			Created:   time.Now(),
			Origins:   req.Origins,
			Meta:      req.Meta,
			Delegates: make([]string, 0),
			Info:      make(map[string]string, 0),

			// Extension
			TenantID:              orgID,
			TenantName:            orgName,
			TenantTimezone:        orgTimezone,
			ID:                    primitive.NewObjectID(),
			ProjectID:             req.ProjectID,
			CreatedFromIPAddress:  ipAdress,
			ModifiedAt:            time.Now(),
			ModifiedFromIPAddress: ipAdress,

			// S3
			Filename: req.Meta["filename"],
			// ObjectKey: objectKey,
			// ObjectURL: "",
		}

		// Save to database.
		if err := impl.PinObjectStorer.Create(sessCtx, res); err != nil {
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

	return res.(*a_d.PinObject), nil
}
