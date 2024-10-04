package domain

import (
	"crypto/rand"
	"fmt"

	"github.com/fxamacker/cbor/v2"
	"github.com/libp2p/go-libp2p/core/crypto"
)

type IdentityKey struct {
	ID              string `json:"id"`
	PrivateKeyBytes []byte `json:"private_key_bytes"`
	PublicKeyBytes  []byte `json:"public_key_bytes"`
}

type IdentityKeyRepository interface {
	GetByID(id string) (*IdentityKey, error)
	Upsert(key *IdentityKey) error
}

func NewIdentityKey(id string) (*IdentityKey, error) {
	r := rand.Reader

	priv, pub, err := crypto.GenerateKeyPairWithReader(crypto.RSA, 2048, r)
	if err != nil {
		return nil, err
	}

	// Marshal the private key in protobuf format
	privBytes, err := crypto.MarshalPrivateKey(priv)
	if err != nil {
		return nil, err
	}

	// Marshal the public key in protobuf format
	pubBytes, err := crypto.MarshalPublicKey(pub)
	if err != nil {
		return nil, err
	}

	return &IdentityKey{
		ID:              id,
		PrivateKeyBytes: privBytes,
		PublicKeyBytes:  pubBytes,
	}, nil
}

func (ik IdentityKey) GetPrivateKey() (crypto.PrivKey, error) {
	// Unmarshal the private key from protobuf format
	priv, err := crypto.UnmarshalPrivateKey(ik.PrivateKeyBytes)
	if err != nil {
		return nil, err
	}
	return priv, nil
}

func (ik IdentityKey) GetPublicKey() (crypto.PubKey, error) {
	pub, err := crypto.UnmarshalPublicKey(ik.PublicKeyBytes)
	if err != nil {
		return nil, err
	}
	return pub, nil
}

func (b *IdentityKey) Serialize() ([]byte, error) {
	dataBytes, err := cbor.Marshal(b)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize identity key: %v", err)
	}
	return dataBytes, nil
}

func NewIdentityKeyFromDeserialize(data []byte) (*IdentityKey, error) {
	// Variable we will use to return.
	identityKey := &IdentityKey{}

	// Defensive code: If programmer entered empty bytes then we will
	// return nil deserialization result.
	if data == nil {
		return nil, nil
	}

	if err := cbor.Unmarshal(data, &identityKey); err != nil {
		return nil, fmt.Errorf("failed to deserialize identity key: %v", err)
	}
	return identityKey, nil
}
