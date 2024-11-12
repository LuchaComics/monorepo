package datastore

import (
	"context"
	"log/slog"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

func (impl NFTAssetStorerImpl) ListByFilter(ctx context.Context, f *NFTAssetPaginationListFilter) (*NFTAssetPaginationListResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 12*time.Second)
	defer cancel()

	filter, err := impl.newPaginationFilter(f)
	if err != nil {
		return nil, err
	}

	// Add filter conditions to the filter
	if !f.TenantID.IsZero() {
		filter["tenant_id"] = f.TenantID
	}
	if !f.NFTMetadataID.IsZero() {
		filter["nftmetadata_id"] = f.NFTMetadataID
	}
	if !f.NFTCollectionID.IsZero() {
		filter["nftcollection_id"] = f.NFTCollectionID
	}

	// Create a slice to store conditions
	var conditions []bson.M

	// Add filter conditions to the slice
	if !f.CreatedAtGTE.IsZero() {
		conditions = append(conditions, bson.M{"created_at": bson.M{"$gte": f.CreatedAtGTE}})
	}
	if !f.CreatedAtGT.IsZero() {
		conditions = append(conditions, bson.M{"created_at": bson.M{"$gt": f.CreatedAtGT}})
	}
	if !f.CreatedAtLTE.IsZero() {
		conditions = append(conditions, bson.M{"created_at": bson.M{"$lte": f.CreatedAtLTE}})
	}
	if !f.CreatedAtLT.IsZero() {
		conditions = append(conditions, bson.M{"created_at": bson.M{"$lt": f.CreatedAtLT}})
	}

	// Combine conditions with $and operator
	if len(conditions) > 0 {
		filter["$and"] = conditions
	}

	impl.Logger.Debug("fetching nftassets list",
		slog.Any("Cursor", f.Cursor),
		slog.Int64("PageSize", f.PageSize),
		slog.String("SortField", f.SortField),
		slog.Any("SortOrder", f.SortOrder),
		slog.Any("TenantID", f.TenantID),
		slog.Any("NFTMetadataID", f.NFTMetadataID),
	)

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
	results := []*NFTAsset{}
	hasNextPage := false
	for cursor.Next(ctx) {
		document := &NFTAsset{}
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

	return &NFTAssetPaginationListResult{
		Results:     results,
		NextCursor:  nextCursor,
		HasNextPage: hasNextPage,
	}, nil
}