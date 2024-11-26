package repo

import (
	"context"
	"log"
	"log/slog"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/config"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/domain"
)

type BlockchainStateRepo struct {
	config     *config.Configuration
	logger     *slog.Logger
	dbClient   *mongo.Client
	collection *mongo.Collection
}

func NewBlockchainStateRepository(cfg *config.Configuration, logger *slog.Logger, client *mongo.Client) domain.BlockchainStateRepository {
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
	defer log.Fatal("Unsupported feature in the `comiccoin-faucet` repository.")
	return nil
}

func (r *BlockchainStateRepo) CommitTransaction() error {
	defer log.Fatal("Unsupported feature in the `comiccoin-faucet` repository.")
	return nil
}

func (r *BlockchainStateRepo) DiscardTransaction() {
	log.Fatal("Unsupported feature in the `comiccoin-faucet` repository.")
}

func (r *BlockchainStateRepo) getUpdateChangeStream(ctx context.Context) (*mongo.ChangeStream, error) {
	pipeline := mongo.Pipeline{bson.D{{"$match", bson.D{{"$or",
		bson.A{
			bson.D{{"operationType", "update"}}}}},
	}}}

	changeStream, err := r.collection.Watch(ctx, pipeline, options.ChangeStream().SetFullDocument(options.UpdateLookup))
	if err != nil {
		return nil, err
	}

	return changeStream, nil
}

func (r *BlockchainStateRepo) GetUpdateChangeStreamChannel(ctx context.Context) (<-chan domain.BlockchainState, chan struct{}, error) {
	changeStream, err := r.getUpdateChangeStream(ctx)
	if err != nil {
		return nil, nil, err
	}

	dataChan := make(chan domain.BlockchainState)
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
				Doc domain.BlockchainState `bson:"fullDocument"`
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
