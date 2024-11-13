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

	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/config"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/domain"
)

type LibP2PNetworkPeerUniqueIdentifierRepo struct {
	config     *config.Configuration
	logger     *slog.Logger
	dbClient   *mongo.Client
	collection *mongo.Collection
}

func NewLibP2PNetworkPeerUniqueIdentifierRepo(cfg *config.Configuration, logger *slog.Logger, client *mongo.Client) domain.LibP2PNetworkPeerUniqueIdentifierRepository {
	// ctx := context.Background()
	uc := client.Database(cfg.DB.Name).Collection("blockchain_network_peer_unique_identifiers")

	// Note:
	// * 1 for ascending
	// * -1 for descending
	// * "text" for text indexes

	// The following few lines of code will create the index for our app for this
	// colleciton.
	_, err := uc.Indexes().CreateMany(context.TODO(), []mongo.IndexModel{
		{Keys: bson.D{{Key: "_id", Value: 1}}},
		{Keys: bson.D{{Key: "label", Value: 1}}},
		{Keys: bson.D{
			{Key: "label", Value: "text"},
		}},
	})
	if err != nil {
		// It is important that we crash the app on startup to meet the
		// requirements of `google/wire` framework.
		log.Fatal(err)
	}

	return &LibP2PNetworkPeerUniqueIdentifierRepo{
		config:     cfg,
		logger:     logger,
		dbClient:   client,
		collection: uc,
	}
}

func (r *LibP2PNetworkPeerUniqueIdentifierRepo) GetOrCreate(ctx context.Context, label string) (*domain.LibP2PNetworkPeerUniqueIdentifier, error) {
	//
	// STEP 1: Get (if it exists)
	//

	data, err := r.GetByLabel(ctx, label)
	if err != nil {
		r.logger.Error("Failed getting by label",
			slog.Any("error", err))
		return nil, err
	}
	if data != nil {
		return data, nil
	}

	//
	// STEP 2: Create new record
	//

	newData, err := domain.NewLibP2PNetworkPeerUniqueIdentifier(label)
	if err != nil {
		r.logger.Error("Failed creating new by label",
			slog.Any("error", err))
		return nil, err
	}
	if err := r.Upsert(ctx, newData); err != nil {
		r.logger.Error("Failed upserting",
			slog.Any("error", err))
		return nil, err
	}

	//
	// STEP 3: Get newly created record.
	//

	return newData, nil
}

func (r *LibP2PNetworkPeerUniqueIdentifierRepo) Upsert(ctx context.Context, data *domain.LibP2PNetworkPeerUniqueIdentifier) error {
	ctxWithTimeout, cancel := context.WithTimeout(ctx, 30*time.Second) // Use to prevent resource leaks.
	defer cancel()

	// Defensive Code: No empty ID values are allowed.
	if data.ID.IsZero() {
		data.ID = primitive.NewObjectID()
	}

	opts := options.Update().SetUpsert(true)
	_, err := r.collection.UpdateOne(ctxWithTimeout, bson.M{
		"_id": data.ID,
	}, bson.M{"$set": data}, opts)
	return err
}

func (r *LibP2PNetworkPeerUniqueIdentifierRepo) GetByLabel(ctx context.Context, label string) (*domain.LibP2PNetworkPeerUniqueIdentifier, error) {
	var data domain.LibP2PNetworkPeerUniqueIdentifier
	err := r.collection.FindOne(ctx, bson.M{"label": label}).Decode(&data)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &data, nil
}
