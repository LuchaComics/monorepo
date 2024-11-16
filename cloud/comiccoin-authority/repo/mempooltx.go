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

type MempoolTransactionRepo struct {
	config     *config.Configuration
	logger     *slog.Logger
	dbClient   *mongo.Client
	collection *mongo.Collection
}

func NewMempoolTransactionRepo(cfg *config.Configuration, logger *slog.Logger, client *mongo.Client) *MempoolTransactionRepo {
	// ctx := context.Background()
	uc := client.Database(cfg.DB.Name).Collection("mempool_transactions")

	// Note:
	// * 1 for ascending
	// * -1 for descending
	// * "text" for text indexes

	// The following few lines of code will create the index for our app for this
	// colleciton.
	_, err := uc.Indexes().CreateMany(context.TODO(), []mongo.IndexModel{
		{Keys: bson.D{{Key: "transaction.chain_id", Value: 1}}},
		{Keys: bson.D{{Key: "transaction.nonce", Value: 1}}},
		{Keys: bson.D{
			{Key: "transaction.data", Value: "text"},
		}},
	})
	if err != nil {
		// It is important that we crash the app on startup to meet the
		// requirements of `google/wire` framework.
		log.Fatal(err)
	}

	return &MempoolTransactionRepo{
		config:     cfg,
		logger:     logger,
		dbClient:   client,
		collection: uc,
	}
}

func (r *MempoolTransactionRepo) GetByID(ctx context.Context, id primitive.ObjectID) (*domain.MempoolTransaction, error) {
	var account domain.MempoolTransaction
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&account)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &account, nil
}

func (r *MempoolTransactionRepo) Upsert(ctx context.Context, mempoolTx *domain.MempoolTransaction) error {
	ctxWithTimeout, cancel := context.WithTimeout(ctx, 30*time.Second) // Use to prevent resource leaks.
	defer cancel()

	// Defensive Code: No empty ID values are allowed.
	if mempoolTx.ID.IsZero() {
		mempoolTx.ID = primitive.NewObjectID()
	}

	opts := options.Update().SetUpsert(true)
	_, err := r.collection.UpdateOne(ctxWithTimeout, bson.M{
		"_id": mempoolTx.ID,
	}, bson.M{"$set": mempoolTx}, opts)
	return err
}

func (r *MempoolTransactionRepo) ListByChainID(ctx context.Context, chainID uint16) ([]*domain.MempoolTransaction, error) {
	ctxWithTimeout, cancel := context.WithTimeout(ctx, 30*time.Second) // Use to prevent resource leaks.
	defer cancel()

	mempoolTxs := make([]*domain.MempoolTransaction, 0)
	cur, err := r.collection.Find(ctxWithTimeout, bson.M{"chain_id": chainID})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctxWithTimeout)
	for cur.Next(ctxWithTimeout) {
		var mempoolTx domain.MempoolTransaction
		err := cur.Decode(&mempoolTx)
		if err != nil {
			return nil, err
		}
		mempoolTxs = append(mempoolTxs, &mempoolTx)
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}
	return mempoolTxs, nil
}

func (r *MempoolTransactionRepo) DeleteByChainID(ctx context.Context, chainID uint16) error {
	ctxWithTimeout, cancel := context.WithTimeout(ctx, 30*time.Second) // Use to prevent resource leaks.
	defer cancel()
	_, err := r.collection.DeleteMany(ctxWithTimeout, bson.M{"chain_id": chainID})
	return err
}

func (r *MempoolTransactionRepo) DeleteByID(ctx context.Context, id primitive.ObjectID) error {
	ctxWithTimeout, cancel := context.WithTimeout(ctx, 30*time.Second) // Use to prevent resource leaks.
	defer cancel()
	_, err := r.collection.DeleteOne(ctxWithTimeout, bson.M{"_id": id})
	return err
}

// -----------------------------------------------------------------------------

func (r *MempoolTransactionRepo) getInsertionChangeStream(ctx context.Context) (*mongo.ChangeStream, error) {
	pipeline := mongo.Pipeline{bson.D{{"$match", bson.D{{"$or",
		bson.A{
			bson.D{{"operationType", "insert"}}}}},
	}}}

	changeStream, err := r.collection.Watch(ctx, pipeline, options.ChangeStream().SetFullDocument(options.UpdateLookup))
	if err != nil {
		return nil, err
	}

	return changeStream, nil
}

func (r *MempoolTransactionRepo) GetInsertionChangeStreamChannel(ctx context.Context) (<-chan domain.MempoolTransaction, chan struct{}, error) {
	changeStream, err := r.getInsertionChangeStream(ctx)
	if err != nil {
		return nil, nil, err
	}

	dataChan := make(chan domain.MempoolTransaction)
	quitChan := make(chan struct{})
	go func() {
		defer close(dataChan)
		for changeStream.Next(ctx) {
			// r.logger.Debug("Running next...")
			select {
			case <-quitChan:
				// r.logger.Debug("Quit chan!")
				changeStream.Close(ctx)
				return
			default:
				// r.logger.Debug("OK...")
			}

			// SPECIAL THANKS: https://stackoverflow.com/a/74519310
			var event struct {
				Doc domain.MempoolTransaction `bson:"fullDocument"`
			}

			if err := changeStream.Decode(&event); err != nil {
				r.logger.Error("Failed to decode event",
					slog.Any("error", err))
				continue
			}

			// r.logger.Debug("Ready to send...",
			// 	slog.Any("dataChan", event.Doc))

			dataChan <- event.Doc
		}
		changeStream.Close(ctx)
	}()
	return dataChan, quitChan, nil
}
