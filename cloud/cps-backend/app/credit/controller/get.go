package controller

import (
	"context"

	"log/slog"

	domain "github.com/LuchaComics/monorepo/cloud/cps-backend/app/credit/datastore"
	"github.com/LuchaComics/monorepo/cloud/cps-backend/utils/httperror"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (c *CreditControllerImpl) GetByID(ctx context.Context, id primitive.ObjectID) (*domain.Credit, error) {
	// Retrieve from our database the record for the specific id.
	m, err := c.CreditStorer.GetByID(ctx, id)
	if err != nil {
		c.Logger.Error("database get by id error", slog.Any("error", err))
		return nil, err
	}
	if m == nil {
		return nil, httperror.NewForBadRequestWithSingleField("id", "credit does not exist")
	}
	return m, err
}
