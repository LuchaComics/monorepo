package domain

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/fxamacker/cbor/v2"
)

// Account struct represents a user's wallet, which contains a private key and a public key.
// The private key is used to sign transactions, while the public key acts as the user's address.
// Note that the keys are stored encrypted at rest and require the user's password to decrypt.
type Account struct {
	// Unique identifier for the account.
	ID string `json:"id"`

	// The balance of the account in coins.
	Balance uint64 `json:"balance"`

	// The file path where the wallet is stored.
	WalletFilepath string `json:"wallet_filepath"`

	// The public address of the account.
	Address *common.Address `json:"address"`
}

// AccountRepository interface defines the methods for interacting with the account repository.
// This interface provides a way to manage accounts, including upserting, getting, listing, and deleting.
type AccountRepository interface {
	// Upsert inserts or updates an account in the repository.
	Upsert(acc *Account) error

	// GetByID retrieves an account by its ID.
	GetByID(id string) (*Account, error)

	// ListAll retrieves all accounts in the repository.
	ListAll() ([]*Account, error)

	// DeleteByID deletes an account by its ID.
	DeleteByID(id string) error

	// HashState returns a hash based on the contents of the accounts and
	// their balances. This is added to each block and checked by peers.
	HashState() (string, error)
}

// Serialize serializes the account into a byte slice.
// This method uses the cbor library to marshal the account into a byte slice.
func (b *Account) Serialize() ([]byte, error) {
	// Marshal the account into a byte slice using the cbor library.
	dataBytes, err := cbor.Marshal(b)
	if err != nil {
		// Return an error if the marshaling fails.
		return nil, fmt.Errorf("failed to serialize account: %v", err)
	}
	return dataBytes, nil
}

// NewAccountFromDeserialize deserializes an account from a byte slice.
// This method uses the cbor library to unmarshal the byte slice into an account.
func NewAccountFromDeserialize(data []byte) (*Account, error) {
	// Create a new account variable to return.
	account := &Account{}

	// Defensive code: If the input data is empty, return a nil deserialization result.
	if data == nil {
		return nil, nil
	}

	// Unmarshal the byte slice into the account variable using the cbor library.
	if err := cbor.Unmarshal(data, &account); err != nil {
		// Return an error if the unmarshaling fails.
		return nil, fmt.Errorf("failed to deserialize account: %v", err)
	}
	return account, nil
}
