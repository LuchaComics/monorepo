package datastore

import (
	"context"
	"log/slog"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (impl ComicSubmissionStorerImpl) CountAll(ctx context.Context) (int64, error) {

	opts := options.Count().SetHint("_id_")
	count, err := impl.Collection.CountDocuments(ctx, bson.D{}, opts)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (impl ComicSubmissionStorerImpl) CountByFilter(ctx context.Context, f *ComicSubmissionPaginationListFilter) (int64, error) {

	filter := bson.M{}

	// Add filter conditions to the filter. This code is similar to the list code.
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

	impl.Logger.Debug("listing filter:",
		slog.Any("filter", filter))

	opts := options.Count().SetHint("_id_")
	count, err := impl.Collection.CountDocuments(ctx, filter, opts)
	if err != nil {
		return 0, err
	}

	return count, nil
}
