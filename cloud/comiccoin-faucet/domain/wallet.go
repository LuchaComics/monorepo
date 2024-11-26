package domain

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/fxamacker/cbor/v2"
)

type Wallet struct {
	// The (Optional) description for this wallet.
	Label string `bson:"label" json:"label"`

	// The public address of the wallet.
	Address *common.Address `bson:"address" json:"address"`

	// FileContent contains
	KeystoreBytes []byte `bson:"keystore_bytes" json:"keystore_bytes"`
}

type WalletRepository interface {
	// Upsert inserts or updates an wallet in the repository.
	Upsert(ctx context.Context, acc *Wallet) error

	// GetByID retrieves an wallet by its Address.
	GetByAddress(ctx context.Context, address *common.Address) (*Wallet, error)

	// ListAll retrieves all wallets in the repository.
	ListAll(ctx context.Context) ([]*Wallet, error)

	// ListAllAddresses retrieves all wallet addresses in the repository.
	ListAllAddresses(ctx context.Context) ([]*common.Address, error)

	// DeleteByID deletes an wallet by its Address.
	DeleteByAddress(ctx context.Context, address *common.Address) error

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
