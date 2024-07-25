package datastore

import (
	"context"
	"log/slog"
	"time"
)

func (impl CreditStorerImpl) ListByFilter(ctx context.Context, f *CreditPaginationListFilter) (*CreditPaginationListResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 12*time.Second)
	defer cancel()

	filter, err := impl.newPaginationFilter(f)
	if err != nil {
		return nil, err
	}

	// Add filter conditions to the filter
	if !f.StoreID.IsZero() {
		filter["store_id"] = f.StoreID
	}
	if !f.UserID.IsZero() {
		filter["user_id"] = f.UserID
	}
	if !f.OfferID.IsZero() {
		filter["offer_id"] = f.OfferID
	}
	if f.OfferServiceType != 0 {
		filter["offer_service_type"] = f.OfferServiceType
	}
	if f.Status != 0 {
		filter["status"] = f.Status
	}
	if f.BusinessFunction != 0 {
		filter["business_function"] = f.BusinessFunction
	}

	impl.Logger.Debug("listing filter:",
		slog.Any("filter", filter))

	// Include additional filters for our cursor-based pagination pertaining to sorting and limit.
	options, err := impl.newPaginationOptions(f)
	if err != nil {
		return nil, err
	}

	// Execute the query
	cursor, err := impl.Collection.Find(ctx, filter, options)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	// Retrieve the documents and check if there is a next page
	results := []*Credit{}
	hasNextPage := false
	for cursor.Next(ctx) {
		document := &Credit{}
		if err := cursor.Decode(document); err != nil {
			return nil, err
		}
		results = append(results, document)
		// Stop fetching documents if we have reached the desired page size
		if int64(len(results)) >= f.PageSize {
			hasNextPage = true
			break
		}
	}

	// Get the next cursor and encode it
	var nextCursor string
	if hasNextPage {
		nextCursor, err = impl.newPaginatorNextCursor(f, results)
		if err != nil {
			return nil, err
		}
	}

	return &CreditPaginationListResult{
		Results:     results,
		NextCursor:  nextCursor,
		HasNextPage: hasNextPage,
	}, nil
}
