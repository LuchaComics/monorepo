package domain

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/fxamacker/cbor/v2"
)

type Token struct {
	// The unique identifier for this blockchain that we are managing the state for.
	ChainID uint16 `bson:"chain_id" json:"chain_id"`

	IDBytes     []byte          `bson:"id_bytes" json:"id_bytes"`
	Owner       *common.Address `bson:"owner" json:"owner"`
	MetadataURI string          `bson:"metadata_uri" json:"metadata_uri"` // ComicCoin: URI pointing to Token metadata file (if this transaciton is an Token).
	NonceBytes  []byte          `bson:"nonce_bytes" json:"nonce_bytes"`   // ComicCoin: Newly minted tokens always start at zero and for every transaction action afterwords (transfer, burn, etc) this value is increment by 1.
}

// TokenRepository interface defines the methods for interacting with the token repository.
// This interface provides a way to manage tokens, including upserting, getting, listing, and deleting.
type TokenRepository interface {
	// Upsert inserts or updates an token in the repository.
	Upsert(ctx context.Context, tok *Token) error

	// GetByAddress retrieves an token by its ID.
	GetByID(ctx context.Context, id *big.Int) (*Token, error)

	// ListByChainID retrieves all accounts in the repository for the particular chain ID.
	ListByChainID(ctx context.Context, chainID uint16) ([]*Token, error)

	// ListByOwner retrieves all the tokens in the repository that belongs
	// to the owner address.
	ListByOwner(ctx context.Context, owner *common.Address) ([]*Token, error)

	// CountByOwner counts all the tokens owned by the owner.
	CountByOwner(ctx context.Context, owner *common.Address) (int64, error)

	// DeleteByID deletes an token by its ID.
	DeleteByID(ctx context.Context, id *big.Int) error

	// HashStateByChainID returns a hash based on the contents of the tokens and
	// their metadata. This is added to each block and checked by peers.
	HashStateByChainID(ctx context.Context, chainID uint16) (string, error)

	OpenTransaction() error
	CommitTransaction() error
	DiscardTransaction()
}

func (tok *Token) GetID() *big.Int {
	return new(big.Int).SetBytes(tok.IDBytes)
}

func (tok *Token) IsLatestTokenIDZero() bool {
	// Special thanks to "Is there another way of testing if a big.Int is 0?" via https://stackoverflow.com/a/64257532
	return len(tok.GetID().Bits()) == 0
}

func (tok *Token) SeLatestTokenID(n *big.Int) {
	tok.IDBytes = n.Bytes()
}

func (tok *Token) GetNonce() *big.Int {
	return new(big.Int).SetBytes(tok.NonceBytes)
}

func (tok *Token) IsNonceZero() bool {
	// Special thanks to "Is there another way of testing if a big.Int is 0?" via https://stackoverflow.com/a/64257532
	return len(tok.GetNonce().Bits()) == 0
}

func (tok *Token) SetNonce(n *big.Int) {
	tok.NonceBytes = n.Bytes()
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

func ToTokenIDsArray(toks []*Token) []*big.Int {
	tokIDs := make([]*big.Int, len(toks))
	for _, tok := range toks {
		tokIDs = append(tokIDs, tok.GetID())
	}
	return tokIDs
}
