package domain

import (
	"context"
	"fmt"

	"github.com/fxamacker/cbor/v2"
)

// MempoolTransaction represents the data that can be serialized to disk and over the network.
type MempoolTransactionDTO MempoolTransaction

// ToDTO method converts a `MempoolTransaction` data type into a `MempoolTransactionDTO` data type.
func (bd *MempoolTransaction) ToDTO() *MempoolTransactionDTO {
	return (*MempoolTransactionDTO)(bd)
}

func (bd *MempoolTransactionDTO) ToIDO() *MempoolTransaction {
	return (*MempoolTransaction)(bd)
}

// MempoolTransactionRepository is an interface that defines the methods for
// handling the blockchain state via the network.
type MempoolTransactionDTORepository interface {
	SubmitToBlockchainAuthority(ctx context.Context, dto *MempoolTransactionDTO) error
}

// Serialize serializes golang struct data into a byte array.
// It returns the serialized byte array and an error if one occurs.
func (b *MempoolTransactionDTO) Serialize() ([]byte, error) {
	// Marshal the block data into a byte array using CBOR.
	dataBytes, err := cbor.Marshal(b)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize mempool transaction dto: %v", err)
	}
	return dataBytes, nil
}

// NewMempoolTransactionDTOFromDeserialize deserializes struct data from a
// byte array. It returns the deserialized struct data and an error if one occurs.
func NewMempoolTransactionDTOFromDeserialize(data []byte) (*MempoolTransactionDTO, error) {
	// Variable we will use to return.
	blockData := &MempoolTransactionDTO{}

	// Defensive code: If programmer entered empty bytes then we will
	// return nil deserialization result.
	if data == nil {
		return nil, nil
	}

	// Unmarshal the byte array into a block data using CBOR.
	if err := cbor.Unmarshal(data, &blockData); err != nil {
		return nil, fmt.Errorf("failed to deserialize mempool transaction dto: %v", err)
	}
	return blockData, nil
}
