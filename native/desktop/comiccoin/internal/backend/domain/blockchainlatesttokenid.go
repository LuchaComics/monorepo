package domain

// BlockchainLastestTokenIDRepository is an interface that defines the methods for interacting with the latest token ID of the blockchain.
// It provides methods for getting and setting the latest token ID of the blockchain.
type BlockchainLastestTokenIDRepository interface {
	// Get returns the latest token ID of the blockchain.
	// It returns the latest token ID as a string and an error if one occurs.
	Get() (uint64, error)

	// Set sets the latest token ID of the blockchain.
	// It takes a token ID as a string and returns an error if one occurs.
	Set(tokenID uint64) error
}
