package controller

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	a_d "github.com/LuchaComics/monorepo/cloud/cps-ipfsstore-backend/app/pinobject/datastore"
	pin_s "github.com/LuchaComics/monorepo/cloud/cps-ipfsstore-backend/app/pinobject/datastore"
	project_s "github.com/LuchaComics/monorepo/cloud/cps-ipfsstore-backend/app/project/datastore"
	"github.com/LuchaComics/monorepo/cloud/cps-ipfsstore-backend/config/constants"
	"github.com/LuchaComics/monorepo/cloud/cps-ipfsstore-backend/utils/httperror"
)

type IpfsAddRequestIDO struct {
	ApiKey      string `bson:"api_key" json:"api_key"`
	Filename    string `bson:"filename" json:"filename"`
	ContentType string `bson:"content_type" json:"content_type"`
	Data        []byte `bson:"data" json:"data"`
}

// IpfsAddResponseIDO represents `PinStatus` spec via https://ipfs.github.io/pinning-services-api-spec/#section/Schemas/Identifiers.
type IpfsAddResponseIDO struct {
	CID string `bson:"cid" json:"cid"`
}

func ValidateIpfsAddRequest(dirtyData *IpfsAddRequestIDO) error {
	e := make(map[string]string)

	if dirtyData.ApiKey == "" {
		e["api_key"] = "missing value"
	}
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

func (impl *PinObjectControllerImpl) validateApiKey(sessCtx mongo.SessionContext, apiKey string) (*project_s.Project, error) {
	// DEVELOPERS NOTE:
	// Verify API token, to see how it was set, please see `app/project/controller/create`.

	apiKeyDecoded, err := impl.JWT.ProcessJWTToken(apiKey)
	if err != nil {
		return nil, httperror.NewForUnauthorizedWithSingleField("api_key", fmt.Sprintf("bad formatting: %v", err))
	}
	apiKeyPayload := strings.Split(apiKeyDecoded, "@")
	if len(apiKeyPayload) < 2 {
		return nil, httperror.NewForUnauthorizedWithSingleField("api_key", "corrupted payload: bad structure")
	}
	if apiKeyPayload[0] == "" {
		return nil, httperror.NewForUnauthorizedWithSingleField("api_key", "corrupted payload: missing `project_id`")
	}
	if apiKeyPayload[1] == "" {
		return nil, httperror.NewForUnauthorizedWithSingleField("api_key", "corrupted payload: missing `secret`")
	}
	projectID, err := primitive.ObjectIDFromHex(apiKeyPayload[0])
	if err != nil {
		return nil, httperror.NewForUnauthorizedWithSingleField("api_key", "corrupted payload: `project_id` is invalid")
	}
	project, err := impl.ProjectStorer.GetByID(sessCtx, projectID)
	if err != nil {
		return nil, httperror.NewForUnauthorizedWithSingleField("api_key", fmt.Sprintf("invalid: error getting project: %v", err))
	}
	if project == nil {
		return nil, httperror.NewForUnauthorizedWithSingleField("api_key", "invalid: `project_id` does not exist")
	}

	// Verify the api key secret and project hashed secret match.
	passwordMatch, _ := impl.Password.ComparePasswordAndHash(apiKeyPayload[1], project.SecretHash)
	if passwordMatch == false {
		return nil, httperror.NewForUnauthorizedWithSingleField("api_key", "unauthorized")
	}
	return project, nil
}

func (impl *PinObjectControllerImpl) IpfsAdd(ctx context.Context, req *IpfsAddRequestIDO) (string, error) {
	if err := ValidateIpfsAddRequest(req); err != nil {
		return "", err
	}

	////
	//// Start the transaction.
	////

	session, err := impl.DbClient.StartSession()
	if err != nil {
		impl.Logger.Error("start session error",
			slog.Any("error", err))
		return "", err
	}
	defer session.EndSession(ctx)

	// Define a transaction function with a series of operations
	transactionFunc := func(sessCtx mongo.SessionContext) (interface{}, error) {
		project, err := impl.validateApiKey(sessCtx, req.ApiKey)
		if err != nil {
			impl.Logger.Error("api key validation failed",
				slog.Any("error", err))
			return nil, err
		}

		// Upload to IPFS network.
		_, cid, err := impl.IPFS.UploadBytes(ctx, req.Data, req.Filename, fmt.Sprintf("dir_%v", project.ID.Hex()))
		if err != nil {
			impl.Logger.Error("failed uploading to IPFS", slog.Any("error", err))
			return nil, err
		}

		// Pin the file so it won't get deleted by IPFS garbage collection.
		if err := impl.IPFS.Pin(ctx, cid); err != nil {
			impl.Logger.Error("failed pinning to IPFS", slog.Any("error", err))
			return nil, err
		}

		// // Add to S3.
		// // Generate the key of our upload.
		// objectKey := fmt.Sprintf("%v/%v/%v/%v/%v", "projects", project.ID.Hex(), "cids", cid, req.Filename)
		// go func(content []byte, objkey string) {
		// 	impl.Logger.Debug("beginning private s3 upload...")
		// 	if err := impl.S3.UploadContentFromBytes(context.Background(), objkey, content); err != nil {
		// 		impl.Logger.Error("private s3 upload error", slog.Any("error", err))
		// 		// Do not return an error, simply continue this function as there might
		// 		// be a case were the file was removed on the s3 bucket by ourselves
		// 		// or some other reason.
		// 	}
		// 	impl.Logger.Debug("Finished private s3 upload")
		// }(req.Data, objectKey)

		// For developer purposes only.
		ipAdress, _ := sessCtx.Value(constants.SessionIPAddress).(string)
		log.Println("Project ID:", project.ID)
		log.Printf("Received file: %s\n", req.Filename)
		log.Printf("Content-Type: %s\n", req.ContentType)

		origins := make([]string, 0)
		meta := make(map[string]string, 0)

		// Handle meta. We will attach meta along with some custom fields.
		meta["filename"] = req.Filename
		meta["content_type"] = req.ContentType

		// Initialize our array which will store all the results from the remote server.
		pinObject := &pin_s.PinObject{
			// Core fields required for a `pin` in IPFS.
			Status:    a_d.StatusPinned,
			CID:       cid,
			RequestID: primitive.NewObjectID(),
			Name:      "", // Blank b/c it's optional.
			Created:   time.Now(),
			Origins:   origins,
			Meta:      meta,
			Delegates: make([]string, 0),
			Info:      make(map[string]string, 0),

			// Extension
			TenantID:              project.TenantID,
			TenantName:            project.TenantName,
			TenantTimezone:        project.TenantTimezone,
			ID:                    primitive.NewObjectID(),
			ProjectID:             project.ID,
			CreatedFromIPAddress:  ipAdress,
			ModifiedAt:            time.Now(),
			ModifiedFromIPAddress: ipAdress,

			// S3
			Filename: req.Filename,
			// ObjectKey: objectKey,
			// ObjectURL: "",
		}
		if err := impl.PinObjectStorer.Create(sessCtx, pinObject); err != nil {
			impl.Logger.Error("database create error", slog.Any("error", err))
			return nil, err
		}

		return cid, nil
	}

	// Start a transaction
	res, err := session.WithTransaction(ctx, transactionFunc)
	if err != nil {
		impl.Logger.Error("session failed error",
			slog.Any("error", err))
		return "", err
	}

	return res.(string), nil
}
