package repo

import (
	"context"
	"log"
	"log/slog"

	"go.mongodb.org/mongo-driver/bson"
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

func (r *MempoolTransactionRepo) Upsert(ctx context.Context, mempoolTx *domain.MempoolTransaction) error {
	opts := options.Update().SetUpsert(true)
	_, err := r.collection.UpdateOne(ctx, bson.M{
		"chain_id":          mempoolTx.ChainID,
		"transaction.nonce": mempoolTx.Transaction.Nonce,
	}, bson.M{"$set": mempoolTx}, opts)
	return err
}

func (r *MempoolTransactionRepo) ListByChainID(ctx context.Context, chainID uint16) ([]*domain.MempoolTransaction, error) {
	mempoolTxs := make([]*domain.MempoolTransaction, 0)
	cur, err := r.collection.Find(ctx, bson.M{"chain_id": chainID})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	for cur.Next(ctx) {
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
	//TODO: Impl.
	return nil
}
