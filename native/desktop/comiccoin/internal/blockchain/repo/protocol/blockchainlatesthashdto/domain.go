package blockchainlatesthashdto

import (
	"fmt"

	"github.com/fxamacker/cbor/v2"
	"github.com/libp2p/go-libp2p/core/peer"
)

type BlockchainLatestHashDTORequest struct {
	Content []byte `json:"content"`

	// Value set by the receiving node, not the sender in the payload!
	FromPeerID peer.ID `json:"from_peer_id"`
}

type BlockchainLatestHashDTOResponse struct {
	Content []byte `json:"content"`

	// Value set by the receiving node, not the sender in the payload!
	FromPeerID peer.ID `json:"from_peer_id"`
}

func (b *BlockchainLatestHashDTORequest) Serialize() ([]byte, error) {
	bytes, err := cbor.Marshal(b)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize stream message dto: %v", err)
	}
	return bytes, nil
}

func NewBlockchainLatestHashDTORequestFromDeserialize(data []byte) (*BlockchainLatestHashDTORequest, error) {
	// Variable we will use to return.
	dto := &BlockchainLatestHashDTORequest{}

	// Defensive code: If programmer entered empty bytes then we will
	// return nil deserialization result.
	if data == nil {
		return nil, nil
	}
	err := cbor.Unmarshal(data, &dto)
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize stream message dto: %v", err)
	}
	return dto, nil
}

func (b *BlockchainLatestHashDTOResponse) Serialize() ([]byte, error) {
	bytes, err := cbor.Marshal(b)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize stream message dto: %v", err)
	}
	return bytes, nil
}

func NewBlockchainLatestHashDTOResponseFromDeserialize(data []byte) (*BlockchainLatestHashDTOResponse, error) {
	// Variable we will use to return.
	dto := &BlockchainLatestHashDTOResponse{}

	// Defensive code: If programmer entered empty bytes then we will
	// return nil deserialization result.
	if data == nil {
		return nil, nil
	}
	err := cbor.Unmarshal(data, &dto)
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize stream message dto: %v", err)
	}
	return dto, nil
}
