package controller

import (
	"context"
	"log/slog"

	"go.mongodb.org/mongo-driver/bson/primitive"

	domain "github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/app/nftmetadata/datastore"
	"github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/utils/httperror"
)

func (c *NFTMetadataControllerImpl) GetByID(ctx context.Context, id primitive.ObjectID) (*domain.NFTMetadata, error) {
	// Retrieve from our database the record for the specific id.
	m, err := c.NFTMetadataStorer.GetByID(ctx, id)
	if err != nil {
		c.Logger.Error("database get by id error", slog.Any("error", err))
		return nil, err
	}
	if m == nil {
		return nil, httperror.NewForBadRequestWithSingleField("id", "nftmetadata does not exist")
	}
	return m, err
}
