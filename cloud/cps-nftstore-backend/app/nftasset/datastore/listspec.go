package datastore

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (impl NFTAssetStorerImpl) ListByNFTMetadataID(ctx context.Context, nftmetadataID primitive.ObjectID) (*NFTAssetPaginationListResult, error) {
	f := &NFTAssetPaginationListFilter{
		PageSize:      1_000_000_000, // Essentially unlimited
		SortField:     "created",
		SortOrder:     SortOrderDescending,
		NFTMetadataID: nftmetadataID,
	}
	return impl.ListByFilter(ctx, f)
}

func (impl NFTAssetStorerImpl) ListByNFTCollectionID(ctx context.Context, nftCollectionID primitive.ObjectID) (*NFTAssetPaginationListResult, error) {
	f := &NFTAssetPaginationListFilter{
		PageSize:        1_000_000_000, // Essentially unlimited
		SortField:       "created",
		SortOrder:       SortOrderDescending,
		NFTCollectionID: nftCollectionID,
	}
	return impl.ListByFilter(ctx, f)
}
