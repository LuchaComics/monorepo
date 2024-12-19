package repo

import (
	"context"
	"errors"
	"log"
	"log/slog"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/config"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/domain"
)

type comicSubmissionImplImpl struct {
	Logger     *slog.Logger
	DbClient   *mongo.Client
	Collection *mongo.Collection
}

func NewComicSubmissionRepository(appCfg *config.Configuration, loggerp *slog.Logger, client *mongo.Client) domain.ComicSubmissionRepository {
	// ctx := context.Background()
	uc := client.Database(appCfg.DB.Name).Collection("comic_submissions")

	// // For debugging purposes only or if you are going to recreate new indexes.
	// if _, err := uc.Indexes().DropAll(context.TODO()); err != nil {
	// 	loggerp.Warn("failed deleting all indexes",
	// 		slog.Any("err", err))
	//
	// 	// Do not crash app, just continue.
	// }

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
			{Key: "name", Value: "text"},
		}},
	})

	// db.comic_submissions.createIndex({ "tenant_id": 1, "created_at": -1 })
	// db.comic_submissions.createIndex({ "tenant_id": 1, "status": 1, "created_at": -1 })
	// db.comic_submissions.createIndex({ "tenant_id": 1, "user_id": 1, "created_at": -1 })
	// db.comic_submissions.createIndex({ "tenant_id": 1, "name": "text" })

	if err != nil {
		// It is important that we crash the app on startup to meet the
		// requirements of `google/wire` framework.
		log.Fatal(err)
	}

	s := &comicSubmissionImplImpl{
		Logger:     loggerp,
		DbClient:   client,
		Collection: uc,
	}
	return s
}

func (impl comicSubmissionImplImpl) Create(ctx context.Context, u *domain.ComicSubmission) error {
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

func (impl comicSubmissionImplImpl) GetByID(ctx context.Context, id primitive.ObjectID) (*domain.ComicSubmission, error) {
	filter := bson.M{"_id": id}

	var result domain.ComicSubmission
	err := impl.Collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// This error means your query did not match any documents.
			return nil, nil
		}
		impl.Logger.Error("database get by user id error", slog.Any("error", err))
		return nil, err
	}
	return &result, nil
}

func (impl comicSubmissionImplImpl) CountByUserID(ctx context.Context, userID primitive.ObjectID) (uint64, error) {
	filter := bson.M{
		"user_id": userID,
	}

	count, err := impl.Collection.CountDocuments(ctx, filter)
	if err != nil {
		impl.Logger.Error("database check if exists by ID error", slog.Any("error", err))
		return uint64(0), err
	}

	return uint64(count), nil
}

func (impl comicSubmissionImplImpl) CountByStatusAndByUserID(ctx context.Context, status int8, userID primitive.ObjectID) (uint64, error) {
	filter := bson.M{
		"user_id": userID,
		"status":  status,
	}

	count, err := impl.Collection.CountDocuments(ctx, filter)
	if err != nil {
		impl.Logger.Error("database check if exists by ID error", slog.Any("error", err))
		return uint64(0), err
	}

	return uint64(count), nil
}

func (impl comicSubmissionImplImpl) CountTotalCreatedTodayByUserID(ctx context.Context, userID primitive.ObjectID, timezone string) (uint64, error) {
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		impl.Logger.Warn("Failed validating",
			slog.Any("error", err))
		return 0, err
	}
	now := time.Now()
	userTime := now.In(loc)

	thisDayStart := time.Date(userTime.Year(), userTime.Month(), userTime.Day()-1, 0, 0, 0, 0, userTime.Location())
	thisDayEnd := time.Date(userTime.Year(), userTime.Month(), userTime.Day()+1, 0, 0, 0, 0, userTime.Location())

	///

	filter := bson.M{
		"user_id": userID,
	}

	var conditions []bson.M
	conditions = append(conditions, bson.M{"created_at": bson.M{"$gte": thisDayStart}})
	conditions = append(conditions, bson.M{"created_at": bson.M{"$lt": thisDayEnd}})
	filter["$and"] = conditions

	// impl.Logger.Debug("counting total created today",
	// 	slog.Any("thisDayStart", thisDayStart),
	// 	slog.Any("thisDayNow", time.Now()),
	// 	slog.Any("thisDayEnd", thisDayEnd),
	// 	slog.Any("filter", filter))

	count, err := impl.Collection.CountDocuments(ctx, filter)
	if err != nil {
		impl.Logger.Error("database check if exists by ID error", slog.Any("error", err))
		return uint64(0), err
	}

	// impl.Logger.Debug("finished counting total created today",
	// 	slog.Any("count", count))

	return uint64(count), nil
}

