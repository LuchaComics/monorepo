package controller

import (
	"context"
	"log/slog"

	"go.mongodb.org/mongo-driver/bson/primitive"

	nft_s "github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/app/nft/datastore"
	"github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/utils/httperror"
)

func (c *NFTControllerImpl) GetByID(ctx context.Context, id primitive.ObjectID) (*nft_s.NFT, error) {
	// Retrieve from our database the record for the specific id.
	m, err := c.NFTStorer.GetByID(ctx, id)
	if err != nil {
		c.Logger.Error("database get by id error", slog.Any("error", err))
		return nil, err
	}
	if m == nil {
		return nil, httperror.NewForBadRequestWithSingleField("id", "nft does not exist")
	}
	return m, err
}

func (c *NFTControllerImpl) GetByTokenID(ctx context.Context, tokenID uint64) (*nft_s.NFT, error) {
	// Retrieve from our database the record for the specific id.
	m, err := c.NFTStorer.GetByTokenID(ctx, tokenID)
	if err != nil {
		c.Logger.Error("database get by token_id error", slog.Any("error", err))
		return nil, err
	}
	if m == nil {
		return nil, httperror.NewForBadRequestWithSingleField("id", "nft does not exist")
	}
	return m, err
}
