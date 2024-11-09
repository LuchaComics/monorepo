package domain

// GenesisBlockData represents the first block (data) in our blockchain.
type GenesisBlockData BlockData

// GenesisBlockDataRepository is an interface that defines the methods for
// loading up the Genesis block from file.
type GenesisBlockDataRepository interface {
	// LoadGenesisData method returns the Genesis block of this blockchain.
	LoadGenesisData() (*GenesisBlockData, error)
}
