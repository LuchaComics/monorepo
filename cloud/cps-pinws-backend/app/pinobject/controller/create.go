package controller

import (
	"context"
	"fmt"
	"log/slog"
	"mime/multipart"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	a_d "github.com/LuchaComics/monorepo/cloud/cps-pinws-backend/app/pinobject/datastore"
	"github.com/LuchaComics/monorepo/cloud/cps-pinws-backend/config/constants"
	"github.com/LuchaComics/monorepo/cloud/cps-pinws-backend/utils/httperror"
)

type PinObjectCreateRequestIDO struct {
	Name      string
	ProjectID primitive.ObjectID
	FileName  string
	FileType  string
	File      multipart.File
	Origins   []string          `bson:"origins" json:"origins"`
	Meta      map[string]string `bson:"meta" json:"meta"`
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
	if dirtyData.FileName == "" {
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

		// The following code will choose the directory we will upload based on the image type.
		var directory string = "projects"

		// Generate the key of our upload.
		objectKey := fmt.Sprintf("%v/%v/%v", directory, req.ProjectID.Hex(), req.FileName)

		// For debugging purposes only.
		impl.Logger.Debug("pre-upload meta",
			slog.String("FileName", req.FileName),
			slog.String("FileType", req.FileType),
			slog.String("Directory", directory),
			slog.String("ObjectKey", objectKey),
			slog.String("Name", req.Name),
		)

		go func(file multipart.File, objkey string) {
			impl.Logger.Debug("beginning private s3 image upload...")
			if err := impl.S3.UploadContentFromMulipart(context.Background(), objkey, file); err != nil {
				impl.Logger.Error("private s3 upload error", slog.Any("error", err))
				// Do not return an error, simply continue this function as there might
				// be a case were the file was removed on the s3 bucket by ourselves
				// or some other reason.
			}
			impl.Logger.Debug("Finished private s3 image upload")
		}(req.File, objectKey)

		// Upload to IPFS network.
		cid, err := impl.IPFS.UploadContentFromMulipart(ctx, req.File)
		if err != nil {
			impl.Logger.Error("failed uploading to IPFS", slog.Any("error", err))
			return nil, err
		}

		// Pin the file so it won't get deleted by IPFS garbage collection.
		if err := impl.IPFS.PinContent(ctx, cid); err != nil {
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
			Filename:  req.FileName,
			ObjectKey: objectKey,
			ObjectURL: "",
		}
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
