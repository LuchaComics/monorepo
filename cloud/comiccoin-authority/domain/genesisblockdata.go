package domain

import (
	"context"
	"fmt"

	"github.com/fxamacker/cbor/v2"
)

// GenesisBlockData represents the first block (data) in our blockchain.
type GenesisBlockData BlockData

// BlockDataToGenesisBlockData method converts a `BlockData` data type into
// a `GenesisBlockData` data type.
func BlockDataToGenesisBlockData(bd *BlockData) *GenesisBlockData {
	return (*GenesisBlockData)(bd)
}

// GenesisBlockDataRepository is an interface that defines the methods for
// handling the Genesis Block Data in our local database.
type GenesisBlockDataRepository interface {
	GetByChainID(ctx context.Context, chainID uint16) (*GenesisBlockData, error)
	UpsertByChainID(ctx context.Context, genesis *GenesisBlockData) error

	OpenTransaction() error
	CommitTransaction() error
	DiscardTransaction()
}

// Serialize serializes a genesis block data into a byte array.
// It returns the serialized byte array and an error if one occurs.
func (b *GenesisBlockData) Serialize() ([]byte, error) {
	// Marshal the block data into a byte array using CBOR.
	dataBytes, err := cbor.Marshal(b)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize genesis block data: %v", err)
	}
	return dataBytes, nil
}

// NewGenesisBlockDataFromDeserialize deserializes a genesis block data from a
// byte array. It returns the deserialized block data and an error if one occurs.
func NewGenesisBlockDataFromDeserialize(data []byte) (*GenesisBlockData, error) {
	// Variable we will use to return.
	blockData := &GenesisBlockData{}

	// Defensive code: If programmer entered empty bytes then we will
	// return nil deserialization result.
	if data == nil {
		return nil, nil
	}

	// Unmarshal the byte array into a block data using CBOR.
	if err := cbor.Unmarshal(data, &blockData); err != nil {
		return nil, fmt.Errorf("failed to deserialize genesis block data: %v", err)
	}
	return blockData, nil
}
