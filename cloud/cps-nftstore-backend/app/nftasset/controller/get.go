package controller

import (
	"context"

	"log/slog"

	domain "github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/app/nftasset/datastore"
	"github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/utils/httperror"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (c *NFTAssetControllerImpl) GetByID(ctx context.Context, id primitive.ObjectID) (*domain.NFTAsset, error) {
	// // Extract from our session the following data.
	// userNFTAssetID := ctx.Value(constants.SessionUserNFTAssetID).(primitive.ObjectID)
	// userRole := ctx.Value(constants.SessionUserRole).(int8)
	//
	// If user is not administrator nor belongs to the nftasset then error.
	// if userRole != user_d.UserRoleRoot && id != userNFTAssetID {
	// 	c.Logger.Error("authenticated user is not staff role nor belongs to the nftasset error",
	// 		slog.Any("userRole", userRole),
	// 		slog.Any("userNFTAssetID", userNFTAssetID))
	// 	return nil, httperror.NewForForbiddenWithSingleField("message", "you do not belong to this nftasset")
	// }

	// Retrieve from our database the record for the specific id.
	m, err := c.NFTAssetStorer.GetByID(ctx, id)
	if err != nil {
		c.Logger.Error("database get by id error", slog.Any("error", err))
		return nil, err
	}
	if m == nil {
		c.Logger.Warn("does not exist", slog.Any("id", id))
		return nil, httperror.NewForNotFoundWithSingleField("request_id", "does not exist")
	}
	return m, err
}

//
// func (c *NFTAssetControllerImpl) GetByRequestID(ctx context.Context, requestID primitive.ObjectID) (*domain.NFTAsset, error) {
// 	// // Extract from our session the following data.
// 	// userNFTAssetID := ctx.Value(constants.SessionUserNFTAssetID).(primitive.ObjectID)
// 	// userRole := ctx.Value(constants.SessionUserRole).(int8)
// 	//
// 	// If user is not administrator nor belongs to the nftasset then error.
// 	// if userRole != user_d.UserRoleRoot && id != userNFTAssetID {
// 	// 	c.Logger.Error("authenticated user is not staff role nor belongs to the nftasset error",
// 	// 		slog.Any("userRole", userRole),
// 	// 		slog.Any("userNFTAssetID", userNFTAssetID))
// 	// 	return nil, httperror.NewForForbiddenWithSingleField("message", "you do not belong to this nftasset")
// 	// }
//
// 	// Retrieve from our database the record for the specific id.
// 	m, err := c.NFTAssetStorer.GetByRequestID(ctx, requestID)
// 	if err != nil {
// 		c.Logger.Error("database get by request id error", slog.Any("error", err))
// 		return nil, err
// 	}
// 	if m == nil {
// 		c.Logger.Warn("does not exist", slog.String("request_id", requestID.Hex()))
// 		return nil, httperror.NewForNotFoundWithSingleField("request_id", "does not exist")
// 	}
//
// 	// Generate the URL.
// 	fileURL, err := c.S3.GetPresignedURL(ctx, m.ObjectKey, 5*time.Minute)
// 	if err != nil {
// 		c.Logger.Error("s3 failed get presigned url error", slog.Any("error", err))
// 		return nil, err
// 	}
//
// 	m.ObjectURL = fileURL
// 	return m, err
// }
//
// func (impl *NFTAssetControllerImpl) GetWithContentByRequestID(ctx context.Context, requestID primitive.ObjectID) (*domain.NFTAsset, error) {
// 	// Retrieve from our database the record for the specific id.
// 	m, err := impl.NFTAssetStorer.GetByRequestID(ctx, requestID)
// 	if err != nil {
// 		impl.Logger.Error("database get by request id error", slog.Any("error", err))
// 		return nil, err
// 	}
// 	if m == nil {
// 		impl.Logger.Warn("does not exist", slog.String("request_id", requestID.Hex()))
// 		return nil, httperror.NewForNotFoundWithSingleField("request_id", "does not exist")
// 	}
//
// 	// content, err := impl.IPFS.GetContent(ctx, m.CID)
// 	// if err != nil {
// 	// 	impl.Logger.Error("get content by cid via ipfs error", slog.Any("error", err))
// 	// 	return nil, err
// 	// }
// 	//
// 	// m.Content = content
//
// 	return m, err
// }
