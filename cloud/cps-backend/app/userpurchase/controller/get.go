package controller

import (
	"context"
	"log/slog"

	"go.mongodb.org/mongo-driver/bson/primitive"

	domain "github.com/LuchaComics/monorepo/cloud/cps-backend/app/userpurchase/datastore"
	"github.com/LuchaComics/monorepo/cloud/cps-backend/utils/httperror"
)

func (c *UserPurchaseControllerImpl) GetByID(ctx context.Context, id primitive.ObjectID) (*domain.UserPurchase, error) {
	// Retrieve from our database the record for the specific id.
	m, err := c.UserPurchaseStorer.GetByID(ctx, id)
	if err != nil {
		c.Logger.Error("database get by id error", slog.Any("error", err))
		return nil, err
	}
	if m == nil {
		return nil, httperror.NewForBadRequestWithSingleField("id", "userpurchase does not exist")
	}
	return m, err
}
