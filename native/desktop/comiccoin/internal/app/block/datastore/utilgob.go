package datastore

import (
	"bytes"
	"encoding/gob"
	"fmt"
)

func (b *Block) Serialize() ([]byte, error) {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)
	err := encoder.Encode(b)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize block: %v", err)
	}
	return result.Bytes(), nil
}

func NewBlockFromDeserialize(data []byte) (*Block, error) {
	// Variable we will use to return.
	block := &Block{}

	// Defensive code: If programmer entered empty bytes then we will
	// return nil deserialization result.
	if data == nil {
		return nil, nil
	}

	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&block)
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize block: %v", err)
	}
	return block, nil
}
