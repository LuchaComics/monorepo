package domain

import (
	"context"
	"fmt"

	"github.com/fxamacker/cbor/v2"
)

// BlockDataDTO is the data-transfer object used by nodes to send back and forth
type BlockDataDTO BlockData

type BlockNumberByHashDTO struct {
	Number uint64 `bson:"number"`
	Hash   string `bson:"hash"`
}

// BlockDataDTORepository is an interface that defines the methods for interacting with block data DTOs.
type BlockDataDTORepository interface {
	GetFromCentralAuthorityByHash(ctx context.Context, hash string) (*BlockDataDTO, error)
	GetFromCentralAuthorityByBlockNumber(ctx context.Context, blockNumber uint64) (*BlockDataDTO, error)
}

// Serialize serializes a block data DTO into a byte array.
// It returns the serialized byte array and an error if one occurs.
func (b *BlockDataDTO) Serialize() ([]byte, error) {
	// Marshal the block data DTO into a byte array using CBOR.
	dataBytes, err := cbor.Marshal(b)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize block data dto: %v", err)
	}
	return dataBytes, nil
}

// NewBlockDataDTOFromDeserialize deserializes a block data DTO from a byte array.
// It returns the deserialized block data DTO and an error if one occurs.
func NewBlockDataDTOFromDeserialize(data []byte) (*BlockDataDTO, error) {
	// Variable we will use to return.
	blockDataDTO := &BlockDataDTO{}

	// Defensive code: If programmer entered empty bytes then we will
	// return nil deserialization result.
	if data == nil {
		return nil, nil
	}

	// Unmarshal the byte array into a block data DTO using CBOR.
	if err := cbor.Unmarshal(data, &blockDataDTO); err != nil {
		return nil, fmt.Errorf("failed to deserialize block data dto: %v", err)
	}
	return blockDataDTO, nil
}
