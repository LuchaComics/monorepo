package controller

import (
	"context"
	"log/slog"

	"go.mongodb.org/mongo-driver/bson/primitive"

	domain "github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/app/nftmetadata/datastore"
	"github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/utils/httperror"
)

func (c *NFTMetadataControllerImpl) ArchiveByID(ctx context.Context, id primitive.ObjectID) error {
	// Retrieve from our database the record for the specific id.
	m, err := c.NFTMetadataStorer.GetByID(ctx, id)
	if err != nil {
		c.Logger.Error("database get by id error", slog.Any("error", err))
		return err
	}
	if m == nil {
		return httperror.NewForBadRequestWithSingleField("id", "nftmetadata does not exist")
	}

	m.Status = domain.StatusArchived
	if err := c.NFTMetadataStorer.UpdateByID(ctx, m); err != nil {
		c.Logger.Error("database update by id error", slog.Any("error", err))
		return err
	}

	return nil
}
