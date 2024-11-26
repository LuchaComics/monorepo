package domain

import (
	"context"
	"fmt"

	"github.com/fxamacker/cbor/v2"
)

// BlockData represents the data that can be serialized to disk and over the network.
type BlockDataDTO BlockData

// BlockDataToBlockData method converts a `BlockData` data type into
// a `BlockDataDTO` data type.
func BlockDataToBlockDataDTO(bd *BlockData) *BlockDataDTO {
	return (*BlockDataDTO)(bd)
}

func BlockDataDTOToBlockData(bd *BlockDataDTO) *BlockData {
	return (*BlockData)(bd)
}

// BlockDataRepository is an interface that defines the methods for
// handling the Block Data via the network.
type BlockDataDTORepository interface {
	GetFromBlockchainAuthorityByHash(ctx context.Context, hash string) (*BlockDataDTO, error)
}

// Serialize serializes a block data into a byte array.
// It returns the serialized byte array and an error if one occurs.
func (b *BlockDataDTO) Serialize() ([]byte, error) {
	// Marshal the block data into a byte array using CBOR.
	dataBytes, err := cbor.Marshal(b)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize block data dto: %v", err)
	}
	return dataBytes, nil
}

// NewBlockDataDTOFromDeserialize deserializes a block data from a
// byte array. It returns the deserialized block data and an error if one occurs.
func NewBlockDataDTOFromDeserialize(data []byte) (*BlockDataDTO, error) {
	// Variable we will use to return.
	blockData := &BlockDataDTO{}

	// Defensive code: If programmer entered empty bytes then we will
	// return nil deserialization result.
	if data == nil {
		return nil, nil
	}

	// Unmarshal the byte array into a block data using CBOR.
	if err := cbor.Unmarshal(data, &blockData); err != nil {
		return nil, fmt.Errorf("failed to deserialize block data dto: %v", err)
	}
	return blockData, nil
}
