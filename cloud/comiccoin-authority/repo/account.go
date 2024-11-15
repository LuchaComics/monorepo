package repo

import (
	"context"
	"log"
	"log/slog"
	"sort"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/common/blockchain/signature"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/config"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/domain"
)

type AccountRepo struct {
	config     *config.Configuration
	logger     *slog.Logger
	dbClient   *mongo.Client
	collection *mongo.Collection
}

func NewAccountRepo(cfg *config.Configuration, logger *slog.Logger, client *mongo.Client) *AccountRepo {
	// ctx := context.Background()
	uc := client.Database(cfg.DB.Name).Collection("accounts")

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

	return &AccountRepo{
		config:     cfg,
		logger:     logger,
		dbClient:   client,
		collection: uc,
	}
}

func (r *AccountRepo) Upsert(ctx context.Context, account *domain.Account) error {
	opts := options.Update().SetUpsert(true)
	_, err := r.collection.UpdateOne(ctx, bson.M{"address": account.Address}, bson.M{"$set": account}, opts)
	return err
}

func (r *AccountRepo) GetByAddress(ctx context.Context, address *common.Address) (*domain.Account, error) {
	var account domain.Account
	err := r.collection.FindOne(ctx, bson.M{"address": address}).Decode(&account)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &account, nil
}

func (r *AccountRepo) ListAll(ctx context.Context) ([]*domain.Account, error) {
	var accounts []*domain.Account
	cur, err := r.collection.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	for cur.Next(ctx) {
		var account domain.Account
		err := cur.Decode(&account)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, &account)
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}
	return accounts, nil
}

func (r *AccountRepo) ListWithFilterByAddresses(ctx context.Context, addrs []*common.Address) ([]*domain.Account, error) {
	var accounts []*domain.Account
	filter := bson.M{
		"address": bson.M{
			"$in": addrs,
		},
	}
	cur, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	for cur.Next(ctx) {
		var account domain.Account
		err := cur.Decode(&account)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, &account)
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}
	return accounts, nil
}

func (r *AccountRepo) DeleteByAddress(ctx context.Context, address *common.Address) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"address": address})
	return err
}

func (r *AccountRepo) HashState(ctx context.Context) (string, error) {
	accounts, err := r.ListAll(ctx)
	if err != nil {
		return "", err
	}

	// Variable used to only store the accounts which have a balance greater
	// then the value of zero.
	accountsWithBalance := make([]*domain.Account, 0)

	// Iterate through all the accounts and only save the accounts with balance.
	for _, account := range accounts {
		if account.Balance > 0 {
			accountsWithBalance = append(accountsWithBalance, account)
		}
	}

	// Sort the accounts by address
	sort.Sort(byAccount(accountsWithBalance))

	// Serialize the accounts to JSON
	accountsBytes := make([]byte, 0)
	for _, account := range accountsWithBalance {
		// DEVELOPERS NOTE:
		// In Go, the order of struct fields is determined by the order in which
		// they are defined in the struct. However, this order is not guaranteed
		// to be the same across different nodes or even different runs of the
		// same program.
		//
		// To fix this issue, you can use a deterministic serialization
		// algorithm, such as JSON or CBOR, to serialize the Account struct
		// before hashing it. This will ensure that the fields are always
		// serialized in the same order, regardless of the node or run.
		accountBytes, err := account.Serialize()
		if err != nil {
			return "", err
		}
		accountsBytes = append(accountsBytes, accountBytes...)
	}

	// Hash the deterministic serialized accounts
	res := signature.Hash(accountsBytes)
	return res, nil
}

// =============================================================================

// byAccount provides sorting support by the account id value.
type byAccount []*domain.Account

// Len returns the number of transactions in the list.
func (ba byAccount) Len() int {
	return len(ba)
}

// Less helps to sort the list by account id in ascending order to keep the
// accounts in the right order of processing.
func (ba byAccount) Less(i, j int) bool {
	return strings.ToLower(ba[i].Address.String()) < strings.ToLower(ba[j].Address.String())
}

// Swap moves accounts in the order of the account id value.
func (ba byAccount) Swap(i, j int) {
	ba[i], ba[j] = ba[j], ba[i]
}

func (r *AccountRepo) OpenTransaction() error {
	defer log.Fatal("Unsupported feature in the `comiccoin-authority` repository.")
	return nil
}

func (r *AccountRepo) CommitTransaction() error {
	defer log.Fatal("Unsupported feature in the `comiccoin-authority` repository.")
	return nil
}

func (r *AccountRepo) DiscardTransaction() {
	log.Fatal("Unsupported feature in the `comiccoin-authority` repository.")

}
