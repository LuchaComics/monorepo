package simple

import (
	"bytes"
	"encoding/gob"
	"fmt"

	"github.com/libp2p/go-libp2p/core/peer"
)

type SimpleMessageRequest struct {
	ID      string `json:"id"`
	Content []byte `json:"content"`

	// Value set by the receiving node, not the sender in the payload!
	FromPeerID peer.ID `json:"from_peer_id"`
}

type SimpleMessageResponse struct {
	ID      string `json:"id"`
	Content []byte `json:"content"`

	// Value set by the receiving node, not the sender in the payload!
	FromPeerID peer.ID `json:"from_peer_id"`
}

func (b *SimpleMessageRequest) Serialize() ([]byte, error) {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)
	err := encoder.Encode(b)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize stream message dto: %v", err)
	}
	return result.Bytes(), nil
}

func NewSimpleMessageRequestFromDeserialize(data []byte) (*SimpleMessageRequest, error) {
	// Variable we will use to return.
	dto := &SimpleMessageRequest{}

	// Defensive code: If programmer entered empty bytes then we will
	// return nil deserialization result.
	if data == nil {
		return nil, nil
	}

	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&dto)
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize stream message dto: %v", err)
	}
	return dto, nil
}

func (b *SimpleMessageResponse) Serialize() ([]byte, error) {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)
	err := encoder.Encode(b)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize stream message dto: %v", err)
	}
	return result.Bytes(), nil
}

func NewSimpleMessageResponseFromDeserialize(data []byte) (*SimpleMessageResponse, error) {
	// Variable we will use to return.
	dto := &SimpleMessageResponse{}

	// Defensive code: If programmer entered empty bytes then we will
	// return nil deserialization result.
	if data == nil {
		return nil, nil
	}

	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&dto)
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize stream message dto: %v", err)
	}
	return dto, nil
}
