package domain

// BlockchainLastestHashRepository is an interface that defines the methods for interacting with the latest hash of the blockchain.
// It provides methods for getting and setting the latest hash of the blockchain.
type BlockchainLastestHashRepository interface {
	// Get returns the latest hash of the blockchain.
	// It returns the latest hash as a string and an error if one occurs.
	Get() (string, error)

	// Set sets the latest hash of the blockchain.
	// It takes a hash as a string and returns an error if one occurs.
	Set(hash string) error
}
