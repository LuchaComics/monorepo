package controller

import (
	"context"
	"log/slog"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	s_d "github.com/LuchaComics/monorepo/cloud/cps-backend/app/userpurchase/datastore"
	"github.com/LuchaComics/monorepo/cloud/cps-backend/utils/httperror"
)

func (impl *UserPurchaseControllerImpl) DeleteByID(ctx context.Context, id primitive.ObjectID) error {
	d, err := impl.GetByID(ctx, id)
	if err != nil {
		impl.Logger.Error("database get by id error", slog.Any("error", err))
		return err
	}
	if d == nil {
		impl.Logger.Error("database returns nothing from get by id")
		return httperror.NewForBadRequestWithSingleField("id", "userpurchase does not exist")
	}
	d.Status = s_d.StatusArchived
	d.ModifiedAt = time.Now()

	// Save to the database the modified store.
	if err := impl.UserPurchaseStorer.UpdateByID(ctx, d); err != nil {
		impl.Logger.Error("database update by id error", slog.Any("error", err))
		return err
	}

	return nil
}
