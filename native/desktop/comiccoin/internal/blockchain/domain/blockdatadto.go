package domain

import (
	"bytes"
	"context"
	"encoding/gob"
	"fmt"
	"log"

	"github.com/libp2p/go-libp2p/core/peer"
)

// BlockDataDTO is the data-transfer object used by nodes to send back and forth
// BlockData across the distributed / P2P network for the blockchain.
type BlockDataDTO struct {
	Hash   string             `json:"hash"`
	Header *BlockHeader       `json:"block_header"`
	Trans  []BlockTransaction `json:"trans"`
}

type BlockDataDTORepository interface {
	SendRequestToRandomPeer(ctx context.Context, blockDataHash string) error
	ReceiveRequestFromNetwork(ctx context.Context) (peer.ID, string, error)
	SendResponseToPeer(ctx context.Context, peerID peer.ID, data BlockDataDTO) error
	ReceiveResponseFromNetwork(ctx context.Context) (*BlockDataDTO, error)
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
