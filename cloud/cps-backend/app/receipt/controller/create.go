package controller

import (
	"context"
	"log/slog"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	s_d "github.com/LuchaComics/monorepo/cloud/cps-backend/app/receipt/datastore"
	u_d "github.com/LuchaComics/monorepo/cloud/cps-backend/app/user/datastore"
	"github.com/LuchaComics/monorepo/cloud/cps-backend/config/constants"
	"github.com/LuchaComics/monorepo/cloud/cps-backend/utils/httperror"
)

func (c *ReceiptControllerImpl) Create(ctx context.Context, m *s_d.Receipt) (*s_d.Receipt, error) {
	// Extract from our session the following data.
	urole := ctx.Value(constants.SessionUserRole).(int8)
	// uid := ctx.Value(constants.SessionUserID).(primitive.ObjectID)
	// uname := ctx.Value(constants.SessionUserName).(string)
	oid := ctx.Value(constants.SessionUserStoreID).(primitive.ObjectID)
	oname := ctx.Value(constants.SessionUserStoreName).(string)

	switch urole { // Security.
	case u_d.UserRoleRoot:
		c.Logger.Debug("access granted")
	default:
		return nil, httperror.NewForForbiddenWithSingleField("message", "you do not have permission")
	}

	// Add defaults.
	m.StoreID = oid
	m.StoreName = oname
	m.ID = primitive.NewObjectID()
	m.CreatedAt = time.Now()
	// m.CreatedByUserID = uid
	// m.CreatedByUserName = uname
	m.ModifiedAt = time.Now()
	// m.ModifiedByUserID = uid
	// m.ModifiedByUserName = uname
	m.Status = s_d.StatusActive

	// Save to our database.
	if err := c.ReceiptStorer.Create(ctx, m); err != nil {
		c.Logger.Error("database create error", slog.Any("error", err))
		return nil, err
	}

	return m, nil
}
