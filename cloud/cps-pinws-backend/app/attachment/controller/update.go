package controller

import (
	"context"
	"fmt"
	"log/slog"
	"mime/multipart"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	a_d "github.com/LuchaComics/monorepo/cloud/cps-pinws-backend/app/attachment/datastore"
	domain "github.com/LuchaComics/monorepo/cloud/cps-pinws-backend/app/attachment/datastore"
	user_d "github.com/LuchaComics/monorepo/cloud/cps-pinws-backend/app/user/datastore"
	"github.com/LuchaComics/monorepo/cloud/cps-pinws-backend/config/constants"
	"github.com/LuchaComics/monorepo/cloud/cps-pinws-backend/utils/httperror"
)

type AttachmentUpdateRequestIDO struct {
	ID            primitive.ObjectID
	Name          string
	Description   string
	OwnershipID   primitive.ObjectID
	OwnershipType int8
	FileName      string
	FileType      string
	File          multipart.File
}

func ValidateUpdateRequest(dirtyData *AttachmentUpdateRequestIDO) error {
	e := make(map[string]string)

	if dirtyData.ID.IsZero() {
		e["id"] = "missing value"
	}
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
	if len(e) != 0 {
		return httperror.NewForBadRequest(&e)
	}
	return nil
}

func (impl *AttachmentControllerImpl) UpdateByID(ctx context.Context, req *AttachmentUpdateRequestIDO) (*domain.Attachment, error) {
	if err := ValidateUpdateRequest(req); err != nil {
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

		// Fetch the original attachment.
		os, err := impl.AttachmentStorer.GetByID(sessCtx, req.ID)
		if err != nil {
			impl.Logger.Error("database get by id error",
				slog.Any("error", err),
				slog.Any("attachment_id", req.ID))
			return nil, err
		}
		if os == nil {
			impl.Logger.Error("attachment does not exist error",
				slog.Any("attachment_id", req.ID))
			return nil, httperror.NewForBadRequestWithSingleField("message", "attachment does not exist")
		}

		// Extract from our session the following data.
		userID, _ := sessCtx.Value(constants.SessionUserID).(primitive.ObjectID)
		userTenantID, _ := sessCtx.Value(constants.SessionUserTenantID).(primitive.ObjectID)
		userTenantTimezone, _ := sessCtx.Value(constants.SessionUserTenantTimezone).(string)
		userRole, _ := sessCtx.Value(constants.SessionUserRole).(int8)
		userName, _ := sessCtx.Value(constants.SessionUserName).(string)

		// If user is not administrator nor belongs to the attachment then error.
		if userRole != user_d.UserRoleRoot && os.TenantID != userTenantID {
			impl.Logger.Error("authenticated user is not staff role nor belongs to the attachment error",
				slog.Any("userRole", userRole),
				slog.Any("userTenantID", userTenantID))
			return nil, httperror.NewForForbiddenWithSingleField("message", "you do not belong to this attachment")
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

			os.CID = cid
		}

		// Modify our original attachment.
		os.TenantTimezone = userTenantTimezone
		os.ModifiedAt = time.Now()
		os.ModifiedByUserID = userID
		os.ModifiedByUserName = userName
		os.Name = req.Name
		os.Description = req.Description
		os.OwnershipID = req.OwnershipID
		os.OwnershipType = req.OwnershipType

		// Save to the database the modified attachment.
		if err := impl.AttachmentStorer.UpdateByID(sessCtx, os); err != nil {
			impl.Logger.Error("database update by id error", slog.Any("error", err))
			return nil, err
		}

		// go func(org *domain.Attachment) {
		// 	impl.updateAttachmentNameForAllUsers(sessCtx, org)
		// }(os)
		// go func(org *domain.Attachment) {
		// 	impl.updateAttachmentNameForAllComicSubmissions(sessCtx, org)
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

	return res.(*a_d.Attachment), nil
}

// func (c *AttachmentControllerImpl) updateAttachmentNameForAllUsers(ctx context.Context, ns *domain.Attachment) error {
// 	impl.Logger.Debug("Beginning to update attachment name for all uses")
// 	f := &user_d.UserListFilter{AttachmentID: ns.ID}
// 	uu, err := impl.UserStorer.ListByFilter(ctx, f)
// 	if err != nil {
// 		impl.Logger.Error("database update by id error", slog.Any("error", err))
// 		return err
// 	}
// 	for _, u := range uu.Results {
// 		u.AttachmentName = ns.Name
// 		log.Println("--->", u)
// 		// if err := impl.UserStorer.UpdateByID(ctx, u); err != nil {
// 		// 	impl.Logger.Error("database update by id error", slog.Any("error", err))
// 		// 	return err
// 		// }
// 	}
// 	return nil
// }
//
// func (c *AttachmentControllerImpl) updateAttachmentNameForAllComicSubmissions(ctx context.Context, ns *domain.Attachment) error {
// 	impl.Logger.Debug("Beginning to update attachment name for all submissions")
// 	f := &domain.ComicSubmissionListFilter{AttachmentID: ns.ID}
// 	uu, err := impl.ComicSubmissionStorer.ListByFilter(ctx, f)
// 	if err != nil {
// 		impl.Logger.Error("database update by id error", slog.Any("error", err))
// 		return err
// 	}
// 	for _, u := range uu.Results {
// 		u.AttachmentName = ns.Name
// 		log.Println("--->", u)
// 		// if err := impl.ComicSubmissionStorer.UpdateByID(ctx, u); err != nil {
// 		// 	impl.Logger.Error("database update by id error", slog.Any("error", err))
// 		// 	return err
// 		// }
// 	}
// 	return nil
// }
