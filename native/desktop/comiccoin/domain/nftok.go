package domain

import (
	"fmt"
	"math/big"

	"github.com/fxamacker/cbor/v2"
)

type NonFungibleToken struct {
	TokenID     *big.Int                  `json:"token_id"`
	MetadataURI string                    `json:"metadata_uri"`
	Metadata    *NonFungibleTokenMetadata `json:"metadata"`
	State       string                    `bson:"state" json:"state"`
}

type NonFungibleTokenMetadata struct {
	Image           string                               `bson:"image" json:"image"`
	ExternalURL     string                               `bson:"external_url" json:"external_url"`
	Description     string                               `bson:"description" json:"description"`
	Name            string                               `bson:"name" json:"name"`
	Attributes      []*NonFungibleTokenMetadataAttribute `bson:"attributes" json:"attributes"`
	BackgroundColor string                               `bson:"background_color" json:"background_color"`
	AnimationURL    string                               `bson:"animation_url" json:"animation_url"`
	YoutubeURL      string                               `bson:"youtube_url" json:"youtube_url"`
}

const (
	NonFungibleTokenStateNotReady = "not_ready"
	NonFungibleTokenStateReady    = "ready"
)

type NonFungibleTokenMetadataAttribute struct {
	DisplayType string `bson:"display_type" json:"display_type"`
	TraitType   string `bson:"trait_type" json:"trait_type"`
	Value       string `bson:"value" json:"value"`
}

type NonFungibleTokenRepository interface {
	// Upsert inserts or updates an nfts in the repository.
	Upsert(acc *NonFungibleToken) error

	// GetByID retrieves an nft by its token id.
	GetByTokenID(tokenID *big.Int) (*NonFungibleToken, error)

	// ListAll retrieves all nfts in the repository.
	ListAll() ([]*NonFungibleToken, error)

	// ListByTokenIDs retrieves nfts that have the token is in the parameter
	ListWithFilterByTokenIDs(tokIDs []*big.Int) ([]*NonFungibleToken, error)

	// DeleteByID deletes an nft by its token id.
	DeleteByTokenID(tokenID *big.Int) error

	OpenTransaction() error

	CommitTransaction() error

	DiscardTransaction()
}

// Serialize serializes the wallet into a byte slice.
// This method uses the cbor library to marshal the wallet into a byte slice.
func (b *NonFungibleToken) Serialize() ([]byte, error) {
	// Marshal the wallet into a byte slice using the cbor library.
	dataBytes, err := cbor.Marshal(b)
	if err != nil {
		// Return an error if the marshaling fails.
		return nil, fmt.Errorf("failed to serialize wallet: %v", err)
	}
	return dataBytes, nil
}

// NewNonFungibleTokenFromDeserialize deserializes an wallet from a byte slice.
// This method uses the cbor library to unmarshal the byte slice into an wallet.
func NewNonFungibleTokenFromDeserialize(data []byte) (*NonFungibleToken, error) {
	// Create a new wallet variable to return.
	wallet := &NonFungibleToken{}

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

func ToNonFungibleTokenIDsArray(nftoks []*NonFungibleToken) []*big.Int {
	nftokIDs := make([]*big.Int, len(nftoks))
	for _, nftok := range nftoks {
		nftokIDs = append(nftokIDs, nftok.TokenID)
	}
	return nftokIDs
}
