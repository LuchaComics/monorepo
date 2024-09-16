package datastore

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (impl NFTStorerImpl) ListByNFTCollectionID(ctx context.Context, nftCollectionID primitive.ObjectID) (*NFTPaginationListResult, error) {
	f := &NFTPaginationListFilter{
		CollectionID:    nftCollectionID,
		Cursor:          "",
		PageSize:        1_000_000_000,
		SortField:       "created_at",
		SortOrder:       -1, // 1=ascending | -1=descending
		ExcludeArchived: true,
	}
	return impl.ListByFilter(ctx, f)
}
