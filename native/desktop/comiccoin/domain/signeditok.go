package domain

import (
	"crypto/ecdsa"
	"fmt"
	"math/big"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/common/blockchain/signature"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/fxamacker/cbor/v2"
)

type IssuedToken struct {
	ID          uint64 `json:"id"`
	MetadataURI string `json:"metadata_uri"` // ComicCoin: URI pointing to Token metadata file (if this transaciton is an Token).
}

type SignedIssuedToken struct {
	IssuedToken
	V *big.Int `json:"v"` // Ethereum: Recovery identifier, either 29 or 30 with ardanID.
	R *big.Int `json:"r"` // Ethereum: First coordinate of the ECDSA signature.
	S *big.Int `json:"s"` // Ethereum: Second coordinate of the ECDSA signature.
}

type SignedIssuedTokenRepository interface {
	Upsert(bd *SignedIssuedToken) error
	ListAll() ([]*SignedIssuedToken, error)
	GetByID(id uint64) (*SignedIssuedToken, error)
	DeleteAll() error
	OpenTransaction() error
	CommitTransaction() error
	DiscardTransaction()
}

// Sign function signs the issued token using the user's private key
// and returns a signed version of that issued token.
func (itok *IssuedToken) Sign(privateKey *ecdsa.PrivateKey) (*SignedIssuedToken, error) {
	// Break the signature into the 3 parts: R, S, and V.
	v, r, s, err := signature.Sign(itok, privateKey)
	if err != nil {
		return &SignedIssuedToken{}, err
	}

	// Create the signed transaction, including the original transaction and the signature parts.
	signedTok := &SignedIssuedToken{
		IssuedToken: *itok,
		V:           v,
		R:           r,
		S:           s,
	}

	return signedTok, nil
}

func (sitok *SignedIssuedToken) Verify() error {
	return signature.VerifySignature(sitok.V, sitok.R, sitok.S)
}

func (sitok *SignedIssuedToken) PublicKey() (*ecdsa.PublicKey, error) {
	// Prepare the data for public key extraction.
	hashedData, err := signature.HashWithComicCoinStamp(sitok)
	if err != nil {
		return nil, err
	}
	return signature.GetPublicKeyFromSignature(hashedData, sitok.V, sitok.R, sitok.S)
}

func (sitok *SignedIssuedToken) PublicKeyBytes() ([]byte, error) {
	publicKeyECDSA, err := sitok.PublicKey()
	if err != nil {
		return nil, err
	}
	publicKeyBytes := crypto.FromECDSAPub(publicKeyECDSA)
	return publicKeyBytes, nil
}

// Serialize serializes the token into a byte slice.
// This method uses the cbor library to marshal the token into a byte slice.
func (b *SignedIssuedToken) Serialize() ([]byte, error) {
	// Marshal the token into a byte slice using the cbor library.
	dataBytes, err := cbor.Marshal(b)
	if err != nil {
		// Return an error if the marshaling fails.
		return nil, fmt.Errorf("failed to serialize signed issue token: %v", err)
	}
	return dataBytes, nil
}

// NewSignedIssuedTokenFromDeserialize deserializes an token from a byte slice.
// This method uses the cbor library to unmarshal the byte slice into an token.
func NewSignedIssuedTokenFromDeserialize(data []byte) (*SignedIssuedToken, error) {
	// Create a new token variable to return.
	token := &SignedIssuedToken{}

	// Defensive code: If the input data is empty, return a nil deserialization result.
	if data == nil {
		return nil, nil
	}

	// Unmarshal the byte slice into the token variable using the cbor library.
	if err := cbor.Unmarshal(data, &token); err != nil {
		// Return an error if the unmarshaling fails.
		return nil, fmt.Errorf("failed to deserialize signed issued token: %v", err)
	}
	return token, nil
}
