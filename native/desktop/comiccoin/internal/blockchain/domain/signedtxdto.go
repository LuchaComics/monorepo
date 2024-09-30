package domain

import (
	"bytes"
	"context"
	"encoding/gob"
	"fmt"
	"math/big"
)

// SignedTransactionDTO is the data-transfer object used to send and receive
// signed transaction via the peer-to-peer (P2P) network.
type SignedTransactionDTO struct {
	Transaction
	V *big.Int `json:"v"` // Ethereum: Recovery identifier, either 29 or 30 with ardanID.
	R *big.Int `json:"r"` // Ethereum: First coordinate of the ECDSA signature.
	S *big.Int `json:"s"` // Ethereum: Second coordinate of the ECDSA signature.
}

type SignedTransactionDTORepository interface {
	Broadcast(ctx context.Context, dto *SignedTransactionDTO) error
	Receive(ctx context.Context) (*SignedTransactionDTO, error)
}

func (dto *SignedTransactionDTO) Serialize() ([]byte, error) {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)
	err := encoder.Encode(dto)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize block data: %v", err)
	}
	return result.Bytes(), nil
}

func NewSignedTransactionDTOFromDeserialize(data []byte) (*SignedTransactionDTO, error) {
	// Variable we will use to return.
	dto := &SignedTransactionDTO{}

	// Defensive code: If programmer entered empty bytes then we will
	// return nil deserialization result.
	if data == nil {
		return nil, nil
	}

	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&dto)
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize block data: %v", err)
	}
	return dto, nil
}
