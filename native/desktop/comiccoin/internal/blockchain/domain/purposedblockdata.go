package domain

import (
	"bytes"
	"encoding/gob"
	"fmt"
)

// PurposedBlockData represents what can be serialized to disk and over the network.
type PurposedBlockData struct {
	Hash   string             `json:"hash"`
	Header *BlockHeader       `json:"block_header"`
	Trans  []BlockTransaction `json:"trans"`
}

func (b *PurposedBlockData) Serialize() ([]byte, error) {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)
	err := encoder.Encode(b)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize purposed block data: %v", err)
	}
	return result.Bytes(), nil
}

func NewPurposedBlockDataFromDeserialize(data []byte) (*PurposedBlockData, error) {
	// Variable we will use to return.
	account := &PurposedBlockData{}

	// Defensive code: If programmer entered empty bytes then we will
	// return nil deserialization result.
	if data == nil {
		return nil, nil
	}

	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&account)
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize purposed block data: %v", err)
	}
	return account, nil
}
