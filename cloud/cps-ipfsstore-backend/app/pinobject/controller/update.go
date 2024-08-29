package controller

import (
	"context"
	"fmt"
	"log/slog"
	"mime/multipart"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	a_d "github.com/LuchaComics/monorepo/cloud/cps-ipfsstore-backend/app/pinobject/datastore"
	domain "github.com/LuchaComics/monorepo/cloud/cps-ipfsstore-backend/app/pinobject/datastore"
	user_d "github.com/LuchaComics/monorepo/cloud/cps-ipfsstore-backend/app/user/datastore"
	"github.com/LuchaComics/monorepo/cloud/cps-ipfsstore-backend/config/constants"
	"github.com/LuchaComics/monorepo/cloud/cps-ipfsstore-backend/utils/httperror"
)

type PinObjectUpdateRequestIDO struct {
	RequestID primitive.ObjectID
	Name      string
	ProjectID primitive.ObjectID
	FileName  string
	FileType  string
	File      multipart.File
}

func ValidateUpdateRequest(dirtyData *PinObjectUpdateRequestIDO) error {
	e := make(map[string]string)

	if dirtyData.RequestID.IsZero() {
		e["requestid"] = "missing value"
	}
	if dirtyData.FileName == "" {
		e["meta"] = "missing `filename` value"
	}
	if dirtyData.FileType == "" {
		e["meta"] = "missing `content_type` value"
	}
	if dirtyData.ProjectID.IsZero() {
		e["project_id"] = "missing value"
	}
	if len(e) != 0 {
		return httperror.NewForBadRequest(&e)
	}
	return nil
}

