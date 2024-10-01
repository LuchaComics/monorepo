package domain

import (
	"bytes"
	"context"
	"encoding/gob"
	"fmt"
)

// PurposedBlockDataDTO represents the newly minted BlockData which we want to
// distribute accross the blockchain network.
type PurposedBlockDataDTO struct {
	Hash   string             `json:"hash"`
	Header *BlockHeader       `json:"block_header"`
	Trans  []BlockTransaction `json:"trans"`
}

type PurposedBlockDataDTORepository interface {
	BroadcastToP2PNetwork(ctx context.Context, dto *PurposedBlockDataDTO) error
	ReceiveFromP2PNetwork(ctx context.Context) (*PurposedBlockDataDTO, error)
}

func (dto *PurposedBlockDataDTO) Serialize() ([]byte, error) {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)
	err := encoder.Encode(dto)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize block data: %v", err)
	}
	return result.Bytes(), nil
}

func NewPurposedBlockDataDTOFromDeserialize(data []byte) (*PurposedBlockDataDTO, error) {
	// Variable we will use to return.
	dto := &PurposedBlockDataDTO{}

	// Defensive code: If programmer entered empty bytes then we will
	// return nil deserialization result.
	if data == nil {
		return nil, nil
	}

	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&dto)
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize block data: %v", err)
	}
	return dto, nil
}
