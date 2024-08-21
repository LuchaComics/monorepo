package controller

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"log/slog"

	org_d "github.com/LuchaComics/monorepo/cloud/cps-pinws-backend/app/tenant/datastore"
	"github.com/LuchaComics/monorepo/cloud/cps-pinws-backend/config/constants"
	"github.com/LuchaComics/monorepo/cloud/cps-pinws-backend/utils/httperror"
)

func (c *TenantControllerImpl) CreateComment(ctx context.Context, tenantID primitive.ObjectID, content string) (*org_d.Tenant, error) {
	// Fetch the original customer.
	s, err := c.TenantStorer.GetByID(ctx, tenantID)
	if err != nil {
		c.Logger.Error("database get by id error", slog.Any("error", err))
		return nil, err
	}
	if s == nil {
		c.Logger.Error("tenant does not exist error",
			slog.Any("tenant_id", tenantID))
		return nil, httperror.NewForBadRequestWithSingleField("message", "tenant does not exist")
	}

	// Create our comment.
	comment := &org_d.TenantComment{
		ID:               primitive.NewObjectID(),
		Content:          content,
		TenantID:   ctx.Value(constants.SessionUserTenantID).(primitive.ObjectID),
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
	if err := c.TenantStorer.UpdateByID(ctx, s); err != nil {
		c.Logger.Error("database update by id error", slog.Any("error", err))
		return nil, err
	}

	return s, nil
}