func (impl *PinObjectControllerImpl) UpdateByRequestID(ctx context.Context, req *PinObjectUpdateRequestIDO) (*domain.PinObject, error) {
	if err := ValidateUpdateRequest(req); err != nil {
		return nil, err
	}

	////
	//// Start the transaction.
	////

	session, err := impl.DbClient.StartSession()
	if err != nil {
		impl.Logger.Error("start session error",
			slog.Any("error", err),
			slog.Any("request_id", req.RequestID))
		return nil, err
	}
	defer session.EndSession(ctx)

	// Define a transaction function with a series of operations
	transactionFunc := func(sessCtx mongo.SessionContext) (interface{}, error) {

		// Fetch the original pinobject.
		os, err := impl.PinObjectStorer.GetByRequestID(sessCtx, req.RequestID)
		if err != nil {
			impl.Logger.Error("database get by id error",
				slog.Any("error", err),
				slog.Any("request_id", req.RequestID))
			return nil, err
		}
		if os == nil {
			impl.Logger.Error("pinobject does not exist error",
				slog.Any("request_id", req.RequestID))
			return nil, httperror.NewForBadRequestWithSingleField("message", "pinobject does not exist")
		}

		// Extract from our session the following data.
		userTenantID, _ := sessCtx.Value(constants.SessionUserTenantID).(primitive.ObjectID)
		userTenantTimezone, _ := sessCtx.Value(constants.SessionUserTenantTimezone).(string)
		userRole, _ := sessCtx.Value(constants.SessionUserRole).(int8)
		ipAdress, _ := sessCtx.Value(constants.SessionIPAddress).(string)

		// If user is not administrator nor belongs to the pinobject then error.
		if userRole != user_d.UserRoleRoot && os.TenantID != userTenantID {
			impl.Logger.Error("authenticated user is not staff role nor belongs to the pinobject error",
				slog.Any("userRole", userRole),
				slog.Any("userTenantID", userTenantID))
			return nil, httperror.NewForForbiddenWithSingleField("message", "you do not belong to this pinobject")
		}

		// Update the file if the user uploaded a new file.
		if req.File != nil {
			// Proceed to delete the physical files from AWS s3.
			if err := impl.S3.DeleteByKeys(sessCtx, []string{os.ObjectKey}); err != nil {
				impl.Logger.Warn("s3 delete by keys error", slog.Any("error", err))
				// Do not return an error, simply continue this function as there might
				// be a case were the file was removed on the s3 bucket by ourselves
				// or some other reason.
			}
			// Proceed to delete the physical files from IPFS.
			if err := impl.IPFS.DeleteContent(sessCtx, os.CID); err != nil {
				impl.Logger.Warn("ipfs delete by CID error", slog.Any("error", err))
				// Do not return an error, simply continue this function as there might
				// be a case were the file was removed on the s3 bucket by ourselves
				// or some other reason.
			}

			// The following code will choose the directory we will upload based on the image type.
			var directory string = "projects"

			// Generate the key of our upload.
			objectKey := fmt.Sprintf("%v/%v/%v", directory, req.ProjectID.Hex(), req.FileName)

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

			// Update file.
			os.ObjectKey = objectKey
			os.Filename = req.FileName

			// Upload to IPFS network.
			cid, err := impl.IPFS.AddFileContentFromMulipartFile(ctx, req.File)
			if err != nil {
				impl.Logger.Error("failed uploading to IPFS", slog.Any("error", err))
				return nil, err
			}

			// Pin the file so it won't get deleted by IPFS garbage collection.
			if err := impl.IPFS.PinContent(ctx, cid); err != nil {
				impl.Logger.Error("failed pinning to IPFS", slog.Any("error", err))
				return nil, err
			}

			os.CID = cid
		}

		// Modify our original pinobject.
		os.TenantTimezone = userTenantTimezone
		os.ModifiedAt = time.Now()
		os.ModifiedFromIPAddress = ipAdress
		os.Name = req.Name
		os.ProjectID = req.ProjectID

		// Save to the database the modified pinobject.
		if err := impl.PinObjectStorer.UpdateByRequestID(sessCtx, os); err != nil {
			impl.Logger.Error("database update by id error", slog.Any("error", err))
			return nil, err
		}

		// go func(org *domain.PinObject) {
		// 	impl.updatePinObjectNameForAllUsers(sessCtx, org)
		// }(os)
		// go func(org *domain.PinObject) {
		// 	impl.updatePinObjectNameForAllComicSubmissions(sessCtx, org)
		// }(os)

		return os, nil
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

// func (c *PinObjectControllerImpl) updatePinObjectNameForAllUsers(ctx context.Context, ns *domain.PinObject) error {
// 	impl.Logger.Debug("Beginning to update pinobject name for all uses")
// 	f := &user_d.UserListFilter{PinObjectID: ns.ID}
// 	uu, err := impl.UserStorer.ListByFilter(ctx, f)
// 	if err != nil {
// 		impl.Logger.Error("database update by id error", slog.Any("error", err))
// 		return err
// 	}
// 	for _, u := range uu.Results {
// 		u.PinObjectName = ns.Name
// 		log.Println("--->", u)
// 		// if err := impl.UserStorer.UpdateByID(ctx, u); err != nil {
// 		// 	impl.Logger.Error("database update by id error", slog.Any("error", err))
// 		// 	return err
// 		// }
// 	}
// 	return nil
// }
//
// func (c *PinObjectControllerImpl) updatePinObjectNameForAllComicSubmissions(ctx context.Context, ns *domain.PinObject) error {
// 	impl.Logger.Debug("Beginning to update pinobject name for all submissions")
// 	f := &domain.ComicSubmissionListFilter{PinObjectID: ns.ID}
// 	uu, err := impl.ComicSubmissionStorer.ListByFilter(ctx, f)
// 	if err != nil {
// 		impl.Logger.Error("database update by id error", slog.Any("error", err))
// 		return err
// 	}
// 	for _, u := range uu.Results {
// 		u.PinObjectName = ns.Name
// 		log.Println("--->", u)
// 		// if err := impl.ComicSubmissionStorer.UpdateByID(ctx, u); err != nil {
// 		// 	impl.Logger.Error("database update by id error", slog.Any("error", err))
// 		// 	return err
// 		// }
// 	}
// 	return nil
// }
