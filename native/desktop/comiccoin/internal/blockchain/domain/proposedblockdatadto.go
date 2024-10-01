package domain

import (
	"bytes"
	"context"
	"encoding/gob"
	"fmt"
)

// ProposedBlockDataDTO represents the newly minted BlockData which we want to
// distribute accross the blockchain network.
type ProposedBlockDataDTO struct {
	Hash   string             `json:"hash"`
	Header *BlockHeader       `json:"block_header"`
	Trans  []BlockTransaction `json:"trans"`
}

type ProposedBlockDataDTORepository interface {
	BroadcastToP2PNetwork(ctx context.Context, dto *ProposedBlockDataDTO) error
	ReceiveFromP2PNetwork(ctx context.Context) (*ProposedBlockDataDTO, error)
}

func (dto *ProposedBlockDataDTO) Serialize() ([]byte, error) {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)
	err := encoder.Encode(dto)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize block data: %v", err)
	}
	return result.Bytes(), nil
}

func NewProposedBlockDataDTOFromDeserialize(data []byte) (*ProposedBlockDataDTO, error) {
	// Variable we will use to return.
	dto := &ProposedBlockDataDTO{}

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
