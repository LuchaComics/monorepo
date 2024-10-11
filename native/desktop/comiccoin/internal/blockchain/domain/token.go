package domain

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/fxamacker/cbor/v2"
)

type Token struct {
	ID          uint64          `json:"id"`
	Owner       *common.Address `json:"owner"`
	MetadataURI string          `json:"metadata_uri"` // ComicCoin: URI pointing to Token metadata file (if this transaciton is an Token).
}

// TokenRepository interface defines the methods for interacting with the token repository.
// This interface provides a way to manage tokens, including upserting, getting, listing, and deleting.
type TokenRepository interface {
	// Upsert inserts or updates an token in the repository.
	Upsert(acc *Token) error

	// GetByAddress retrieves an token by its ID.
	GetByID(id uint64) (*Token, error)

	// ListAll retrieves all tokens in the repository.
	ListAll() ([]*Token, error)

	// DeleteByID deletes an token by its ID.
	DeleteByID(id uint64) error

	// HashState returns a hash based on the contents of the tokens and
	// their metadata. This is added to each block and checked by peers.
	HashState() (string, error)
}

// Serialize serializes the token into a byte slice.
// This method uses the cbor library to marshal the token into a byte slice.
func (b *Token) Serialize() ([]byte, error) {
	// Marshal the token into a byte slice using the cbor library.
	dataBytes, err := cbor.Marshal(b)
	if err != nil {
		// Return an error if the marshaling fails.
		return nil, fmt.Errorf("failed to serialize token: %v", err)
	}
	return dataBytes, nil
}

// NewTokenFromDeserialize deserializes an token from a byte slice.
// This method uses the cbor library to unmarshal the byte slice into an token.
func NewTokenFromDeserialize(data []byte) (*Token, error) {
	// Create a new token variable to return.
	token := &Token{}

	// Defensive code: If the input data is empty, return a nil deserialization result.
	if data == nil {
		return nil, nil
	}

	// Unmarshal the byte slice into the token variable using the cbor library.
	if err := cbor.Unmarshal(data, &token); err != nil {
		// Return an error if the unmarshaling fails.
		return nil, fmt.Errorf("failed to deserialize token: %v", err)
	}
	return token, nil
}
