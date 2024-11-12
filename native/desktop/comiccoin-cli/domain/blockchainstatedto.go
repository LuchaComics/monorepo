package domain

import (
	"context"
	"fmt"

	"github.com/fxamacker/cbor/v2"
)

// BlockchainStateDTO represents the data that can be serialized to disk and over the network.
type BlockchainStateDTO BlockchainState

type BlockchainStateDTORepository interface {
	GetFromCentralAuthorityByChainID(ctx context.Context, chainID uint16) (*BlockchainStateDTO, error)
}

// Serialize serializes a block data into a byte array.
// It returns the serialized byte array and an error if one occurs.
func (b *BlockchainStateDTO) Serialize() ([]byte, error) {
	// Marshal the block data into a byte array using CBOR.
	dataBytes, err := cbor.Marshal(b)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize block data: %v", err)
	}
	return dataBytes, nil
}

// NewBlockchainStateDTOFromDeserialize deserializes a block data from a byte array.
// It returns the deserialized block data and an error if one occurs.
func NewBlockchainStateDTOFromDeserialize(data []byte) (*BlockchainStateDTO, error) {
	// Variable we will use to return.
	blockData := &BlockchainStateDTO{}

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
