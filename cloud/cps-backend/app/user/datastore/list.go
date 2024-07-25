package datastore

import (
	"context"
	"log/slog"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (impl UserStorerImpl) ListByFilter(ctx context.Context, f *UserPaginationListFilter) (*UserPaginationListResult, error) {
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
	if f.Role > 0 {
		filter["role"] = f.Role
	}
	if f.FirstName != "" {
		filter["first_name"] = f.FirstName
	}
	if f.LastName != "" {
		filter["last_name"] = f.LastName
	}
	if f.Email != "" {
		filter["email"] = f.Email
	}
	if f.Phone != "" {
		filter["phone"] = f.Phone
	}
	if f.Status != 0 {
		filter["status"] = f.Status
	}
	if !f.CreatedAtGTE.IsZero() {
		filter["created_at"] = bson.M{"$gt": f.CreatedAtGTE} // Add the cursor condition to the filter
	}
	switch f.IsStarred {
	case 1:
		filter["is_starred"] = true
		break
	case 2:
		filter["is_starred"] = false
		break
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
	results := []*User{}
	hasNextPage := false
	for cursor.Next(ctx) {
		document := &User{}
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

	return &UserPaginationListResult{
		Results:     results,
		NextCursor:  nextCursor,
		HasNextPage: hasNextPage,
	}, nil
}

func (impl UserStorerImpl) ListAllRootStaff(ctx context.Context) (*UserPaginationListResult, error) {
	f := &UserPaginationListFilter{
		Cursor:    "",
		PageSize:  1_000_000,
		SortField: "created_at",
		SortOrder: -1,
		Role:      UserRoleRoot,
		Status:    UserStatusActive,
	}
	return impl.ListByFilter(ctx, f)
}

func (impl UserStorerImpl) ListAllRetailerStaffForStoreID(ctx context.Context, storeID primitive.ObjectID) (*UserPaginationListResult, error) {
	f := &UserPaginationListFilter{
		Cursor:    "",
		PageSize:  1_000_000,
		SortField: "created_at",
		SortOrder: -1,
		Role:      UserRoleRetailer,
		StoreID:   storeID,
		Status:    UserStatusActive,
	}
	return impl.ListByFilter(ctx, f)
}
