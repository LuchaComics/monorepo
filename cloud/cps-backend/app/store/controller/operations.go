package controller

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"log/slog"

	org_d "github.com/LuchaComics/monorepo/cloud/cps-backend/app/store/datastore"
	"github.com/LuchaComics/monorepo/cloud/cps-backend/config/constants"
	"github.com/LuchaComics/monorepo/cloud/cps-backend/utils/httperror"
)

func (c *StoreControllerImpl) CreateComment(ctx context.Context, storeID primitive.ObjectID, content string) (*org_d.Store, error) {
	// Fetch the original customer.
	s, err := c.StoreStorer.GetByID(ctx, storeID)
	if err != nil {
		c.Logger.Error("database get by id error", slog.Any("error", err))
		return nil, err
	}
	if s == nil {
		c.Logger.Error("store does not exist error",
			slog.Any("store_id", storeID))
		return nil, httperror.NewForBadRequestWithSingleField("message", "store does not exist")
	}

	// Create our comment.
	comment := &org_d.StoreComment{
		ID:               primitive.NewObjectID(),
		Content:          content,
		StoreID:   ctx.Value(constants.SessionUserStoreID).(primitive.ObjectID),
		CreatedByUserID:  ctx.Value(constants.SessionUserID).(primitive.ObjectID),
		CreatedByName:    ctx.Value(constants.SessionUserName).(string),
		CreatedAt:        time.Now(),
		ModifiedByUserID: ctx.Value(constants.SessionUserID).(primitive.ObjectID),
		ModifiedByName:   ctx.Value(constants.SessionUserName).(string),
		ModifiedAt:       time.Now(),
	}

	// Add our comment to the comments.
	s.ModifiedByUserID = ctx.Value(constants.SessionUserID).(primitive.ObjectID)
	s.ModifiedAt = time.Now()
	s.Comments = append(s.Comments, comment)

	// Save to the database the modified customer.
	if err := c.StoreStorer.UpdateByID(ctx, s); err != nil {
		c.Logger.Error("database update by id error", slog.Any("error", err))
		return nil, err
	}

	return s, nil
}
