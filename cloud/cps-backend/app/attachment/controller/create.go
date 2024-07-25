package controller

import (
	"context"
	"fmt"
	"log/slog"
	"mime/multipart"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	a_d "github.com/LuchaComics/monorepo/cloud/cps-backend/app/attachment/datastore"
	"github.com/LuchaComics/monorepo/cloud/cps-backend/config/constants"
	"github.com/LuchaComics/monorepo/cloud/cps-backend/utils/httperror"
)

type AttachmentCreateRequestIDO struct {
	Name          string
	Description   string
	OwnershipID   primitive.ObjectID
	OwnershipType int8
	FileName      string
	FileType      string
	File          multipart.File
}

func ValidateCreateRequest(dirtyData *AttachmentCreateRequestIDO) error {
	e := make(map[string]string)

	if dirtyData.Name == "" {
		e["name"] = "missing value"
	}
	if dirtyData.Description == "" {
		e["description"] = "missing value"
	}
	if dirtyData.OwnershipID.IsZero() {
		e["ownership_id"] = "missing value"
	}
	if dirtyData.OwnershipType == 0 {
		e["ownership_type"] = "missing value"
	}
	if dirtyData.FileName == "" {
		e["file"] = "missing value"
	}
	if len(e) != 0 {
		return httperror.NewForBadRequest(&e)
	}
	return nil
}

func (impl *AttachmentControllerImpl) Create(ctx context.Context, req *AttachmentCreateRequestIDO) (*a_d.Attachment, error) {
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
		var directory string
		switch req.OwnershipType {
		case a_d.OwnershipTypeUser:
			directory = "user"
		case a_d.OwnershipTypeSubmission:
			directory = "submission"
		case a_d.OwnershipTypeStore:
			directory = "store"
		default:
			impl.Logger.Error("unsupported ownership type format", slog.Any("ownership_type", req.OwnershipType))
			return nil, fmt.Errorf("unsuported iownership type  of %v, please pick another type", req.OwnershipType)
		}

		// Generate the key of our upload.
		objectKey := fmt.Sprintf("%v/%v/%v", directory, req.OwnershipID.Hex(), req.FileName)

		// For debugging purposes only.
		impl.Logger.Debug("pre-upload meta",
			slog.String("FileName", req.FileName),
			slog.String("FileType", req.FileType),
			slog.String("Directory", directory),
			slog.String("ObjectKey", objectKey),
			slog.String("Name", req.Name),
			slog.String("Desc", req.Description),
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

		// Extract from our session the following data.
		orgID, _ := sessCtx.Value(constants.SessionUserStoreID).(primitive.ObjectID)
		orgName, _ := sessCtx.Value(constants.SessionUserStoreName).(string)
		orgTimezone, _ := sessCtx.Value(constants.SessionUserStoreTimezone).(string)
		userID, _ := sessCtx.Value(constants.SessionUserID).(primitive.ObjectID)
		userName, _ := sessCtx.Value(constants.SessionUserName).(string)

		// Create our meta record in the database.
		res := &a_d.Attachment{
			StoreID:            orgID,
			StoreName:          orgName,
			StoreTimezone:      orgTimezone,
			ID:                 primitive.NewObjectID(),
			CreatedAt:          time.Now(),
			CreatedByUserName:  userName,
			CreatedByUserID:    userID,
			ModifiedAt:         time.Now(),
			ModifiedByUserName: userName,
			ModifiedByUserID:   userID,
			Name:               req.Name,
			Description:        req.Description,
			Filename:           req.FileName,
			ObjectKey:          objectKey,
			ObjectURL:          "",
			OwnershipID:        req.OwnershipID,
			OwnershipType:      req.OwnershipType,
			Status:             a_d.StatusActive,
		}
		err := impl.AttachmentStorer.Create(sessCtx, res)
		if err != nil {
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

	return res.(*a_d.Attachment), nil
}
