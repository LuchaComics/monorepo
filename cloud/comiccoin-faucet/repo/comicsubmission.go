package repo

import (
	"context"
	"log"
	"log/slog"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

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

	// Note:
	// * 1 for ascending
	// * -1 for descending
	// * "text" for text indexes

	// The following few lines of code will create the index for our app for this
	// colleciton.
	_, err := uc.Indexes().CreateMany(context.TODO(), []mongo.IndexModel{
		{Keys: bson.D{{Key: "user_id", Value: 1}}},
		{Keys: bson.D{
			{Key: "name", Value: "text"},
		}},
	})
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
