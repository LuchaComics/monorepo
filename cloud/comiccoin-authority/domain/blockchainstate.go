package domain

import (
	"context"
	"fmt"

	"github.com/fxamacker/cbor/v2"
)

// BlockchainState represents the entire blockchain state at the current moment
// in time of operation.
type BlockchainState struct {
	// The unique identifier for this blockchain that we are managing the state for.
	ChainID uint16 `bson:"chain_id" json:"chain_id"`

	LatestBlockNumber uint64 `bson:"latest_block_number" json:"latest_block_number"`
	LatestHash        string `bson:"latest_hash" json:"latest_hash"`
	LatestTokenID     uint64 `bson:"latest_token_id" json:"latest_token_id"`

	AccountHashState string `bson:"account_hash_state" json:"account_hash_state"`
	TokenHashState   string `bson:"token_hash_state" json:"token_hash_state"`

	GenesisBlockData *GenesisBlockData `bson:"genesis_block_data" json:"genesis_block_data"`
}

type BlockchainStateRepository interface {
	// Upsert inserts or updates an blockchain state in the repository.
	Upsert(ctx context.Context, acc *BlockchainState) error

	// GetByChainID retrieves an blockchain state by its chain ID.
	GetByChainID(ctx context.Context, chainID uint16) (*BlockchainState, error)

	// GetForMainNet retrieves an blockchain state by the MainNet chain.
	GetForMainNet(ctx context.Context) (*BlockchainState, error)

	// GetForTestNet retrieves an blockchain state by the TestNet chain.
	GetForTestNet(ctx context.Context) (*BlockchainState, error)

	// ListAll retrieves all blockchain states in the repository.
	ListAll(ctx context.Context) ([]*BlockchainState, error)

	// DeleteByChainID deletes an blockchain state by its chain ID.
	DeleteByChainID(ctx context.Context, chainID uint16) error
}

// Serialize serializes the blockchain state into a byte slice.
// This method uses the cbor library to marshal the blockchainState into a byte slice.
func (b *BlockchainState) Serialize() ([]byte, error) {
	// Marshal the blockchainState into a byte slice using the cbor library.
	dataBytes, err := cbor.Marshal(b)
	if err != nil {
		// Return an error if the marshaling fails.
		return nil, fmt.Errorf("failed to serialize blockchainState: %v", err)
	}
	return dataBytes, nil
}

// NewBlockchainStateFromDeserialize deserializes an blockchainState from a byte slice.
// This method uses the cbor library to unmarshal the byte slice into an blockchainState.
func NewBlockchainStateFromDeserialize(data []byte) (*BlockchainState, error) {
	// Create a new blockchainState variable to return.
	blockchainState := &BlockchainState{}

	// Defensive code: If the input data is empty, return a nil deserialization result.
	if data == nil {
		return nil, nil
	}

	// Unmarshal the byte slice into the blockchainState variable using the cbor library.
	if err := cbor.Unmarshal(data, &blockchainState); err != nil {
		// Return an error if the unmarshaling fails.
		return nil, fmt.Errorf("failed to deserialize blockchainState: %v", err)
	}
	return blockchainState, nil
}
