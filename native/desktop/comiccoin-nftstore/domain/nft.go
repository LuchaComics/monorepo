package domain

import (
	"fmt"

	"github.com/fxamacker/cbor/v2"
)

type NFT struct {
	TokenID     uint64      `json:"token_id"`
	MetadataURI string      `json:"metadata_uri"`
	Metadata    NFTMetadata `json:"metadata"`
}

type NFTMetadata struct {
	Image           string                  `bson:"image" json:"image"`
	ExternalURL     string                  `bson:"external_url" json:"external_url"`
	Description     string                  `bson:"description" json:"description"`
	Name            string                  `bson:"name" json:"name"`
	Attributes      []*NFTMetadataAttribute `bson:"attributes" json:"attributes"`
	BackgroundColor string                  `bson:"background_color" json:"background_color"`
	AnimationURL    string                  `bson:"animation_url" json:"animation_url"`
	YoutubeURL      string                  `bson:"youtube_url" json:"youtube_url"`
}

type NFTMetadataAttribute struct {
	DisplayType string `bson:"display_type" json:"display_type"`
	TraitType   string `bson:"trait_type" json:"trait_type"`
	Value       string `bson:"value" json:"value"`
}

type NFTRepository interface {
	// Upsert inserts or updates an wallet in the repository.
	Upsert(acc *NFT) error

	// GetByID retrieves an wallet by its Address.
	GetByTokenID(tokenID uint64) (*NFT, error)

	// ListAll retrieves all wallets in the repository.
	ListAll() ([]*NFT, error)

	// DeleteByID deletes an wallet by its Address.
	DeleteByTokenID(tokenID uint64) error

	OpenTransaction() error

	CommitTransaction() error

	DiscardTransaction()
}

// Serialize serializes the wallet into a byte slice.
// This method uses the cbor library to marshal the wallet into a byte slice.
func (b *NFT) Serialize() ([]byte, error) {
	// Marshal the wallet into a byte slice using the cbor library.
	dataBytes, err := cbor.Marshal(b)
	if err != nil {
		// Return an error if the marshaling fails.
		return nil, fmt.Errorf("failed to serialize wallet: %v", err)
	}
	return dataBytes, nil
}

// NewNFTFromDeserialize deserializes an wallet from a byte slice.
// This method uses the cbor library to unmarshal the byte slice into an wallet.
func NewNFTFromDeserialize(data []byte) (*NFT, error) {
	// Create a new wallet variable to return.
	wallet := &NFT{}

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