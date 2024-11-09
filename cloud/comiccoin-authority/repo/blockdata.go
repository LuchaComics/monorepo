package repo

import (
	"context"
	"log"
	"log/slog"

	"github.com/ethereum/go-ethereum/common"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/config"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/domain"
)

type BlockDataRepo struct {
	config     *config.Configuration
	logger     *slog.Logger
	dbClient   *mongo.Client
	collection *mongo.Collection
}

func NewBlockDataRepo(cfg *config.Configuration, logger *slog.Logger, client *mongo.Client) *BlockDataRepo {
	// ctx := context.Background()
	uc := client.Database(cfg.DB.Name).Collection("blockdata")

	// Note:
	// * 1 for ascending
	// * -1 for descending
	// * "text" for text indexes

	// The following few lines of code will create the index for our app for this
	// colleciton.
	_, err := uc.Indexes().CreateMany(context.TODO(), []mongo.IndexModel{
		{Keys: bson.D{{Key: "chainID", Value: 1}}},
		{Keys: bson.D{
			{Key: "chainID", Value: "text"},
		}},
	})
	if err != nil {
		// It is important that we crash the app on startup to meet the
		// requirements of `google/wire` framework.
		log.Fatal(err)
	}

	return &BlockDataRepo{
		config:     cfg,
		logger:     logger,
		dbClient:   client,
		collection: uc,
	}
}

func (r *BlockDataRepo) Upsert(ctx context.Context, blockdata *domain.BlockData) error {
	opts := options.Update().SetUpsert(true)
	_, err := r.collection.UpdateOne(ctx, bson.M{"hash": blockdata.Hash}, bson.M{"$set": blockdata}, opts)
	return err
}

func (r *BlockDataRepo) GetByHash(ctx context.Context, hash string) (*domain.BlockData, error) {
	var blockData domain.BlockData
	err := r.collection.FindOne(ctx, bson.M{"hash": hash}).Decode(&blockData)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &blockData, nil
}

func (r *BlockDataRepo) ListAll(ctx context.Context) ([]*domain.BlockData, error) {
	var blockDatas []*domain.BlockData
	cur, err := r.collection.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	for cur.Next(ctx) {
		var blockData domain.BlockData
		err := cur.Decode(&blockData)
		if err != nil {
			return nil, err
		}
		blockDatas = append(blockDatas, &blockData)
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}
	return blockDatas, nil
}

func (r *BlockDataRepo) DeleteByHash(ctx context.Context, hash string) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"hash": hash})
	return err
}

func (r *BlockDataRepo) ListAllBlockTransactionsByAddress(ctx context.Context, address *common.Address) ([]*domain.BlockTransaction, error) {
	var blockTransactions []*domain.BlockTransaction
	cur, err := r.collection.Find(ctx, bson.M{
		"$or": []bson.M{
			{"trans.signedtransaction.from": address},
			{"trans.signedtransaction.to": address},
		},
	})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	for cur.Next(ctx) {
		var blockData domain.BlockData
		err := cur.Decode(&blockData)
		if err != nil {
			return nil, err
		}
		for _, trans := range blockData.Trans {
			if trans.SignedTransaction.From == address || trans.SignedTransaction.To == address {
				blockTransactions = append(blockTransactions, &trans)
			}
		}
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}
	return blockTransactions, nil
}

func (r *BlockDataRepo) GetByBlockTransactionTimestamp(ctx context.Context, timestamp uint64) (*domain.BlockData, error) {
	var blockData domain.BlockData
	err := r.collection.FindOne(ctx, bson.M{"trans.timestamp": timestamp}).Decode(&blockData)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &blockData, nil
}
