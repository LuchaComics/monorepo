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

type GenesisBlockDataRepo struct {
	config     *config.Configuration
	logger     *slog.Logger
	dbClient   *mongo.Client
	collection *mongo.Collection
}

func NewGenesisBlockDataRepo(cfg *config.Configuration, logger *slog.Logger, client *mongo.Client) *GenesisBlockDataRepo {
	// ctx := context.Background()
	uc := client.Database(cfg.DB.Name).Collection("genesis_blockdata")

	// Note:
	// * 1 for ascending
	// * -1 for descending
	// * "text" for text indexes

	// The following few lines of code will create the index for our app for this
	// colleciton.
	_, err := uc.Indexes().CreateMany(context.TODO(), []mongo.IndexModel{
		{Keys: bson.D{{Key: "header.chain_id", Value: 1}}},
		{Keys: bson.D{
			{Key: "hash", Value: "text"},
		}},
	})
	if err != nil {
		// It is important that we crash the app on startup to meet the
		// requirements of `google/wire` framework.
		log.Fatal(err)
	}

	return &GenesisBlockDataRepo{
		config:     cfg,
		logger:     logger,
		dbClient:   client,
		collection: uc,
	}
}

func (r *GenesisBlockDataRepo) UpsertByChainID(ctx context.Context, blockdata *domain.GenesisBlockData) error {
	opts := options.Update().SetUpsert(true)
	_, err := r.collection.UpdateOne(ctx, bson.M{"header.chain_id": blockdata.Header.ChainID}, bson.M{"$set": blockdata}, opts)
	return err
}

func (r *GenesisBlockDataRepo) GetByChainID(ctx context.Context, chainID uint16) (*domain.GenesisBlockData, error) {
	var blockData domain.GenesisBlockData
	err := r.collection.FindOne(ctx, bson.M{"header.chain_id": chainID}).Decode(&blockData)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &blockData, nil
}

func (r *GenesisBlockDataRepo) OpenTransaction() error {
	defer log.Fatal("Unsupported feature in the `comiccoin-authority` repository.")
	return nil
}

func (r *GenesisBlockDataRepo) CommitTransaction() error {
	defer log.Fatal("Unsupported feature in the `comiccoin-authority` repository.")
	return nil
}

func (r *GenesisBlockDataRepo) DiscardTransaction() {
	log.Fatal("Unsupported feature in the `comiccoin-authority` repository.")
}
