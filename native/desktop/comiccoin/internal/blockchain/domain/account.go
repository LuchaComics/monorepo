package domain

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/fxamacker/cbor/v2"
)

// Account struct represents an entity in our blockchain whom has transfered
// between another account some amount of coins or non-fungible tokens.
//
// When you create a wallet it will still not exist on the blockchain, only
// when you use this wallet to receive or send coins to another wallet will
// it become logged on the blockchain.
//
// When an account exists on the blockchain, that means every node in the peer-
// -to-peer network will have this account record in their local storage. This
// is important because every-time a coin is mined, the miner takes a hash of
// the entire accounts database to verify the values in the account are
// consistent across all the peer-to-peer nodes in the distributed network;
// therefore, preventing fraud from occuring.
type Account struct {
	// The public address of the account.
	Address *common.Address `json:"address"`

	// The value of the `nonce` found in the last transaction this account made
	// on the blockchain.
	Nonce uint64 `json:"nonce"`

	// The balance of the account in coins.
	Balance uint64 `json:"balance"`
}

// AccountRepository interface defines the methods for interacting with the account repository.
// This interface provides a way to manage accounts, including upserting, getting, listing, and deleting.
type AccountRepository interface {
	// Upsert inserts or updates an account in the repository.
	Upsert(acc *Account) error

	// GetByAddress retrieves an account by its ID.
	GetByAddress(addr *common.Address) (*Account, error)

	// ListAll retrieves all accounts in the repository.
	ListAll() ([]*Account, error)

	// DeleteByID deletes an account by its ID.
	DeleteByAddress(addr *common.Address) error

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
