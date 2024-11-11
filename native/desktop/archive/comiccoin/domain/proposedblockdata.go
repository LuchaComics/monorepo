package domain

import (
	"fmt"

	"github.com/fxamacker/cbor/v2"
)

// ProposedBlockData represents what can be serialized to disk and over the network.
type ProposedBlockData struct {
	Hash            string             `json:"hash"`
	Header          *BlockHeader       `json:"block_header"`
	HeaderSignatureBytes []byte             `json:"header_signature_bytes"`
	Trans           []BlockTransaction `json:"trans"`
	Validator       *Validator         `json:"validator"`
}

func (b *ProposedBlockData) Serialize() ([]byte, error) {
	dataBytes, err := cbor.Marshal(b)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize proposed block data: %v", err)
	}
	return dataBytes, nil
}

func NewProposedBlockDataFromDeserialize(data []byte) (*ProposedBlockData, error) {
	// Variable we will use to return.
	proposedBlockData := &ProposedBlockData{}

	// Defensive code: If programmer entered empty bytes then we will
	// return nil deserialization result.
	if data == nil {
		return nil, nil
	}

	if err := cbor.Unmarshal(data, &proposedBlockData); err != nil {
		return nil, fmt.Errorf("failed to deserialize proposed block data: %v", err)
	}
	return proposedBlockData, nil
}
