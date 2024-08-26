package controller

import (
	"context"
	"time"

	domain "github.com/LuchaComics/monorepo/cloud/cps-pinws-backend/app/pinobject/datastore"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log/slog"
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

	// Generate the URL.
	fileURL, err := c.S3.GetPresignedURL(ctx, m.ObjectKey, 5*time.Minute)
	if err != nil {
		c.Logger.Error("s3 failed get presigned url error", slog.Any("error", err))
		return nil, err
	}

	m.ObjectURL = fileURL
	return m, err
}
