package datastore

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (impl PinObjectStorerImpl) ListByProjectID(ctx context.Context, projectID primitive.ObjectID) (*PinObjectPaginationListResult, error) {
	f := &PinObjectPaginationListFilter{
		PageSize:  1_000_000_000, // Essentially unlimited
		SortField: "created",
		SortOrder: SortOrderDescending,
		ProjectID: projectID,
	}
	return impl.ListByFilter(ctx, f)
}
