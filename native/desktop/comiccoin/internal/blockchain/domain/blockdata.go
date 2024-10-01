package domain

import (
	"bytes"
	"encoding/gob"
	"fmt"
)

// BlockData represents what can be serialized to disk and over the network.
type BlockData struct {
	Hash   string              `json:"hash"`
	Header *BlockHeader        `json:"block_header"`
	Trans  []*BlockTransaction `json:"trans"`
}

type BlockDataRepository interface {
	Upsert(bd *BlockData) error
	GetByHash(hash string) (*BlockData, error)
	ListAll() ([]*BlockData, error)
	DeleteByHash(hash string) error
}

func (b *BlockData) Serialize() ([]byte, error) {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)
	err := encoder.Encode(b)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize block data: %v", err)
	}
	return result.Bytes(), nil
}

func NewBlockDataFromDeserialize(data []byte) (*BlockData, error) {
	// Variable we will use to return.
	account := &BlockData{}

	// Defensive code: If programmer entered empty bytes then we will
	// return nil deserialization result.
	if data == nil {
		return nil, nil
	}

	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&account)
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize block data: %v", err)
	}
	return account, nil
}
