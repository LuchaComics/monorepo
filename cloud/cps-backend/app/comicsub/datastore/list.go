package datastore

import (
	"context"
	"time"

	"log/slog"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (impl ComicSubmissionStorerImpl) ListByFilter(ctx context.Context, f *ComicSubmissionPaginationListFilter) (*ComicSubmissionPaginationListResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 12*time.Second)
	defer cancel()

	filter, err := impl.newPaginationFilter(f)
	if err != nil {
		return nil, err
	}

	// Add filter conditions to the filter
	if f.UserID != primitive.NilObjectID {
		filter["user_id"] = f.UserID
	}
	if f.UserEmail != "" {
		filter["user.email"] = f.UserEmail
	}
	if f.CreatedByUserRole != 0 {
		filter["created_by_user_role"] = f.CreatedByUserRole
	}
	if f.StoreID != primitive.NilObjectID {
		filter["store_id"] = f.StoreID
	}
	if f.StoreSpecialCollection != 0 {
		filter["store_special_colleciton"] = f.StoreSpecialCollection
	}
	if f.Status != 0 {
		filter["status"] = f.Status
	}
	if !f.CreatedAtGTE.IsZero() {
		filter["created_at"] = bson.M{"$gt": f.CreatedAtGTE} // Add the cursor condition to the filter
	}
	if f.ServiceType != 0 {
		filter["service_type"] = f.ServiceType
	}

	if len(f.ExcludeStoreSpecialCollections) > 0 {
		filter["store_special_colleciton"] = bson.M{"$nin": f.ExcludeStoreSpecialCollections}
	}
	if f.CPSRNClassification != "" {
		filter["cpsrn_classification"] = f.CPSRNClassification
	}
	if !f.InspectorID.IsZero() {
		filter["inspector_id"] = f.InspectorID
	}
	if !f.CustomerID.IsZero() {
		filter["customer_id"] = f.CustomerID
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
	results := []*ComicSubmission{}
	hasNextPage := false
	for cursor.Next(ctx) {
		document := &ComicSubmission{}
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

	return &ComicSubmissionPaginationListResult{
		Results:     results,
		NextCursor:  nextCursor,
		HasNextPage: hasNextPage,
	}, nil
}
