package datastore

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (impl NFTMetadataStorerImpl) ListByNFTCollection(ctx context.Context, nftCollectionID primitive.ObjectID) (*NFTMetadataPaginationListResult, error) {
	f := &NFTMetadataPaginationListFilter{
		CollectionID:    nftCollectionID,
		Cursor:          "",
		PageSize:        1_000_000_000,
		SortField:       "created_at",
		SortOrder:       -1, // 1=ascending | -1=descending
		ExcludeArchived: true,
	}
	return impl.ListByFilter(ctx, f)
}
