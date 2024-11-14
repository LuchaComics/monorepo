package domain

import (
	"context"
	"fmt"

	"github.com/fxamacker/cbor/v2"
)

// SignedIssuedTokenDTO represents what can be serialized to disk and over the network.
type SignedIssuedTokenDTO SignedIssuedToken

type SignedIssuedTokenDTORepository interface {
	BroadcastToP2PNetwork(ctx context.Context, dto *SignedIssuedTokenDTO) error
	ReceiveFromP2PNetwork(ctx context.Context) (*SignedIssuedTokenDTO, error)
}

func (b *SignedIssuedTokenDTO) Serialize() ([]byte, error) {
	dataBytes, err := cbor.Marshal(b)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize signed issued token dto: %v", err)
	}
	return dataBytes, nil
}

func NewSignedIssuedTokenDTOFromDeserialize(data []byte) (*SignedIssuedTokenDTO, error) {
	// Variable we will use to return.
	signedIssuedTok := &SignedIssuedTokenDTO{}

	// Defensive code: If programmer entered empty bytes then we will
	// return nil deserialization result.
	if data == nil {
		return nil, nil
	}

	if err := cbor.Unmarshal(data, &signedIssuedTok); err != nil {
		return nil, fmt.Errorf("failed to deserialize signed issued token dto: %v", err)
	}
	return signedIssuedTok, nil
}
