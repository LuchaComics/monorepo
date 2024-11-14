package domain

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/fxamacker/cbor/v2"
)

type Wallet struct {
	// The (Optional) description for this wallet.
	Label string `json:"label"`

	// The public address of the wallet.
	Address *common.Address `json:"address"`

	// The file path where the wallet is stored.
	Filepath string `json:"filepath"`
}

type WalletRepository interface {
	// Upsert inserts or updates an wallet in the repository.
	Upsert(acc *Wallet) error

	// GetByID retrieves an wallet by its Address.
	GetByAddress(address *common.Address) (*Wallet, error)

	// ListAll retrieves all wallets in the repository.
	ListAll() ([]*Wallet, error)

	// DeleteByID deletes an wallet by its Address.
	DeleteByAddress(address *common.Address) error

	OpenTransaction() error
	CommitTransaction() error
	DiscardTransaction()
}

// Serialize serializes the wallet into a byte slice.
// This method uses the cbor library to marshal the wallet into a byte slice.
func (b *Wallet) Serialize() ([]byte, error) {
	// Marshal the wallet into a byte slice using the cbor library.
	dataBytes, err := cbor.Marshal(b)
	if err != nil {
		// Return an error if the marshaling fails.
		return nil, fmt.Errorf("failed to serialize wallet: %v", err)
	}
	return dataBytes, nil
}

// NewWalletFromDeserialize deserializes an wallet from a byte slice.
// This method uses the cbor library to unmarshal the byte slice into an wallet.
func NewWalletFromDeserialize(data []byte) (*Wallet, error) {
	// Create a new wallet variable to return.
	wallet := &Wallet{}

	// Defensive code: If the input data is empty, return a nil deserialization result.
	if data == nil {
		return nil, nil
	}

	// Unmarshal the byte slice into the wallet variable using the cbor library.
	if err := cbor.Unmarshal(data, &wallet); err != nil {
		// Return an error if the unmarshaling fails.
		return nil, fmt.Errorf("failed to deserialize wallet: %v", err)
	}
	return wallet, nil
}