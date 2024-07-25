package controller

import (
	"context"

	"log/slog"

	domain "github.com/LuchaComics/monorepo/cloud/cps-backend/app/offer/datastore"
	"github.com/LuchaComics/monorepo/cloud/cps-backend/utils/httperror"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (c *OfferControllerImpl) GetByID(ctx context.Context, id primitive.ObjectID) (*domain.Offer, error) {
	// Retrieve from our database the record for the specific id.
	m, err := c.OfferStorer.GetByID(ctx, id)
	if err != nil {
		c.Logger.Error("database get by id error", slog.Any("error", err))
		return nil, err
	}
	if m == nil {
		return nil, httperror.NewForBadRequestWithSingleField("id", "offer does not exist")
	}
	return m, err
}

func (c *OfferControllerImpl) GetByServiceType(ctx context.Context, serviceType int8) (*domain.Offer, error) {
	// Retrieve from our database the record for the specific id.
	m, err := c.OfferStorer.GetByServiceType(ctx, serviceType)
	if err != nil {
		c.Logger.Error("database get by service type error", slog.Any("error", err))
		return nil, err
	}
	if m == nil {
		return nil, httperror.NewForBadRequestWithSingleField("service_type", "offer does not exist")
	}
	return m, err
}
