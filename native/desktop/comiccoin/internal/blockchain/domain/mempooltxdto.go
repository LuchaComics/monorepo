package domain

import (
	"context"
	"fmt"
	"math/big"

	"github.com/fxamacker/cbor/v2"
)

// MempoolTransactionDTO is the data-transfer object used to transmit signed transactions
// between nodes in the peer-to-peer (P2P) network.
// It contains the transaction data, as well as the ECDSA signature and recovery identifier.
type MempoolTransactionDTO struct {
	// The transaction data, including sender, recipient, amount, and other metadata.
	Transaction

	// The recovery identifier, either 29 or 30, depending on the Ethereum network.
	V *big.Int `json:"v"`

	// The first coordinate of the ECDSA signature.
	R *big.Int `json:"r"`

	// The second coordinate of the ECDSA signature.
	S *big.Int `json:"s"`
}

// MempoolTransactionDTORepository interface defines the methods for interacting with
// the mempool transaction DTO repository.
// This interface provides a way to broadcast and receive mempool transactions over the P2P network.
type MempoolTransactionDTORepository interface {
	// BroadcastToP2PNetwork sends the mempool transaction DTO to the P2P network.
	BroadcastToP2PNetwork(ctx context.Context, dto *MempoolTransactionDTO) error

	// ReceiveFromP2PNetwork receives a mempool transaction DTO from the P2P network.
	ReceiveFromP2PNetwork(ctx context.Context) (*MempoolTransactionDTO, error)
}

// Serialize serializes the mempool transaction DTO into a byte slice.
// This method uses the cbor library to marshal the DTO into a byte slice.
func (dto *MempoolTransactionDTO) Serialize() ([]byte, error) {
	// Marshal the DTO into a byte slice using the cbor library.
	dataBytes, err := cbor.Marshal(dto)
	if err != nil {
		// Return an error if the marshaling fails.
		return nil, fmt.Errorf("failed to serialize mempool transaction dto: %v", err)
	}
	return dataBytes, nil
}

// NewMempoolTransactionDTOFromDeserialize deserializes a mempool transaction DTO from a byte slice.
// This method uses the cbor library to unmarshal the byte slice into a DTO.
func NewMempoolTransactionDTOFromDeserialize(data []byte) (*MempoolTransactionDTO, error) {
	// Create a new DTO variable to return.
	dto := &MempoolTransactionDTO{}

	// Defensive code: If the input data is empty, return a nil deserialization result.
	if data == nil {
		return nil, nil
	}

	// Unmarshal the byte slice into the DTO variable using the cbor library.
	if err := cbor.Unmarshal(data, &dto); err != nil {
		// Return an error if the unmarshaling fails.
		return nil, fmt.Errorf("failed to deserialize mempool transaction dto: %v", err)
	}
	return dto, nil
}
