package repo

import (
	"bytes"
	"context"
	"encoding/hex"
	"fmt"
	"log"
	"log/slog"

	"github.com/ethereum/go-ethereum/common"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/config"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/domain"
)

type BlockDataRepo struct {
	config     *config.Configuration
	logger     *slog.Logger
	dbClient   *mongo.Client
	collection *mongo.Collection
}

func NewBlockDataRepository(cfg *config.Configuration, logger *slog.Logger, client *mongo.Client) domain.BlockDataRepository {
	// ctx := context.Background()
	uc := client.Database(cfg.DB.Name).Collection("blockdata")

	// Note:
	// * 1 for ascending
	// * -1 for descending
	// * "text" for text indexes

	// The following few lines of code will create the index for our app for this
	// colleciton.
	_, err := uc.Indexes().CreateMany(context.TODO(), []mongo.IndexModel{
		{Keys: bson.D{{Key: "hash", Value: 1}}},
		{Keys: bson.D{{Key: "header.chain_id", Value: 1}}},
		{Keys: bson.D{{Key: "header.number", Value: 1}}},
		{Keys: bson.D{{Key: "header.timestamp", Value: 1}}},
		{Keys: bson.D{
			{Key: "hash", Value: "text"},
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

func (r *BlockDataRepo) ListByChainID(ctx context.Context, chainID uint16) ([]*domain.BlockData, error) {
	blockDatas := make([]*domain.BlockData, 0)
	filter := bson.M{"header.chain_id": chainID}
	cur, err := r.collection.Find(ctx, filter)
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

// ListInHashes method is deprecated.
func (r *BlockDataRepo) ListInHashes(ctx context.Context, hashes []string) ([]*domain.BlockData, error) {
	blockDatas := make([]*domain.BlockData, 0)
	filter := bson.M{
		"hash": bson.M{
			"$in": hashes,
		},
	}
	cur, err := r.collection.Find(ctx, filter)
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

// ListInBetweenBlockNumbersForChainID method is deprecated.
func (r *BlockDataRepo) ListInBetweenBlockNumbersForChainID(ctx context.Context, startBlockNumber, finishBlockNumber uint64, chainID uint16) ([]*domain.BlockData, error) {
	blockDatas := make([]*domain.BlockData, 0)
	filter := bson.M{
		"header.chain_id": chainID,
		"header.number": bson.M{
			"$gte": startBlockNumber,
			"$lte": finishBlockNumber,
		},
	}
	cur, err := r.collection.Find(ctx, filter)
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

// ListBlockNumberByHashArrayForChainID method is deprecated.
func (r *BlockDataRepo) ListBlockNumberByHashArrayForChainID(ctx context.Context, chainID uint16) ([]domain.BlockNumberByHash, error) {
	blockNumberByHashArray := make([]domain.BlockNumberByHash, 0)
	projection := bson.M{
		"_id":    0,
		"number": 1,
		"hash":   1,
	}
	filter := bson.M{
		"header.chain_id": chainID,
	}
	opts := options.Find().SetProjection(projection)
	cur, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	var blockData domain.BlockNumberByHash
	for cur.Next(ctx) {
		err := cur.Decode(&blockData)
		if err != nil {
			return nil, err
		}
		blockNumberByHashArray = append(blockNumberByHashArray, blockData)
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}
	return blockNumberByHashArray, nil
}

// ListUnorderedHashArrayForChainID method is deprecated.
func (r *BlockDataRepo) ListUnorderedHashArrayForChainID(ctx context.Context, chainID uint16) ([]string, error) {
	hashArray := make([]string, 0)
	projection := bson.M{
		"_id":  0,
		"hash": 1,
	}
	filter := bson.M{
		"header.chain_id": chainID,
	}
	opts := options.Find().SetProjection(projection)
	cur, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	var blockData struct {
		Hash string `bson:"hash"`
	}
	for cur.Next(ctx) {
		err := cur.Decode(&blockData)
		if err != nil {
			return nil, err
		}
		hashArray = append(hashArray, blockData.Hash)
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}
	return hashArray, nil
}

func (r *BlockDataRepo) ListBlockTransactionsByAddress(ctx context.Context, address *common.Address) ([]*domain.BlockTransaction, error) {
	// Result slice
	var blockTransactions []*domain.BlockTransaction

	// Dereference the address pointer for the query
	if address == nil {
		return nil, fmt.Errorf("address cannot be nil")
	}
	// Decode the hex string (strip "0x" prefix)
	addressBytes, err := hex.DecodeString(address.Hex()[2:])
	if err != nil {
		return nil, fmt.Errorf("failed to decode address: %v", err)
	}

	// MongoDB filter
	filter := bson.M{
		"trans": bson.M{
			"$elemMatch": bson.M{
				"$or": []bson.M{
					{"signedtransaction.transaction.from": addressBytes},
					{"signedtransaction.transaction.to": addressBytes},
				},
			},
		},
	}

	// Execute the query
	cur, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	// Iterate over the cursor
	for cur.Next(ctx) {
		var blockData domain.BlockData
		if err := cur.Decode(&blockData); err != nil {
			return nil, err
		}

		// Filter transactions by `from` or `to` address
		for _, trans := range blockData.Trans {
			if (trans.SignedTransaction.Transaction.From != nil && bytes.Equal(trans.SignedTransaction.Transaction.From.Bytes(), addressBytes)) ||
				(trans.SignedTransaction.Transaction.To != nil && bytes.Equal(trans.SignedTransaction.Transaction.To.Bytes(), addressBytes)) {
				blockTransactions = append(blockTransactions, &trans)
			}
		}
	}

	// Check for cursor errors
	if err := cur.Err(); err != nil {
		return nil, err
	}

	return blockTransactions, nil
}

func (r *BlockDataRepo) ListWithLimitForBlockTransactionsByAddress(ctx context.Context, address *common.Address, limit int64) ([]*domain.BlockTransaction, error) {
	// Result slice
	// Result slice
	var blockTransactions []*domain.BlockTransaction

	// Dereference the address pointer for the query
	if address == nil {
		return nil, fmt.Errorf("address cannot be nil")
	}
	// Decode the hex string (strip "0x" prefix)
	addressBytes, err := hex.DecodeString(address.Hex()[2:])
	if err != nil {
		return nil, fmt.Errorf("failed to decode address: %v", err)
	}

	// MongoDB filter
	filter := bson.M{
		"trans": bson.M{
			"$elemMatch": bson.M{
				"$or": []bson.M{
					{"signedtransaction.transaction.from": addressBytes},
					{"signedtransaction.transaction.to": addressBytes},
				},
			},
		},
	}
	options := options.Find().SetLimit(limit)

	// Execute the query
	cur, err := r.collection.Find(ctx, filter, options)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	// Iterate over the cursor
	for cur.Next(ctx) {
		var blockData domain.BlockData
		if err := cur.Decode(&blockData); err != nil {
			return nil, err
		}

		// Filter transactions by `from` or `to` address
		for _, trans := range blockData.Trans {
			if (trans.SignedTransaction.Transaction.From != nil && bytes.Equal(trans.SignedTransaction.Transaction.From.Bytes(), addressBytes)) ||
				(trans.SignedTransaction.Transaction.To != nil && bytes.Equal(trans.SignedTransaction.Transaction.To.Bytes(), addressBytes)) {
				blockTransactions = append(blockTransactions, &trans)
			}
		}
	}

	// Check for cursor errors
	if err := cur.Err(); err != nil {
		return nil, err
	}

	return blockTransactions, nil
}

func (r *BlockDataRepo) GetByBlockTransactionTimestamp(ctx context.Context, timestamp uint64) (*domain.BlockData, error) {
	var blockData domain.BlockData
	err := r.collection.FindOne(ctx, bson.M{"trans": bson.M{"$elemMatch": bson.M{"signed_transaction.timestamp": timestamp}}}).
		Decode(&blockData)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &blockData, nil
}

func (r *BlockDataRepo) OpenTransaction() error {
	defer log.Fatal("Unsupported feature in the `comiccoin-faucet` repository.")
	return nil
}

func (r *BlockDataRepo) CommitTransaction() error {
	defer log.Fatal("Unsupported feature in the `comiccoin-faucet` repository.")
	return nil
}

func (r *BlockDataRepo) DiscardTransaction() {
	log.Fatal("Unsupported feature in the `comiccoin-faucet` repository.")
}
