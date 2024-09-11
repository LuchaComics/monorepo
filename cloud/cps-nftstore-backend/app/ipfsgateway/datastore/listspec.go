package datastore

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (impl PinObjectStorerImpl) ListByTenantID(ctx context.Context, tenantID primitive.ObjectID) (*PinObjectPaginationListResult, error) {
	f := &PinObjectPaginationListFilter{
		PageSize:  1_000_000_000, // Essentially unlimited
		SortField: "created",
		SortOrder: SortOrderDescending,
		TenantID:  tenantID,
	}
	return impl.ListByFilter(ctx, f)
}
