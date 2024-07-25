package datastore

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (impl CreditStorerImpl) GetNextAvailable(ctx context.Context, userID primitive.ObjectID, serviceType int8) (*Credit, error) {

	f := &CreditPaginationListFilter{
		// Pagination related.
		Cursor:    "",
		PageSize:  1_000,
		SortField: "created_at",
		SortOrder: 1, // 1=ascending

		// Filter related.
		UserID:           userID,
		OfferServiceType: serviceType,
		Status:           StatusActive,
		BusinessFunction: BusinessFunctionGrantFreeSubmission,
	}
	res, err := impl.ListByFilter(ctx, f)
	if err != nil {
		return nil, err
	}
	for _, credit := range res.Results {
		if credit != nil {
			return credit, nil
		}
	}

	return nil, nil
}
