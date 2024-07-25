package datastore

import (
	"context"
	"log/slog"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

func (impl AttachmentStorerImpl) ListByFilter(ctx context.Context, f *AttachmentPaginationListFilter) (*AttachmentPaginationListResult, error) {
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
	if !f.OwnershipID.IsZero() {
		filter["ownership_id"] = f.OwnershipID
	}
	if !f.CreatedByUserID.IsZero() {
		filter["created_by_user_id"] = f.CreatedByUserID
	}
	if !f.ModifiedByUserID.IsZero() {
		filter["modified_by_user_id"] = f.ModifiedByUserID
	}
	if f.ExcludeArchived {
		filter["status"] = bson.M{"$ne": StatusArchived} // Do not list archived items! This code
	}

	impl.Logger.Debug("fetching attachments list",
		slog.Any("Cursor", f.Cursor),
		slog.Int64("PageSize", f.PageSize),
		slog.String("SortField", f.SortField),
		slog.Any("SortOrder", f.SortOrder),
		slog.Any("StoreID", f.StoreID),
		slog.Any("OwnershipID", f.OwnershipID),
		slog.Any("ExcludeArchived", f.ExcludeArchived),
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
	results := []*Attachment{}
	hasNextPage := false
	for cursor.Next(ctx) {
		document := &Attachment{}
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

	return &AttachmentPaginationListResult{
		Results:     results,
		NextCursor:  nextCursor,
		HasNextPage: hasNextPage,
	}, nil
}
