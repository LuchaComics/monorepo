package domain

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/fxamacker/cbor/v2"
)

// BlockData represents the data that can be serialized to disk and over the network.
// It contains the hash of the block, the block header, and the list of transactions in the block.
type BlockData struct {
	// Hash is the unique hash of the block.
	Hash string `bson:"hash" json:"hash"`

	// Header is the block header, which contains metadata about the block.
	Header *BlockHeader `bson:"header" json:"header"`

	// The signature of this block's "Header" field which was applied by the
	// proof-of-authority validator.
	HeaderSignatureBytes []byte `bson:"header_signature_bytes" json:"header_signature_bytes"`

	// Trans is the list of (coin) transactions in the block.
	Trans []BlockTransaction `bson:"trans" json:"trans"`

	// The proof-of-authority validator whom executed the validation of
	// this block data in our blockchain.
	Validator *Validator `bson:"validator" json:"validator"`
}

type BlockNumberByHash struct {
	Number uint64 `bson:"number"`
	Hash   string `bson:"hash"`
}

// BlockDataRepository is an interface that defines the methods for interacting with block data.
// It provides methods for upserting, getting, listing, and deleting block data.
type BlockDataRepository interface {
	// Upsert upserts a block data into the repository.
	// It takes a block data and returns an error if one occurs.
	Upsert(ctx context.Context, bd *BlockData) error

	// GetByHash gets a block data by its hash.
	// It takes a hash and returns the block data and an error if one occurs.
	GetByHash(ctx context.Context, hash string) (*BlockData, error)

	// ListByChainID lists all block data in the repository for the particular chain.
	ListByChainID(ctx context.Context, chainID uint16) ([]*BlockData, error)

	// DeleteByHash deletes a block data by its hash.
	// It takes a hash and returns an error if one occurs.
	DeleteByHash(ctx context.Context, hash string) error

	// ListBlockTransactionsByAddress lists all the transactions for a particular address.
	ListBlockTransactionsByAddress(ctx context.Context, address *common.Address) ([]*BlockTransaction, error)

	// ListBlockTransactionsByAddress lists all the transactions for a particular address.
	ListWithLimitForBlockTransactionsByAddress(ctx context.Context, address *common.Address, limit int64) ([]*BlockTransaction, error)

	GetByBlockTransactionTimestamp(ctx context.Context, timestamp uint64) (*BlockData, error)

	OpenTransaction() error
	CommitTransaction() error
	DiscardTransaction()
}

// Serialize serializes a block data into a byte array.
// It returns the serialized byte array and an error if one occurs.
func (b *BlockData) Serialize() ([]byte, error) {
	// Marshal the block data into a byte array using CBOR.
	dataBytes, err := cbor.Marshal(b)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize block data: %v", err)
	}
	return dataBytes, nil
}

// NewBlockDataFromDeserialize deserializes a block data from a byte array.
// It returns the deserialized block data and an error if one occurs.
func NewBlockDataFromDeserialize(data []byte) (*BlockData, error) {
	// Variable we will use to return.
	blockData := &BlockData{}

	// Defensive code: If programmer entered empty bytes then we will
	// return nil deserialization result.
	if data == nil {
		return nil, nil
	}

	// Unmarshal the byte array into a block data using CBOR.
	if err := cbor.Unmarshal(data, &blockData); err != nil {
		return nil, fmt.Errorf("failed to deserialize block data: %v", err)
	}
	return blockData, nil
}
