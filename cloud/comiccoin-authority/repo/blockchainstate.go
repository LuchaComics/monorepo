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

type BlockchainStateRepo struct {
	config     *config.Configuration
	logger     *slog.Logger
	dbClient   *mongo.Client
	collection *mongo.Collection
}

func NewBlockchainStateRepo(cfg *config.Configuration, logger *slog.Logger, client *mongo.Client) *BlockchainStateRepo {
	// ctx := context.Background()
	uc := client.Database(cfg.DB.Name).Collection("blockchain_states")

	// Note:
	// * 1 for ascending
	// * -1 for descending
	// * "text" for text indexes

	// The following few lines of code will create the index for our app for this
	// colleciton.
	_, err := uc.Indexes().CreateMany(context.TODO(), []mongo.IndexModel{
		{Keys: bson.D{{Key: "chain_id", Value: 1}}},
		{Keys: bson.D{
			{Key: "chain_id", Value: "text"},
		}},
	})
	if err != nil {
		// It is important that we crash the app on startup to meet the
		// requirements of `google/wire` framework.
		log.Fatal(err)
	}

	return &BlockchainStateRepo{
		config:     cfg,
		logger:     logger,
		dbClient:   client,
		collection: uc,
	}
}

func (r *BlockchainStateRepo) UpsertByChainID(ctx context.Context, blockchainState *domain.BlockchainState) error {
	opts := options.Update().SetUpsert(true)
	_, err := r.collection.UpdateOne(ctx, bson.M{"chain_id": blockchainState.ChainID}, bson.M{"$set": blockchainState}, opts)
	return err
}

func (r *BlockchainStateRepo) GetByChainID(ctx context.Context, chainID uint16) (*domain.BlockchainState, error) {
	var blockchainState domain.BlockchainState
	err := r.collection.FindOne(ctx, bson.M{"chain_id": chainID}).Decode(&blockchainState)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &blockchainState, nil
}

func (r *BlockchainStateRepo) GetForMainNet(ctx context.Context) (*domain.BlockchainState, error) {
	return r.GetByChainID(ctx, 1)
}

func (r *BlockchainStateRepo) GetForTestNet(ctx context.Context) (*domain.BlockchainState, error) {
	return r.GetByChainID(ctx, 2)
}

func (r *BlockchainStateRepo) ListAll(ctx context.Context) ([]*domain.BlockchainState, error) {
	var blockchainStates []*domain.BlockchainState
	cur, err := r.collection.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	for cur.Next(ctx) {
		var blockchainState domain.BlockchainState
		err := cur.Decode(&blockchainState)
		if err != nil {
			return nil, err
		}
		blockchainStates = append(blockchainStates, &blockchainState)
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}
	return blockchainStates, nil
}

func (r *BlockchainStateRepo) DeleteByChainID(ctx context.Context, chainID uint16) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"chain_id": chainID})
	return err
}

func (r *BlockchainStateRepo) OpenTransaction() error {
	defer log.Fatal("Unsupported feature in the `comiccoin-authority` repository.")
	return nil
}

func (r *BlockchainStateRepo) CommitTransaction() error {
	defer log.Fatal("Unsupported feature in the `comiccoin-authority` repository.")
	return nil
}

func (r *BlockchainStateRepo) DiscardTransaction() {
	log.Fatal("Unsupported feature in the `comiccoin-authority` repository.")
}
