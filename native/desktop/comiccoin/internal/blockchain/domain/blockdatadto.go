package domain

import (
	"bytes"
	"context"
	"encoding/gob"
	"fmt"
)

// BlockDataDTO is the data-transfer object used by nodes to send back and forth
// BlockData across the distributed / P2P network for the blockchain.
type BlockDataDTO struct {
	Hash   string             `json:"hash"`
	Header *BlockHeader       `json:"block_header"`
	Trans  []BlockTransaction `json:"trans"`
}

type BlockDataDTORepository interface {
	// // Function will randomly pick a connected peer and send them a request.
	// SendRequestToRandomPeer(ctx context.Context, hash string) error
	//
	// // Function will block your current execution and wait until it receives
	// // any request from the peer-to-peer network. Function will return the
	// // `peerID` that sent the request and the hash value.
	// ReceiveRequestFromNetwork(ctx context.Context) (string, string, error)
	//
	// // Function will send sync data to the peer whom requested block data.
	// SendResponseToPeer(ctx context.Context, peerID string, data []*BlockDataDTO) error
	//
	// ReceiveResponseFromNetwork(ctx context.Context) ([]*BlockDataDTO, error)

	ListLatestAfterHash(ctx context.Context, afterBlockDataHash string) ([]*BlockDataDTO, error)
}

func (b *BlockDataDTO) Serialize() ([]byte, error) {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)
	err := encoder.Encode(b)
	if err != nil {
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
		return nil, fmt.Errorf("failed to deserialize block data: %v", err)
	}
	return account, nil
}
