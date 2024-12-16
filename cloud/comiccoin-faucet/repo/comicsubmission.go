package repo

import (
	"context"
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

func (s *comicSubmissionImplImpl) ListByFilter(ctx context.Context, filter *domain.ComicSubmissionFilter) (*domain.ComicSubmissionFilterResult, error) {
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

	// Check if there are more results
	hasMore := false
	if len(submissions) > int(filter.Limit) {
		hasMore = true
		submissions = submissions[:len(submissions)-1]
	}

	// Get last document info for next page
	var lastID primitive.ObjectID
	var lastCreatedAt time.Time
	if len(submissions) > 0 {
		lastDoc := submissions[len(submissions)-1]
		lastID = lastDoc.ID
		lastCreatedAt = lastDoc.CreatedAt
	}

	return &domain.ComicSubmissionFilterResult{
		Submissions:   submissions,
		HasMore:       hasMore,
		LastID:        lastID,
		LastCreatedAt: lastCreatedAt,
	}, nil
}

func buildMatchStage(filter *domain.ComicSubmissionFilter) bson.M {
	match := bson.M{
		"tenant_id": filter.TenantID, // Always include tenant_id for data partitioning
	}

	// Build cursor-based pagination condition
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
	if filter.Status != nil {
		match["status"] = *filter.Status
	}

	if filter.Type != nil {
		match["type"] = *filter.Type
	}

	if !filter.UserID.IsZero() {
		match["user_id"] = filter.UserID
	}

	// Date range
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

	// Text search (if name is provided)
	if filter.Name != nil && *filter.Name != "" {
		match["$text"] = bson.M{"$search": *filter.Name}
	}

	return match
}

// func (impl comicSubmissionImplImpl) GetByEmail(ctx context.Context, email string) (*domain.ComicSubmission, error) {
// 	filter := bson.M{"email": email}
//
// 	var result domain.ComicSubmission
// 	err := impl.Collection.FindOne(ctx, filter).Decode(&result)
// 	if err != nil {
// 		if err == mongo.ErrNoDocuments {
// 			// This error means your query did not match any documents.
// 			return nil, nil
// 		}
// 		impl.Logger.Error("database get by email error", slog.Any("error", err))
// 		return nil, err
// 	}
// 	return &result, nil
// }
//
// func (impl comicSubmissionImplImpl) GetByVerificationCode(ctx context.Context, verificationCode string) (*domain.ComicSubmission, error) {
// 	filter := bson.M{"email_verification_code": verificationCode}
//
// 	var result domain.ComicSubmission
// 	err := impl.Collection.FindOne(ctx, filter).Decode(&result)
// 	if err != nil {
// 		if err == mongo.ErrNoDocuments {
// 			// This error means your query did not match any documents.
// 			return nil, nil
// 		}
// 		impl.Logger.Error("database get by verification code error", slog.Any("error", err))
// 		return nil, err
// 	}
// 	return &result, nil
// }
//
// func (impl comicSubmissionImplImpl) DeleteByID(ctx context.Context, id primitive.ObjectID) error {
// 	_, err := impl.Collection.DeleteOne(ctx, bson.M{"_id": id})
// 	if err != nil {
// 		impl.Logger.Error("database failed deletion error",
// 			slog.Any("error", err))
// 		return err
// 	}
// 	return nil
// }
//
// func (impl comicSubmissionImplImpl) CheckIfExistsByID(ctx context.Context, id primitive.ObjectID) (bool, error) {
// 	filter := bson.M{"_id": id}
// 	count, err := impl.Collection.CountDocuments(ctx, filter)
// 	if err != nil {
// 		impl.Logger.Error("database check if exists by ID error", slog.Any("error", err))
// 		return false, err
// 	}
// 	return count >= 1, nil
// }
//
// func (impl comicSubmissionImplImpl) CheckIfExistsByEmail(ctx context.Context, email string) (bool, error) {
// 	filter := bson.M{"email": email}
// 	count, err := impl.Collection.CountDocuments(ctx, filter)
// 	if err != nil {
// 		impl.Logger.Error("database check if exists by email error", slog.Any("error", err))
// 		return false, err
// 	}
// 	return count >= 1, nil
// }
//
// func (impl comicSubmissionImplImpl) UpdateByID(ctx context.Context, m *domain.ComicSubmission) error {
// 	filter := bson.M{"_id": m.ID}
//
// 	update := bson.M{ // DEVELOPERS NOTE: https://stackoverflow.com/a/60946010
// 		"$set": m,
// 	}
//
// 	// execute the UpdateOne() function to update the first matching document
// 	_, err := impl.Collection.UpdateOne(ctx, filter, update)
// 	if err != nil {
// 		impl.Logger.Error("database update by id error", slog.Any("error", err))
// 		return err
// 	}
//
// 	// // display the number of documents updated
// 	// impl.Logger.Debug("number of documents updated", slog.Int64("modified_count", result.ModifiedCount))
//
// 	return nil
// }
