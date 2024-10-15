package domain

import (
	"context"
	"fmt"

	"github.com/fxamacker/cbor/v2"
)

// ProposedBlockDataDTO represents the newly minted BlockData which we want to
// distribute accross the blockchain network.
type ProposedBlockDataDTO struct {
	Hash            string             `json:"hash"`
	Header          *BlockHeader       `json:"block_header"`
	HeaderSignature []byte             `json:"header_signature"`
	Trans           []BlockTransaction `json:"trans"`
	Validator       *Validator         `json:"validator"`
}

type ProposedBlockDataDTORepository interface {
	BroadcastToP2PNetwork(ctx context.Context, dto *ProposedBlockDataDTO) error
	ReceiveFromP2PNetwork(ctx context.Context) (*ProposedBlockDataDTO, error)
}

func (dto *ProposedBlockDataDTO) Serialize() ([]byte, error) {
	dataBytes, err := cbor.Marshal(dto)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize proposed block data dto: %v", err)
	}
	return dataBytes, nil
}

func NewProposedBlockDataDTOFromDeserialize(data []byte) (*ProposedBlockDataDTO, error) {
	// Variable we will use to return.
	dto := &ProposedBlockDataDTO{}

	// Defensive code: If programmer entered empty bytes then we will
	// return nil deserialization result.
	if data == nil {
		return nil, nil
	}

	if err := cbor.Unmarshal(data, &dto); err != nil {
		return nil, fmt.Errorf("failed to deserialize proposed block data dto: %v", err)
	}
	return dto, nil
}
