package repo

import (
	"context"
	"errors"
	"log"
	"log/slog"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/config"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/domain"
)

type attachmentImpl struct {
	Logger     *slog.Logger
	DbClient   *mongo.Client
	Collection *mongo.Collection
}

func NewAttachmentRepository(appCfg *config.Configuration, loggerp *slog.Logger, client *mongo.Client) domain.AttachmentRepository {
	// ctx := context.Background()
	uc := client.Database(appCfg.DB.Name).Collection("attachments")

	// For debugging purposes only or if you are going to recreate new indexes.
	if _, err := uc.Indexes().DropAll(context.TODO()); err != nil {
		loggerp.Warn("failed deleting all indexes",
			slog.Any("err", err))

		// Do not crash app, just continue.
	}

	// Note:
	// * 1 for ascending
	// * -1 for descending
	// * "text" for text indexes

	// The following few lines of code will create the index for our app for this
	// colleciton.
	_, err := uc.Indexes().CreateMany(context.TODO(), []mongo.IndexModel{
		{Keys: bson.D{
			{Key: "tenant_id", Value: 1},
			{Key: "created_at", Value: -1},
		}},
		{Keys: bson.D{
			{Key: "tenant_id", Value: 1},
			{Key: "status", Value: 1},
			{Key: "created_at", Value: -1},
		}},
		{Keys: bson.D{
			{Key: "tenant_id", Value: 1},
			{Key: "user_id", Value: 1},
			{Key: "status", Value: 1},
			{Key: "created_at", Value: -1},
		}},
		{Keys: bson.D{
			{Key: "tenant_id", Value: 1},
			{Key: "user_id", Value: 1},
			{Key: "created_at", Value: -1},
		}},
		{Keys: bson.D{
			{Key: "tenant_id", Value: 1},
			{Key: "belong_type", Value: 1},
			{Key: "created_at", Value: -1},
		}},
		{Keys: bson.D{{Key: "object_key", Value: 1}}},
		{Keys: bson.D{
			{Key: "name", Value: "text"},
			{Key: "description", Value: "text"},
			{Key: "filename", Value: "text"},
		}},
	})
	if err != nil {
		// It is important that we crash the app on startup to meet the
		// requirements of `google/wire` framework.
		log.Fatal(err)
	}

	s := &attachmentImpl{
		Logger:     loggerp,
		DbClient:   client,
		Collection: uc,
	}
	return s
}

func (impl attachmentImpl) Create(ctx context.Context, u *domain.Attachment) error {
	// DEVELOPER NOTES:
	// According to mongodb documentaiton:
	//     Non-existent Databases and Collections
	//     If the necessary database and collection don't exist when you perform a write operation, the server implicitly creates them.
	//     Source: https://www.mongodb.com/docs/drivers/go/current/usage-examples/insertOne/

	if u.ID == primitive.NilObjectID {
		u.ID = primitive.NewObjectID()
		impl.Logger.Warn("database insert attachment not included id value, created id now.", slog.Any("id", u.ID))
	}

	_, err := impl.Collection.InsertOne(ctx, u)

	// check for errors in the insertion
	if err != nil {
		impl.Logger.Error("database failed create error",
			slog.Any("error", err))
		return err
	}

	return nil
}

func (impl attachmentImpl) GetByID(ctx context.Context, id primitive.ObjectID) (*domain.Attachment, error) {
	filter := bson.M{"_id": id}

	var result domain.Attachment
	err := impl.Collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// This error means your query did not match any documents.
			return nil, nil
		}
		impl.Logger.Error("database get by attachment id error", slog.Any("error", err))
		return nil, err
	}
	return &result, nil
}

func (impl attachmentImpl) UpdateByID(ctx context.Context, m *domain.Attachment) error {
	filter := bson.M{"_id": m.ID}

	update := bson.M{ // DEVELOPERS NOTE: https://stackoverflow.com/a/60946010
		"$set": m,
	}

	// execute the UpdateOne() function to update the first matching document
	_, err := impl.Collection.UpdateOne(ctx, filter, update)
	if err != nil {
		impl.Logger.Error("database update attachment by id error", slog.Any("error", err))
		return err
	}

	// // display the number of documents updated
	// impl.Logger.Debug("number of documents updated", slog.Int64("modified_count", result.ModifiedCount))

	return nil
}

// hasActiveFilters checks if any filters besides tenant_id are active
func (impl attachmentImpl) hasActiveFilters(filter *domain.AttachmentFilter) bool {
	return filter.Name != nil ||
		filter.Status != 0 ||
		filter.BelongsToType != 0 ||
		!filter.UserID.IsZero() ||
		filter.CreatedAtStart != nil ||
		filter.CreatedAtEnd != nil
}

// buildCountMatchStage creates the match stage for the aggregation pipeline
func (impl attachmentImpl) buildCountMatchStage(filter *domain.AttachmentFilter) bson.M {
	match := bson.M{
		"tenant_id": filter.TenantID,
	}

	if filter.Status != 0 {
		match["status"] = filter.Status
	}

	if filter.BelongsToType != 0 {
		match["belongs_to_type"] = filter.BelongsToType
	}

	if !filter.UserID.IsZero() {
		match["user_id"] = filter.UserID
	}

	// Date range filtering
	if filter.CreatedAtStart != nil || filter.CreatedAtEnd != nil {
		createdAtFilter := bson.M{}
		if filter.CreatedAtStart != nil {
			createdAtFilter["$gte"] = filter.CreatedAtStart
		}
		if filter.CreatedAtEnd != nil {
			createdAtFilter["$lte"] = filter.CreatedAtEnd
		}
		match["created_at"] = createdAtFilter
	}

	// Text search for name
	if filter.Name != nil && *filter.Name != "" {
		match["$text"] = bson.M{"$search": *filter.Name}
	}

	return match
}

