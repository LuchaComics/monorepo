package controller

import (
	"context"
	"log/slog"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	u_d "github.com/LuchaComics/monorepo/cloud/cps-backend/app/user/datastore"
	domain "github.com/LuchaComics/monorepo/cloud/cps-backend/app/userpurchase/datastore"
	"github.com/LuchaComics/monorepo/cloud/cps-backend/config/constants"
	"github.com/LuchaComics/monorepo/cloud/cps-backend/utils/httperror"
)

func (c *UserPurchaseControllerImpl) UpdateByID(ctx context.Context, ns *domain.UserPurchase) (*domain.UserPurchase, error) {
	// Extract from our session the following data.
	urole, _ := ctx.Value(constants.SessionUserRole).(int8)
	// uid := ctx.Value(constants.SessionUserID).(primitive.ObjectID)
	// uname := ctx.Value(constants.SessionUserName).(string)
	oid, _ := ctx.Value(constants.SessionUserStoreID).(primitive.ObjectID)
	oname, _ := ctx.Value(constants.SessionUserStoreName).(string)
	otz, _ := ctx.Value(constants.SessionUserStoreTimezone).(string)

	switch urole { // Security.
	case u_d.UserRoleRoot:
		c.Logger.Debug("access granted")
	default:
		return nil, httperror.NewForForbiddenWithSingleField("message", "you do not have permission")
	}

	// Fetch the original store.
	os, err := c.UserPurchaseStorer.GetByID(ctx, ns.ID)
	if err != nil {
		c.Logger.Error("database get by id error", slog.Any("error", err))
		return nil, err
	}
	if os == nil {
		return nil, httperror.NewForBadRequestWithSingleField("id", "userpurchase type does not exist")
	}

	// Modify our original store.
	os.StoreID = oid
	os.StoreName = oname
	os.StoreTimezone = otz
	os.ModifiedAt = time.Now()
	os.Status = ns.Status
	// os.Name = ns.Name

	// Save to the database the modified store.
	if err := c.UserPurchaseStorer.UpdateByID(ctx, os); err != nil {
		c.Logger.Error("database update by id error", slog.Any("error", err))
		return nil, err
	}

	return os, nil
}
