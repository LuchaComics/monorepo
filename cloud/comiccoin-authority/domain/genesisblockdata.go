package domain

import "context"

// GenesisBlockData represents the first block (data) in our blockchain.
type GenesisBlockData BlockData

// BlockDataToGenesisBlockData method converts a `BlockData` data type into
// a `GenesisBlockData` data type.
func BlockDataToGenesisBlockData(bd *BlockData) *GenesisBlockData {
	return (*GenesisBlockData)(bd)
}

// GenesisBlockDataRepository is an interface that defines the methods for
// loading up the Genesis block from file.
type GenesisBlockDataRepository interface {
	GetByChainID(ctx context.Context, chainID uint16) (*GenesisBlockData, error)
	UpsertByChainID(ctx context.Context, genesis *GenesisBlockData) error
}