func (impl attachmentImpl) CountByFilter(ctx context.Context, filter *domain.AttachmentFilter) (uint64, error) {
	if filter == nil {
		return 0, errors.New("filter cannot be nil")
	}

	// For exact counts with filters
	if impl.hasActiveFilters(filter) {
		pipeline := []bson.D{
			{{Key: "$match", Value: impl.buildCountMatchStage(filter)}},
			{{Key: "$count", Value: "total"}},
		}

		// Use aggregation with allowDiskUse for large datasets
		opts := options.Aggregate().SetAllowDiskUse(true)
		cursor, err := impl.Collection.Aggregate(ctx, pipeline, opts)
		if err != nil {
			return 0, err
		}
		defer cursor.Close(ctx)

		// Decode the result
		var results []struct {
			Total int64 `bson:"total"`
		}
		if err := cursor.All(ctx, &results); err != nil {
			return 0, err
		}

		if len(results) == 0 {
			return 0, nil
		}

		return uint64(results[0].Total), nil
	}

	// For unfiltered counts (much faster for basic tenant-only filtering)
	countOpts := options.Count().SetHint("tenant_id_1_created_at_-1")
	match := bson.M{"tenant_id": filter.TenantID}

	count, err := impl.Collection.CountDocuments(ctx, match, countOpts)
	if err != nil {
		return 0, err
	}

	return uint64(count), nil
}

func (impl attachmentImpl) buildMatchStage(filter *domain.AttachmentFilter) bson.M {
	match := bson.M{
		"tenant_id": filter.TenantID,
	}

	// Handle cursor-based pagination
	if filter.LastID != nil && filter.LastCreatedAt != nil {
		match["$or"] = []bson.M{
			{
				"created_at": bson.M{"$lt": filter.LastCreatedAt},
			},
			{
				"created_at": filter.LastCreatedAt,
				"_id":        bson.M{"$lt": filter.LastID},
			},
		}
	}

	// Add other filters
	if filter.Status != 0 {
		match["status"] = filter.Status
	}

	if filter.BelongsToType != 0 {
		match["belongs_to_type"] = filter.BelongsToType
	}

	if !filter.UserID.IsZero() {
		match["user_id"] = filter.UserID
	}

	if filter.CreatedAtStart != nil || filter.CreatedAtEnd != nil {
		createdAtFilter := bson.M{}
		if filter.CreatedAtStart != nil {
			createdAtFilter["$gte"] = filter.CreatedAtStart
		}
		if filter.CreatedAtEnd != nil {
			createdAtFilter["$lte"] = filter.CreatedAtEnd
		}
		match["created_at"] = createdAtFilter
	}

	// Text search for name
	if filter.Name != nil && *filter.Name != "" {
		match["$text"] = bson.M{"$search": *filter.Name}
	}

	return match
}

func (impl attachmentImpl) ListByFilter(ctx context.Context, filter *domain.AttachmentFilter) (*domain.AttachmentFilterResult, error) {
	if filter == nil {
		return nil, errors.New("filter cannot be nil")
	}

	// Default limit if not specified
	if filter.Limit <= 0 {
		filter.Limit = 100
	}

	// Request one more document than needed to determine if there are more results
	limit := filter.Limit + 1

	// Build the aggregation pipeline
	pipeline := make([]bson.D, 0)

	// Match stage - initial filtering
	matchStage := bson.D{{"$match", impl.buildMatchStage(filter)}}
	pipeline = append(pipeline, matchStage)

	// Sort stage
	sortStage := bson.D{{"$sort", bson.D{
		{"created_at", -1},
		{"_id", -1},
	}}}
	pipeline = append(pipeline, sortStage)

	// Limit stage
	limitStage := bson.D{{"$limit", limit}}
	pipeline = append(pipeline, limitStage)

	// Execute aggregation
	opts := options.Aggregate().SetAllowDiskUse(true)
	cursor, err := impl.Collection.Aggregate(ctx, pipeline, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	// Decode results
	var attachments []*domain.Attachment
	if err := cursor.All(ctx, &attachments); err != nil {
		return nil, err
	}

	// Handle empty results case
	if len(attachments) == 0 {
		// For debugging purposes only.
		// impl.Logger.Debug("Empty list",
		// 	slog.Any("filter", filter))
		return &domain.AttachmentFilterResult{
			Attachments: make([]*domain.Attachment, 0),
			HasMore:     false,
		}, nil
	}

	// Check if there are more results
	hasMore := false
	if len(attachments) > int(filter.Limit) {
		hasMore = true
		attachments = attachments[:len(attachments)-1]
	}

	// Get last document info for next page
	lastDoc := attachments[len(attachments)-1]

	return &domain.AttachmentFilterResult{
		Attachments:   attachments,
		HasMore:       hasMore,
		LastID:        lastDoc.ID,
		LastCreatedAt: lastDoc.CreatedAt,
	}, nil
}

func (impl attachmentImpl) DeleteByID(ctx context.Context, id primitive.ObjectID) error {
	_, err := impl.Collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		impl.Logger.Error("database failed deletion error",
			slog.Any("error", err))
		return err
	}
	return nil
}
