package domain

import (
	"context"
	"fmt"

	"github.com/fxamacker/cbor/v2"
)

// BlockchainState represents the first block (data) in our blockchain.
type BlockchainState struct {
	// The unique identifier for this blockchain that we are managing the state for.
	ChainID uint16 `bson:"chain_id" json:"chain_id"`

	LatestBlockNumber uint64 `bson:"latest_block_number" json:"latest_block_number"`

	LatestHash string `bson:"latest_hash" json:"latest_hash"`

	LatestTokenID uint64 `bson:"latest_token_id" json:"latest_token_id"`

	AccountHashState string `bson:"account_hash_state" json:"account_hash_state"`

	TokenHashState string `bson:"token_hash_state" json:"token_hash_state"`
}

// BlockchainStateRepository is an interface that defines the methods for
// loading up the Genesis block from file.
type BlockchainStateRepository interface {
	GetByChainID(ctx context.Context, chainID uint16) (*BlockchainState, error)
	UpsertByChainID(ctx context.Context, genesis *BlockchainState) error
}

// Serialize serializes a block data into a byte array.
// It returns the serialized byte array and an error if one occurs.
func (b *BlockchainState) Serialize() ([]byte, error) {
	// Marshal the block data into a byte array using CBOR.
	dataBytes, err := cbor.Marshal(b)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize block data: %v", err)
	}
	return dataBytes, nil
}

// NewBlockchainStateFromDeserialize deserializes a block data from a byte array.
// It returns the deserialized block data and an error if one occurs.
func NewBlockchainStateFromDeserialize(data []byte) (*BlockchainState, error) {
	// Variable we will use to return.
	blockData := &BlockchainState{}

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
