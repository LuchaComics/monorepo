package repo

import (
	"context"
	"log"
	"log/slog"
	"math/big"
	"sort"

	"github.com/ethereum/go-ethereum/common"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/common/blockchain/signature"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/config"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/domain"
)

type TokenRepo struct {
	config     *config.Configuration
	logger     *slog.Logger
	dbClient   *mongo.Client
	collection *mongo.Collection
}

func NewTokenRepo(cfg *config.Configuration, logger *slog.Logger, client *mongo.Client) *TokenRepo {
	// ctx := context.Background()
	uc := client.Database(cfg.DB.Name).Collection("tokens")

	// // For debugging purposes only or if you are going to recreate new indexes.
	// if _, err := uc.Indexes().DropAll(context.TODO()); err != nil {
	// 	logger.Warn("failed deleting all indexes",
	// 		slog.Any("err", err))
	//
	// 	// Do not crash app, just continue.
	// }

	// Note:
	// * 1 for ascending
	// * -1 for descending
	// * "text" for text indexes

	// The following few lines of code will create the index for our app for this
	// colleciton.
	_, err := uc.Indexes().CreateMany(context.TODO(), []mongo.IndexModel{
		{Keys: bson.D{{Key: "id_bytes", Value: 1}}},
		{Keys: bson.D{
			{Key: "id_bytes", Value: "text"},
		}},
	})
	if err != nil {
		// It is important that we crash the app on startup to meet the
		// requirements of `google/wire` framework.
		log.Fatal(err)
	}

	return &TokenRepo{
		config:     cfg,
		logger:     logger,
		dbClient:   client,
		collection: uc,
	}
}

func (r *TokenRepo) Upsert(ctx context.Context, token *domain.Token) error {
	opts := options.Update().SetUpsert(true)
	_, err := r.collection.UpdateOne(ctx, bson.M{"id_bytes": token.IDBytes}, bson.M{"$set": token}, opts)
	return err
}

func (r *TokenRepo) GetByID(ctx context.Context, id *big.Int) (*domain.Token, error) {
	var token domain.Token
	idBytes := id.Bytes()
	err := r.collection.FindOne(ctx, bson.M{"id_bytes": idBytes}).Decode(&token)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &token, nil
}

func (r *TokenRepo) ListAll(ctx context.Context) ([]*domain.Token, error) {
	var tokens []*domain.Token
	cur, err := r.collection.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	for cur.Next(ctx) {
		var token domain.Token
		err := cur.Decode(&token)
		if err != nil {
			return nil, err
		}
		tokens = append(tokens, &token)
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}
	return tokens, nil
}

func (r *TokenRepo) ListByOwner(ctx context.Context, owner *common.Address) ([]*domain.Token, error) {
	var tokens []*domain.Token
	cur, err := r.collection.Find(ctx, bson.M{"owner": owner})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	for cur.Next(ctx) {
		var token domain.Token
		err := cur.Decode(&token)
		if err != nil {
			return nil, err
		}
		tokens = append(tokens, &token)
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}
	return tokens, nil
}

func (r *TokenRepo) DeleteByID(ctx context.Context, id *big.Int) error {
	idBytes := id.Bytes()
	_, err := r.collection.DeleteOne(ctx, bson.M{"id_bytes": idBytes})
	return err
}

func (r *TokenRepo) HashState(ctx context.Context) (string, error) {
	tokens, err := r.ListAll(ctx)
	if err != nil {
		return "", err
	}

	// Sort and hash our tokens.
	sort.Sort(byToken(tokens))

	// Serialize the tokens to JSON
	tokensBytes := make([]byte, 0)
	for _, tok := range tokens {
		// DEVELOPERS NOTE:
		// In Go, the order of struct fields is determined by the order in which
		// they are defined in the struct. However, this order is not guaranteed
		// to be the same across different nodes or even different runs of the
		// same program.
		//
		// To fix this issue, you can use a deterministic serialization
		// algorithm, such as JSON or CBOR, to serialize the Token struct
		// before hashing it. This will ensure that the fields are always
		// serialized in the same order, regardless of the node or run.
		tokBytes, err := tok.Serialize()
		if err != nil {
			return "", err
		}
		tokensBytes = append(tokensBytes, tokBytes...)
	}

	// Hash the deterministic serialized tokens.
	return signature.Hash(tokensBytes), nil
}

// =============================================================================

// byToken provides sorting support by the token id value.
type byToken []*domain.Token

// Len returns the number of transactions in the list.
func (ba byToken) Len() int {
	return len(ba)
}

// Less helps to sort the list by token id in ascending order to keep the
// tokens in the right order of processing.
func (ba byToken) Less(i, j int) bool {
	return ba[i].GetID().Cmp(ba[j].GetID()) < 0
}

// Swap moves tokens in the order of the token id value.
func (ba byToken) Swap(i, j int) {
	ba[i], ba[j] = ba[j], ba[i]
}
