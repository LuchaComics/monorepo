package datastore

import (
	"bytes"
	"encoding/gob"
	"fmt"
)

func (b *BlockData) Serialize() ([]byte, error) {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)
	err := encoder.Encode(b)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize blockdata: %v", err)
	}
	return result.Bytes(), nil
}

func NewBlockDataFromDeserialize(data []byte) (*BlockData, error) {
	// Variable we will use to return.
	blockdata := &BlockData{}

	// Defensive code: If programmer entered empty bytes then we will
	// return nil deserialization result.
	if data == nil {
		return nil, nil
	}

	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&blockdata)
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize blockdata: %v", err)
	}
	return blockdata, nil
}
