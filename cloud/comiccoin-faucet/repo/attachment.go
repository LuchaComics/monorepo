package repo

import (
	"context"
	"log"
	"log/slog"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

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

	// Note:
	// * 1 for ascending
	// * -1 for descending
	// * "text" for text indexes

	// The following few lines of code will create the index for our app for this
	// colleciton.
	_, err := uc.Indexes().CreateMany(context.TODO(), []mongo.IndexModel{
		{Keys: bson.D{{Key: "object_key", Value: 1}}},
		{Keys: bson.D{{Key: "user_id", Value: 1}}},
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

//	func (impl attachmentImpl) GetByEmail(ctx context.Context, email string) (*domain.Attachment, error) {
//		filter := bson.M{"email": email}
//
//		var result domain.Attachment
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
//	func (impl attachmentImpl) GetByVerificationCode(ctx context.Context, verificationCode string) (*domain.Attachment, error) {
//		filter := bson.M{"email_verification_code": verificationCode}
//
//		var result domain.Attachment
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
//	func (impl attachmentImpl) DeleteByID(ctx context.Context, id primitive.ObjectID) error {
//		_, err := impl.Collection.DeleteOne(ctx, bson.M{"_id": id})
//		if err != nil {
//			impl.Logger.Error("database failed deletion error",
//				slog.Any("error", err))
//			return err
//		}
//		return nil
//	}
//
//	func (impl attachmentImpl) CheckIfExistsByID(ctx context.Context, id primitive.ObjectID) (bool, error) {
//		filter := bson.M{"_id": id}
//		count, err := impl.Collection.CountDocuments(ctx, filter)
//		if err != nil {
//			impl.Logger.Error("database check if exists by ID error", slog.Any("error", err))
//			return false, err
//		}
//		return count >= 1, nil
//	}
//
//	func (impl attachmentImpl) CheckIfExistsByEmail(ctx context.Context, email string) (bool, error) {
//		filter := bson.M{"email": email}
//		count, err := impl.Collection.CountDocuments(ctx, filter)
//		if err != nil {
//			impl.Logger.Error("database check if exists by email error", slog.Any("error", err))
//			return false, err
//		}
//		return count >= 1, nil
//	}
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
