package datastore

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
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
	block := &Block{}
	err := json.Unmarshal(data, block)
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize block: %v", err)
	}
	return block, nil
}
