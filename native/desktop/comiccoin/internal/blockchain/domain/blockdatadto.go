package domain

import (
	"context"
	"fmt"

	"github.com/fxamacker/cbor/v2"
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
	SendResponseToPeer(ctx context.Context, peerID peer.ID, data *BlockDataDTO) error
	ReceiveResponseFromNetwork(ctx context.Context) (*BlockDataDTO, error)
}

func (b *BlockDataDTO) Serialize() ([]byte, error) {
	dataBytes, err := cbor.Marshal(b)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize block data dto: %v", err)
	}
	return dataBytes, nil
}

func NewBlockDataDTOFromDeserialize(data []byte) (*BlockDataDTO, error) {
	// Variable we will use to return.
	blockDataDTO := &BlockDataDTO{}

	// Defensive code: If programmer entered empty bytes then we will
	// return nil deserialization result.
	if data == nil {
		return nil, nil
	}

	if err := cbor.Unmarshal(data, &blockDataDTO); err != nil {
		return nil, fmt.Errorf("failed to deserialize account: %v", err)
	}
	return blockDataDTO, nil
}
