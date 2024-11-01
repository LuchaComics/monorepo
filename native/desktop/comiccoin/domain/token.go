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
	Nonce       uint64          `json:"nonce"`        // ComicCoin: Newly minted tokens always start at zero and for every transaction action afterwords (transfer, burn, etc) this value is increment by 1.
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

	// ListByOwner retrieves all the tokens in the repository that belongs
	// to the owner address.
	ListByOwner(owner *common.Address) ([]*Token, error)

	// DeleteByID deletes an token by its ID.
	DeleteByID(id uint64) error

	// HashState returns a hash based on the contents of the tokens and
	// their metadata. This is added to each block and checked by peers.
	HashState() (string, error)

	OpenTransaction() error
	CommitTransaction() error
	DiscardTransaction()
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

func ToTokenIDsArray(toks []*Token) []uint64 {
	tokIDs := make([]uint64, len(toks))
	for _, tok := range toks {
		tokIDs = append(tokIDs, tok.ID)
	}
	return tokIDs
}
