package controller

import (
	"context"

	"log/slog"

	domain "github.com/LuchaComics/monorepo/cloud/cps-backend/app/receipt/datastore"
	"github.com/LuchaComics/monorepo/cloud/cps-backend/utils/httperror"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (c *ReceiptControllerImpl) GetByID(ctx context.Context, id primitive.ObjectID) (*domain.Receipt, error) {
	// Retrieve from our database the record for the specific id.
	m, err := c.ReceiptStorer.GetByID(ctx, id)
	if err != nil {
		c.Logger.Error("database get by id error", slog.Any("error", err))
		return nil, err
	}
	if m == nil {
		return nil, httperror.NewForBadRequestWithSingleField("id", "receipt does not exist")
	}
	return m, err
}
