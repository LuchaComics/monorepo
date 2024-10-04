package blockdatadto

import (
	"fmt"

	"github.com/fxamacker/cbor/v2"
	"github.com/libp2p/go-libp2p/core/peer"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/domain"
)

type BlockDataDTORequest struct {
	// Value set by the receiving node, not the sender in the payload!
	FromPeerID peer.ID `json:"from_peer_id"`

	ParamHash string `json:"param_hash"`
}

type BlockDataDTOResponse struct {
	Payload *domain.BlockDataDTO `json:"payload"`
}

func (b *BlockDataDTORequest) Serialize() ([]byte, error) {
	bytes, err := cbor.Marshal(b)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize stream message dto: %v", err)
	}
	return bytes, nil
}

func NewBlockDataDTORequestFromDeserialize(data []byte) (*BlockDataDTORequest, error) {
	// Variable we will use to return.
	dto := &BlockDataDTORequest{}

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

func (b *BlockDataDTOResponse) Serialize() ([]byte, error) {
	bytes, err := cbor.Marshal(b)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize stream message dto: %v", err)
	}
	return bytes, nil
}

func NewBlockDataDTOResponseFromDeserialize(data []byte) (*BlockDataDTOResponse, error) {
	// Variable we will use to return.
	dto := &BlockDataDTOResponse{}

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