func (impl comicSubmissionImplImpl) CountCoinsRewardByUserID(ctx context.Context, userID primitive.ObjectID) (uint64, error) {
	// Define the aggregation pipeline
	pipeline := mongo.Pipeline{
		// Match documents with the given user_id
		{{"$match", bson.D{{"user_id", userID}}}},
		// Group by user_id and calculate the total coins_reward
		{{"$group", bson.D{
			{"_id", nil}, // No grouping key, we just want the total
			{"totalCoins", bson.D{{"$sum", "$coins_reward"}}},
		}}},
	}

	// Execute the aggregation
	cursor, err := impl.Collection.Aggregate(ctx, pipeline)
	if err != nil {
		return 0, err
	}
	defer cursor.Close(ctx)

	// Parse the result
	var result []struct {
		TotalCoins uint64 `bson:"totalCoins"`
	}
	if err := cursor.All(ctx, &result); err != nil {
		return 0, err
	}

	// Return the total coins if found, otherwise 0
	if len(result) > 0 {
		return result[0].TotalCoins, nil
	}

	return 0, nil
}

func (impl comicSubmissionImplImpl) CountCoinsRewardByStatusAndByUserID(ctx context.Context, status int8, userID primitive.ObjectID) (uint64, error) {
	// Define the aggregation pipeline
	pipeline := mongo.Pipeline{
		// Match documents with the given user_id
		{{"$match", bson.D{
			{"user_id", userID},
			{"status", status},
		}}},
		// Group by user_id and calculate the total coins_reward
		{{"$group", bson.D{
			{"_id", nil}, // No grouping key, we just want the total
			{"totalCoins", bson.D{{"$sum", "$coins_reward"}}},
		}}},
	}

	// Execute the aggregation
	cursor, err := impl.Collection.Aggregate(ctx, pipeline)
	if err != nil {
		return 0, err
	}
	defer cursor.Close(ctx)

	// Parse the result
	var result []struct {
		TotalCoins uint64 `bson:"totalCoins"`
	}
	if err := cursor.All(ctx, &result); err != nil {
		return 0, err
	}

	// Return the total coins if found, otherwise 0
	if len(result) > 0 {
		return result[0].TotalCoins, nil
	}

	return 0, nil
}

// hasActiveFilters checks if any filters besides tenant_id are active
func hasActiveFilters(filter *domain.ComicSubmissionFilter) bool {
	return filter.Name != nil ||
		filter.Status != 0 ||
		!filter.UserID.IsZero() ||
		filter.CreatedAtStart != nil ||
		filter.CreatedAtEnd != nil
}

