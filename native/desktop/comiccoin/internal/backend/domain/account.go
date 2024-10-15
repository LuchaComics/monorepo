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

	// The Nonce field in the Account struct is not directly related to the
	// Nonce field in the BlockHeader struct. Instead, it's used to prevent
	// replay attacks on transactions.
	//
	// In Ethereum and other blockchain systems, a nonce is a counter that is
	// used to prevent a transaction from being replayed on the network. When a
	// user wants to send a transaction, they need to specify the nonce value
	// for their account. The nonce value is incremented each time a transaction
	// is sent from the account.
	//
	// In this context, the Nonce field in the Account struct represents the
	// nonce value of the last transaction sent from the account. This value is
	// used to prevent replay attacks by ensuring that each transaction has a
	// unique nonce value.
	//
	// Here's how it works:
	// 1. When a user wants to send a transaction, they retrieve the current nonce value for their account from the blockchain.
	// 2. They increment the nonce value by 1 and include it in the transaction.
	// 3. The transaction is sent to the network and verified by the nodes.
	// 4. If the transaction is valid, the nonce value is incremented again and stored in the account's state.
	//
	// By including the Nonce field in the Account struct, the blockchain can
	// keep track of the nonce value for each account and prevent replay attacks.
	//
	// It's worth noting that the Nonce field in the BlockHeader struct has a
	// different purpose. It's used to solve the proof-of-work puzzle required
	// to mine a new block.
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
