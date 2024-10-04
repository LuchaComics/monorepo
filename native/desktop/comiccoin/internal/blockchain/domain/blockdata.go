package domain

import (
	"fmt"

	"github.com/fxamacker/cbor/v2"
)

// BlockData represents what can be serialized to disk and over the network.
type BlockData struct {
	Hash   string             `json:"hash"`
	Header *BlockHeader       `json:"block_header"`
	Trans  []BlockTransaction `json:"trans"`
}

type BlockDataRepository interface {
	Upsert(bd *BlockData) error
	GetByHash(hash string) (*BlockData, error)
	ListAll() ([]*BlockData, error)
	DeleteByHash(hash string) error
}

func (b *BlockData) Serialize() ([]byte, error) {
	dataBytes, err := cbor.Marshal(b)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize block data: %v", err)
	}
	return dataBytes, nil
}

func NewBlockDataFromDeserialize(data []byte) (*BlockData, error) {
	// Variable we will use to return.
	blockData := &BlockData{}

	// Defensive code: If programmer entered empty bytes then we will
	// return nil deserialization result.
	if data == nil {
		return nil, nil
	}

	if err := cbor.Unmarshal(data, &blockData); err != nil {
		return nil, fmt.Errorf("failed to deserialize block data: %v", err)
	}
	return blockData, nil
}