// buildCountMatchStage creates the match stage for the aggregation pipeline
func buildCountMatchStage(filter *domain.ComicSubmissionFilter) bson.M {
	match := bson.M{
		"tenant_id": filter.TenantID,
	}

	if filter.Status != 0 {
		match["status"] = filter.Status
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

func (s *comicSubmissionImplImpl) CountByFilter(ctx context.Context, filter *domain.ComicSubmissionFilter) (uint64, error) {
	if filter == nil {
		return 0, errors.New("filter cannot be nil")
	}

	// For exact counts with filters
	if hasActiveFilters(filter) {
		pipeline := []bson.D{
			{{Key: "$match", Value: buildCountMatchStage(filter)}},
			{{Key: "$count", Value: "total"}},
		}

		// Use aggregation with allowDiskUse for large datasets
		opts := options.Aggregate().SetAllowDiskUse(true)
		cursor, err := s.Collection.Aggregate(ctx, pipeline, opts)
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

	count, err := s.Collection.CountDocuments(ctx, match, countOpts)
	if err != nil {
		return 0, err
	}

	return uint64(count), nil
}

func (s *comicSubmissionImplImpl) CountCoinsRewardByFilter(ctx context.Context, filter *domain.ComicSubmissionFilter) (uint64, error) {
	if filter == nil {
		return 0, errors.New("filter cannot be nil")
	}

	// Build aggregation pipeline
	pipeline := []bson.D{
		// Match stage
		{{Key: "$match", Value: buildMatchStage(filter)}},

		// Group stage to sum coins
		{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: nil},
			{Key: "totalCoins", Value: bson.D{
				{Key: "$sum", Value: "$coins_reward"},
			}},
		}}},
	}

	// Use aggregation with allowDiskUse for large datasets
	opts := options.Aggregate().SetAllowDiskUse(true)
	cursor, err := s.Collection.Aggregate(ctx, pipeline, opts)
	if err != nil {
		return 0, err
	}
	defer cursor.Close(ctx)

	// Decode the result
	var results []struct {
		TotalCoins int64 `bson:"totalCoins"`
	}
	if err := cursor.All(ctx, &results); err != nil {
		return 0, err
	}

	// If no results found, return 0
	if len(results) == 0 {
		return 0, nil
	}

	return uint64(results[0].TotalCoins), nil
}

