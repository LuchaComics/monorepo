package controller

import (
	"context"
	"time"

	"log/slog"

	domain "github.com/LuchaComics/monorepo/cloud/cps-ipfsstore-backend/app/pinobject/datastore"
	"github.com/LuchaComics/monorepo/cloud/cps-ipfsstore-backend/utils/httperror"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (c *PinObjectControllerImpl) GetByID(ctx context.Context, id primitive.ObjectID) (*domain.PinObject, error) {
	// // Extract from our session the following data.
	// userPinObjectID := ctx.Value(constants.SessionUserPinObjectID).(primitive.ObjectID)
	// userRole := ctx.Value(constants.SessionUserRole).(int8)
	//
	// If user is not administrator nor belongs to the pinobject then error.
	// if userRole != user_d.UserRoleRoot && id != userPinObjectID {
	// 	c.Logger.Error("authenticated user is not staff role nor belongs to the pinobject error",
	// 		slog.Any("userRole", userRole),
	// 		slog.Any("userPinObjectID", userPinObjectID))
	// 	return nil, httperror.NewForForbiddenWithSingleField("message", "you do not belong to this pinobject")
	// }

	// Retrieve from our database the record for the specific id.
	m, err := c.PinObjectStorer.GetByID(ctx, id)
	if err != nil {
		c.Logger.Error("database get by id error", slog.Any("error", err))
		return nil, err
	}
	if m == nil {
		c.Logger.Warn("does not exist", slog.Any("id", id))
		return nil, httperror.NewForNotFoundWithSingleField("request_id", "does not exist")
	}

	// Generate the URL.
	fileURL, err := c.S3.GetPresignedURL(ctx, m.ObjectKey, 5*time.Minute)
	if err != nil {
		c.Logger.Error("s3 failed get presigned url error", slog.Any("error", err))
		return nil, err
	}

	m.ObjectURL = fileURL
	return m, err
}

func (c *PinObjectControllerImpl) GetByRequestID(ctx context.Context, requestID primitive.ObjectID) (*domain.PinObject, error) {
	// // Extract from our session the following data.
	// userPinObjectID := ctx.Value(constants.SessionUserPinObjectID).(primitive.ObjectID)
	// userRole := ctx.Value(constants.SessionUserRole).(int8)
	//
	// If user is not administrator nor belongs to the pinobject then error.
	// if userRole != user_d.UserRoleRoot && id != userPinObjectID {
	// 	c.Logger.Error("authenticated user is not staff role nor belongs to the pinobject error",
	// 		slog.Any("userRole", userRole),
	// 		slog.Any("userPinObjectID", userPinObjectID))
	// 	return nil, httperror.NewForForbiddenWithSingleField("message", "you do not belong to this pinobject")
	// }

	// Retrieve from our database the record for the specific id.
	m, err := c.PinObjectStorer.GetByRequestID(ctx, requestID)
	if err != nil {
		c.Logger.Error("database get by request id error", slog.Any("error", err))
		return nil, err
	}
	if m == nil {
		c.Logger.Warn("does not exist", slog.String("request_id", requestID.Hex()))
		return nil, httperror.NewForNotFoundWithSingleField("request_id", "does not exist")
	}

	// Generate the URL.
	fileURL, err := c.S3.GetPresignedURL(ctx, m.ObjectKey, 5*time.Minute)
	if err != nil {
		c.Logger.Error("s3 failed get presigned url error", slog.Any("error", err))
		return nil, err
	}

	m.ObjectURL = fileURL
	return m, err
}

func (impl *PinObjectControllerImpl) GetWithContentByRequestID(ctx context.Context, requestID primitive.ObjectID) (*domain.PinObject, error) {
	// Retrieve from our database the record for the specific id.
	m, err := impl.PinObjectStorer.GetByRequestID(ctx, requestID)
	if err != nil {
		impl.Logger.Error("database get by request id error", slog.Any("error", err))
		return nil, err
	}
	if m == nil {
		impl.Logger.Warn("does not exist", slog.String("request_id", requestID.Hex()))
		return nil, httperror.NewForNotFoundWithSingleField("request_id", "does not exist")
	}

	content, err := impl.IPFS.GetContent(ctx, m.CID)
	if err != nil {
		impl.Logger.Error("get content by cid via ipfs error", slog.Any("error", err))
		return nil, err
	}

	m.Content = content

	return m, err
}
