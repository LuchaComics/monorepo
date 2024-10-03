package domain

import (
	"bytes"
	"context"
	"encoding/gob"
	"fmt"
	"log"
)

// BlockDataDTO is the data-transfer object used by nodes to send back and forth
// BlockData across the distributed / P2P network for the blockchain.
type BlockDataDTO struct {
	Hash   string             `json:"hash"`
	Header *BlockHeader       `json:"block_header"`
	Trans  []BlockTransaction `json:"trans"`
}

type BlockDataDTORepository interface {
	// Function will add this block data to the distributed / peer-to-peer
	// network for all peers to discover this data and download remotely to
	// their peer.
	UploadToNetwork(ctx context.Context, data *BlockDataDTO) error

	// Function will lookup this hash in the distributed / peer-to-peer
	// network and return the block data.
	DownloadFromNetwork(ctx context.Context, BlockDataHash string) (*BlockDataDTO, error)
}

func (b *BlockDataDTO) Serialize() ([]byte, error) {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)
	err := encoder.Encode(b)
	if err != nil {
		log.Printf("BlockDataDTO: Serialize: result: %v\n", result)
		return nil, fmt.Errorf("failed to serialize block data: %v", err)
	}
	return result.Bytes(), nil
}

func NewBlockDataDTOFromDeserialize(data []byte) (*BlockDataDTO, error) {
	// Variable we will use to return.
	account := &BlockDataDTO{}

	// Defensive code: If programmer entered empty bytes then we will
	// return nil deserialization result.
	if data == nil {
		return nil, nil
	}

	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&account)
	if err != nil {
		log.Printf("BlockDataDTO: NewBlockDataDTOFromDeserialize: data: %v\n", data)
		return nil, fmt.Errorf("failed to deserialize block data: %v", err)
	}
	return account, nil
}