func buildMatchStage(filter *domain.ComicSubmissionFilter) bson.M {
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

func (s *comicSubmissionImplImpl) ListByFilter(ctx context.Context, filter *domain.ComicSubmissionFilter) (*domain.ComicSubmissionFilterResult, error) {
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
	matchStage := bson.D{{"$match", buildMatchStage(filter)}}
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
	cursor, err := s.Collection.Aggregate(ctx, pipeline, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	// Decode results
	var submissions []*domain.ComicSubmission
	if err := cursor.All(ctx, &submissions); err != nil {
		return nil, err
	}

	// Handle empty results case
	if len(submissions) == 0 {
		s.Logger.Debug("Empty list", slog.Any("filter", filter))
		return &domain.ComicSubmissionFilterResult{
			Submissions: make([]*domain.ComicSubmission, 0),
			HasMore:     false,
		}, nil
	}

	// Check if there are more results
	hasMore := false
	if len(submissions) > int(filter.Limit) {
		hasMore = true
		submissions = submissions[:len(submissions)-1]
	}

	// Get last document info for next page
	lastDoc := submissions[len(submissions)-1]

	return &domain.ComicSubmissionFilterResult{
		Submissions:   submissions,
		HasMore:       hasMore,
		LastID:        lastDoc.ID,
		LastCreatedAt: lastDoc.CreatedAt,
	}, nil
}

func (impl comicSubmissionImplImpl) TotalCoinsAwarded(ctx context.Context) (uint64, error) {
	// Define the aggregation pipeline
	pipeline := mongo.Pipeline{
		// Match documents with the given criteria.
		{{"$match", bson.D{{"was_awarded", true}}}},
		// Group by user_id and calculate the total coins_reward
		{{"$group", bson.D{
			{"_id", nil}, // No grouping key, we just want the total
			{"totalCoins", bson.D{{"$sum", "$coins_reward"}}},
		}}},
	}

	// Execute the aggregation
	cursor, err := impl.Collection.Aggregate(ctx, pipeline)
	if err != nil {
		return 0, err
	}
	defer cursor.Close(ctx)

	// Parse the result
	var result []struct {
		TotalCoins uint64 `bson:"totalCoins"`
	}
	if err := cursor.All(ctx, &result); err != nil {
		return 0, err
	}

	// Return the total coins if found, otherwise 0
	if len(result) > 0 {
		return result[0].TotalCoins, nil
	}

	return 0, nil
}

func (impl comicSubmissionImplImpl) TotalCoinsAwardedByUserID(ctx context.Context, userID primitive.ObjectID) (uint64, error) {
	filter := &domain.ComicSubmissionFilter{
		UserID: userID,
	}

	// Define the aggregation pipeline
	pipeline := mongo.Pipeline{
		// Match documents with the given criteria.
		{{"$match", buildMatchStage(filter)}},
		// Group by user_id and calculate the total coins_reward
		{{"$group", bson.D{
			{"_id", nil}, // No grouping key, we just want the total
			{"totalCoins", bson.D{{"$sum", "$coins_reward"}}},
		}}},
	}

	// Execute the aggregation
	cursor, err := impl.Collection.Aggregate(ctx, pipeline)
	if err != nil {
		return 0, err
	}
	defer cursor.Close(ctx)

	// Parse the result
	var result []struct {
		TotalCoins uint64 `bson:"totalCoins"`
	}
	if err := cursor.All(ctx, &result); err != nil {
		return 0, err
	}

	// Return the total coins if found, otherwise 0
	if len(result) > 0 {
		return result[0].TotalCoins, nil
	}

	return 0, nil
}

//	func (impl comicSubmissionImplImpl) GetByEmail(ctx context.Context, email string) (*domain.ComicSubmission, error) {
//		filter := bson.M{"email": email}
//
//		var result domain.ComicSubmission
//		err := impl.Collection.FindOne(ctx, filter).Decode(&result)
//		if err != nil {
//			if err == mongo.ErrNoDocuments {
//				// This error means your query did not match any documents.
//				return nil, nil
//			}
//			impl.Logger.Error("database get by email error", slog.Any("error", err))
//			return nil, err
//		}
//		return &result, nil
//	}
//
//	func (impl comicSubmissionImplImpl) GetByVerificationCode(ctx context.Context, verificationCode string) (*domain.ComicSubmission, error) {
//		filter := bson.M{"email_verification_code": verificationCode}
//
//		var result domain.ComicSubmission
//		err := impl.Collection.FindOne(ctx, filter).Decode(&result)
//		if err != nil {
//			if err == mongo.ErrNoDocuments {
//				// This error means your query did not match any documents.
//				return nil, nil
//			}
//			impl.Logger.Error("database get by verification code error", slog.Any("error", err))
//			return nil, err
//		}
//		return &result, nil
//	}
//
//	func (impl comicSubmissionImplImpl) DeleteByID(ctx context.Context, id primitive.ObjectID) error {
//		_, err := impl.Collection.DeleteOne(ctx, bson.M{"_id": id})
//		if err != nil {
//			impl.Logger.Error("database failed deletion error",
//				slog.Any("error", err))
//			return err
//		}
//		return nil
//	}
//
//	func (impl comicSubmissionImplImpl) CheckIfExistsByID(ctx context.Context, id primitive.ObjectID) (bool, error) {
//		filter := bson.M{"_id": id}
//		count, err := impl.Collection.CountDocuments(ctx, filter)
//		if err != nil {
//			impl.Logger.Error("database check if exists by ID error", slog.Any("error", err))
//			return false, err
//		}
//		return count >= 1, nil
//	}
//
//	func (impl comicSubmissionImplImpl) CheckIfExistsByEmail(ctx context.Context, email string) (bool, error) {
//		filter := bson.M{"email": email}
//		count, err := impl.Collection.CountDocuments(ctx, filter)
//		if err != nil {
//			impl.Logger.Error("database check if exists by email error", slog.Any("error", err))
//			return false, err
//		}
//		return count >= 1, nil
//	}
func (impl comicSubmissionImplImpl) UpdateByID(ctx context.Context, m *domain.ComicSubmission) error {
	filter := bson.M{"_id": m.ID}

	update := bson.M{ // DEVELOPERS NOTE: https://stackoverflow.com/a/60946010
		"$set": m,
	}

	// execute the UpdateOne() function to update the first matching document
	_, err := impl.Collection.UpdateOne(ctx, filter, update)
	if err != nil {
		impl.Logger.Error("database update by id error", slog.Any("error", err))
		return err
	}

	// // display the number of documents updated
	// impl.Logger.Debug("number of documents updated", slog.Int64("modified_count", result.ModifiedCount))

	return nil
}
