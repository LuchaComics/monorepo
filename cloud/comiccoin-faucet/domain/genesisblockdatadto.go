package domain

import (
	"context"
	"fmt"

	"github.com/fxamacker/cbor/v2"
)

// GenesisBlockData represents the data that can be serialized to disk and over the network.
type GenesisBlockDataDTO GenesisBlockData

// BlockDataToGenesisBlockData method converts a `GenesisBlockData` data type into
// a `GenesisBlockDataDTO` data type.
func GenesisBlockDataToGenesisBlockDataDTO(bd *GenesisBlockData) *GenesisBlockDataDTO {
	return (*GenesisBlockDataDTO)(bd)
}

func GenesisBlockDataDTOToGenesisBlockData(bd *GenesisBlockDataDTO) *GenesisBlockData {
	return (*GenesisBlockData)(bd)
}

// GenesisBlockDataRepository is an interface that defines the methods for
// handling the Genesis Block Data via the network.
type GenesisBlockDataDTORepository interface {
	GetFromBlockchainAuthorityByChainID(ctx context.Context, chainID uint16) (*GenesisBlockDataDTO, error)
}

// Serialize serializes a genesis block data into a byte array.
// It returns the serialized byte array and an error if one occurs.
func (b *GenesisBlockDataDTO) Serialize() ([]byte, error) {
	// Marshal the block data into a byte array using CBOR.
	dataBytes, err := cbor.Marshal(b)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize genesis block data dto: %v", err)
	}
	return dataBytes, nil
}

// NewGenesisBlockDataDTOFromDeserialize deserializes a genesis block data from a
// byte array. It returns the deserialized block data and an error if one occurs.
func NewGenesisBlockDataDTOFromDeserialize(data []byte) (*GenesisBlockDataDTO, error) {
	// Variable we will use to return.
	blockData := &GenesisBlockDataDTO{}

	// Defensive code: If programmer entered empty bytes then we will
	// return nil deserialization result.
	if data == nil {
		return nil, nil
	}

	// Unmarshal the byte array into a block data using CBOR.
	if err := cbor.Unmarshal(data, &blockData); err != nil {
		return nil, fmt.Errorf("failed to deserialize genesis block data dto: %v", err)
	}
	return blockData, nil
}
