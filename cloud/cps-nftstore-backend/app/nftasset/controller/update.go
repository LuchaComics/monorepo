package controller

import (
	"mime/multipart"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/utils/httperror"
)

type NFTAssetUpdateRequestIDO struct {
	RequestID     primitive.ObjectID
	Name          string
	NFTMetadataID primitive.ObjectID
	FileName      string
	FileType      string
	File          multipart.File
}

func ValidateUpdateRequest(dirtyData *NFTAssetUpdateRequestIDO) error {
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
	if dirtyData.NFTMetadataID.IsZero() {
		e["nftmetadata_id"] = "missing value"
	}
	if len(e) != 0 {
		return httperror.NewForBadRequest(&e)
	}
	return nil
}

// func (impl *NFTAssetControllerImpl) UpdateByRequestID(ctx context.Context, req *NFTAssetUpdateRequestIDO) (*domain.NFTAsset, error) {
// 	if err := ValidateUpdateRequest(req); err != nil {
// 		return nil, err
// 	}
//
// 	////
// 	//// Start the transaction.
// 	////
//
// 	session, err := impl.DbClient.StartSession()
// 	if err != nil {
// 		impl.Logger.Error("start session error",
// 			slog.Any("error", err),
// 			slog.Any("request_id", req.RequestID))
// 		return nil, err
// 	}
// 	defer session.EndSession(ctx)
//
// 	// Define a transaction function with a series of operations
// 	transactionFunc := func(sessCtx mongo.SessionContext) (interface{}, error) {
// 	    return nil, nil
//
// 		// // Fetch the original nftasset.
// 		// os, err := impl.NFTAssetStorer.GetByRequestID(sessCtx, req.RequestID)
// 		// if err != nil {
// 		// 	impl.Logger.Error("database get by id error",
// 		// 		slog.Any("error", err),
// 		// 		slog.Any("request_id", req.RequestID))
// 		// 	return nil, err
// 		// }
// 		// if os == nil {
// 		// 	impl.Logger.Error("nftasset does not exist error",
// 		// 		slog.Any("request_id", req.RequestID))
// 		// 	return nil, httperror.NewForBadRequestWithSingleField("message", "nftasset does not exist")
// 		// }
// 		//
// 		// // Extract from our session the following data.
// 		// userTenantID, _ := sessCtx.Value(constants.SessionUserTenantID).(primitive.ObjectID)
// 		// userTenantTimezone, _ := sessCtx.Value(constants.SessionUserTenantTimezone).(string)
// 		// userRole, _ := sessCtx.Value(constants.SessionUserRole).(int8)
// 		// ipAddress, _ := sessCtx.Value(constants.SessionIPAddress).(string)
// 		//
// 		// // If user is not administrator nor belongs to the nftasset then error.
// 		// if userRole != user_d.UserRoleRoot && os.TenantID != userTenantID {
// 		// 	impl.Logger.Error("authenticated user is not staff role nor belongs to the nftasset error",
// 		// 		slog.Any("userRole", userRole),
// 		// 		slog.Any("userTenantID", userTenantID))
// 		// 	return nil, httperror.NewForForbiddenWithSingleField("message", "you do not belong to this nftasset")
// 		// }
// 		//
// 		// 	// The following code will choose the directory we will upload based on the image type.
// 		// 	var directory string = "nftmetadatas"
// 		//
// 		// 	// Generate the key of our upload.
// 		// 	objectKey := fmt.Sprintf("%v/%v/%v", directory, req.NFTMetadataID.Hex(), req.FileName)
// 		//
// 		// 	// Update file.
// 		// 	os.ObjectKey = objectKey
// 		// 	os.Filename = req.FileName
// 		//
// 		// 	// // Upload to IPFS network.
// 		// 	// cid, err := impl.IPFS.AddFileContentFromMulipartFile(ctx, req.File)
// 		// 	// if err != nil {
// 		// 	// 	impl.Logger.Error("failed uploading to IPFS", slog.Any("error", err))
// 		// 	// 	return nil, err
// 		// 	// }
// 		//
// 		// 	// // Pin the file so it won't get deleted by IPFS garbage collection.
// 		// 	// if err := impl.IPFS.PinContent(ctx, cid); err != nil {
// 		// 	// 	impl.Logger.Error("failed pinning to IPFS", slog.Any("error", err))
// 		// 	// 	return nil, err
// 		// 	// }
// 		//
// 		// 	// os.CID = cid
// 		}
//
// 		// // Modify our original nftasset.
// 		// os.TenantTimezone = userTenantTimezone
// 		// os.ModifiedAt = time.Now()
// 		// os.ModifiedFromIPAddress = ipAddress
// 		// os.Name = req.Name
// 		// os.NFTMetadataID = req.NFTMetadataID
//
// 		// // Save to the database the modified nftasset.
// 		// if err := impl.NFTAssetStorer.UpdateByRequestID(sessCtx, os); err != nil {
// 		// 	impl.Logger.Error("database update by id error", slog.Any("error", err))
// 		// 	return nil, err
// 		// }
//
// 		// go func(org *domain.NFTAsset) {
// 		// 	impl.updateNFTAssetNameForAllUsers(sessCtx, org)
// 		// }(os)
// 		// go func(org *domain.NFTAsset) {
// 		// 	impl.updateNFTAssetNameForAllComicSubmissions(sessCtx, org)
// 		// }(os)
//
// 		return os, nil
// 	}
//
// 	// Start a transaction
// 	res, err := session.WithTransaction(ctx, transactionFunc)
// 	if err != nil {
// 		impl.Logger.Error("session failed error",
// 			slog.Any("error", err))
// 		return nil, err
// 	}
//
// 	return res.(*a_d.NFTAsset), nil
// }

// func (c *NFTAssetControllerImpl) updateNFTAssetNameForAllUsers(ctx context.Context, ns *domain.NFTAsset) error {
// 	impl.Logger.Debug("Beginning to update nftasset name for all uses")
// 	f := &user_d.UserListFilter{NFTAssetID: ns.ID}
// 	uu, err := impl.UserStorer.ListByFilter(ctx, f)
// 	if err != nil {
// 		impl.Logger.Error("database update by id error", slog.Any("error", err))
// 		return err
// 	}
// 	for _, u := range uu.Results {
// 		u.NFTAssetName = ns.Name
// 		log.Println("--->", u)
// 		// if err := impl.UserStorer.UpdateByID(ctx, u); err != nil {
// 		// 	impl.Logger.Error("database update by id error", slog.Any("error", err))
// 		// 	return err
// 		// }
// 	}
// 	return nil
// }
//
// func (c *NFTAssetControllerImpl) updateNFTAssetNameForAllComicSubmissions(ctx context.Context, ns *domain.NFTAsset) error {
// 	impl.Logger.Debug("Beginning to update nftasset name for all submissions")
// 	f := &domain.ComicSubmissionListFilter{NFTAssetID: ns.ID}
// 	uu, err := impl.ComicSubmissionStorer.ListByFilter(ctx, f)
// 	if err != nil {
// 		impl.Logger.Error("database update by id error", slog.Any("error", err))
// 		return err
// 	}
// 	for _, u := range uu.Results {
// 		u.NFTAssetName = ns.Name
// 		log.Println("--->", u)
// 		// if err := impl.ComicSubmissionStorer.UpdateByID(ctx, u); err != nil {
// 		// 	impl.Logger.Error("database update by id error", slog.Any("error", err))
// 		// 	return err
// 		// }
// 	}
// 	return nil
// }
