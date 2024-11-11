package domain

import (
	"crypto/ecdsa"
	"fmt"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/fxamacker/cbor/v2"
)

type IssuedToken struct {
	ID          uint64 `json:"id"`
	MetadataURI string `json:"metadata_uri"` // ComicCoin: URI pointing to Token metadata file (if this transaciton is an Token).
}

type SignedIssuedToken struct {
	IssuedToken

	// The signature of this block's "IssuedToken" field which was applied by the
	// proof-of-authority validator.
	IssuedTokenSignatureBytes []byte `json:"issued_token_signature_bytes"`

	// The proof-of-authority validator whom executed the validation of
	// this NFT in our blockchain. Must match genesis block validator.
	Validator *Validator `json:"validator"`
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

// SignUsingProofOfAuthorityValidator function signs the issued token using the
// proof of authorities private key and returns a signed version of that issued token.
func (itok *IssuedToken) SignUsingProofOfAuthorityValidator(poaValidator *Validator, poaPrivateKey *ecdsa.PrivateKey) (*SignedIssuedToken, error) {
	issuedTokenSignatureBytes, err := poaValidator.Sign(poaPrivateKey, itok)
	if err != nil {
		return nil, err
	}

	// Create the signed transaction, including the original transaction and the signature parts.
	signedIssuedTok := &SignedIssuedToken{
		IssuedToken:               *itok,
		IssuedTokenSignatureBytes: issuedTokenSignatureBytes,
		Validator:                 poaValidator,
	}

	return signedIssuedTok, nil
}

func (sitok *SignedIssuedToken) Verify() bool {
	return sitok.Validator.Verify(sitok.IssuedTokenSignatureBytes, sitok.IssuedToken)
}

func (sitok *SignedIssuedToken) PublicKey() (*ecdsa.PublicKey, error) {
	return sitok.Validator.GetPublicKeyECDSA()
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
