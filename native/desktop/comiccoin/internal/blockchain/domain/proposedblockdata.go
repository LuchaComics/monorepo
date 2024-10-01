package domain

import (
	"bytes"
	"encoding/gob"
	"fmt"
)

// ProposedBlockData represents what can be serialized to disk and over the network.
type ProposedBlockData struct {
	Hash   string             `json:"hash"`
	Header *BlockHeader       `json:"block_header"`
	Trans  []BlockTransaction `json:"trans"`
}

func (b *ProposedBlockData) Serialize() ([]byte, error) {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)
	err := encoder.Encode(b)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize proposed block data: %v", err)
	}
	return result.Bytes(), nil
}

func NewProposedBlockDataFromDeserialize(data []byte) (*ProposedBlockData, error) {
	// Variable we will use to return.
	account := &ProposedBlockData{}

	// Defensive code: If programmer entered empty bytes then we will
	// return nil deserialization result.
	if data == nil {
		return nil, nil
	}

	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&account)
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize proposed block data: %v", err)
	}
	return account, nil
}
