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

type WalletRepo struct {
	config     *config.Configuration
	logger     *slog.Logger
	dbClient   *mongo.Client
	collection *mongo.Collection
}

func NewWalletRepo(cfg *config.Configuration, logger *slog.Logger, client *mongo.Client) *WalletRepo {
	// ctx := context.Background()
	uc := client.Database(cfg.DB.Name).Collection("wallets")

	// Note:
	// * 1 for ascending
	// * -1 for descending
	// * "text" for text indexes

	// The following few lines of code will create the index for our app for this
	// colleciton.
	_, err := uc.Indexes().CreateMany(context.TODO(), []mongo.IndexModel{
		{Keys: bson.D{{Key: "address", Value: 1}}},
		{Keys: bson.D{
			{Key: "address", Value: "text"},
		}},
	})
	if err != nil {
		// It is important that we crash the app on startup to meet the
		// requirements of `google/wire` framework.
		log.Fatal(err)
	}

	return &WalletRepo{
		config:     cfg,
		logger:     logger,
		dbClient:   client,
		collection: uc,
	}
}

func (r *WalletRepo) Upsert(ctx context.Context, wallet *domain.Wallet) error {
	opts := options.Update().SetUpsert(true)
	_, err := r.collection.UpdateOne(ctx, bson.M{"address": wallet.Address}, bson.M{"$set": wallet}, opts)
	return err
}

func (r *WalletRepo) GetByAddress(ctx context.Context, address *common.Address) (*domain.Wallet, error) {
	var wallet domain.Wallet
	err := r.collection.FindOne(ctx, bson.M{"address": address}).Decode(&wallet)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &wallet, nil
}

func (r *WalletRepo) ListAll(ctx context.Context) ([]*domain.Wallet, error) {
	var wallets []*domain.Wallet
	cur, err := r.collection.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	for cur.Next(ctx) {
		var wallet domain.Wallet
		err := cur.Decode(&wallet)
		if err != nil {
			return nil, err
		}
		wallets = append(wallets, &wallet)
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}
	return wallets, nil
}

func (r *WalletRepo) DeleteByAddress(ctx context.Context, address *common.Address) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"address": address})
	return err
}

func (r *WalletRepo) OpenTransaction() error {
	defer log.Fatal("Unsupported feature in the `comiccoin-authority` repository.")
	return nil
}

func (r *WalletRepo) CommitTransaction() error {
	defer log.Fatal("Unsupported feature in the `comiccoin-authority` repository.")
	return nil
}

func (r *WalletRepo) DiscardTransaction() {
	log.Fatal("Unsupported feature in the `comiccoin-authority` repository.")

}
