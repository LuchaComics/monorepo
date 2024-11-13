package domain

import (
	"context"
	"crypto/rand"
	"fmt"

	"github.com/fxamacker/cbor/v2"
	"github.com/libp2p/go-libp2p/core/crypto"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type LibP2PNetworkPeerUniqueIdentifier struct {
	ID              primitive.ObjectID `bson:"_id" json:"id"`
	Label           string             `bson:"label" json:"label"`
	PrivateKeyBytes []byte             `bson:"private_key_bytes" json:"private_key_bytes"`
	PublicKeyBytes  []byte             `bson:"public_key_bytes" json:"public_key_bytes"`
}

type LibP2PNetworkPeerUniqueIdentifierRepository interface {
	GetOrCreate(ctx context.Context, label string) (*LibP2PNetworkPeerUniqueIdentifier, error)
	GetByLabel(ctx context.Context, label string) (*LibP2PNetworkPeerUniqueIdentifier, error)
	Upsert(ctx context.Context, key *LibP2PNetworkPeerUniqueIdentifier) error
}

func NewLibP2PNetworkPeerUniqueIdentifier(label string) (*LibP2PNetworkPeerUniqueIdentifier, error) {
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

	return &LibP2PNetworkPeerUniqueIdentifier{
		ID:              primitive.NewObjectID(),
		Label:           label,
		PrivateKeyBytes: privBytes,
		PublicKeyBytes:  pubBytes,
	}, nil
}

func (ik LibP2PNetworkPeerUniqueIdentifier) GetPrivateKey() (crypto.PrivKey, error) {
	// Unmarshal the private key from protobuf format
	priv, err := crypto.UnmarshalPrivateKey(ik.PrivateKeyBytes)
	if err != nil {
		return nil, err
	}
	return priv, nil
}

func (ik LibP2PNetworkPeerUniqueIdentifier) GetPublicKey() (crypto.PubKey, error) {
	pub, err := crypto.UnmarshalPublicKey(ik.PublicKeyBytes)
	if err != nil {
		return nil, err
	}
	return pub, nil
}

func (b *LibP2PNetworkPeerUniqueIdentifier) Serialize() ([]byte, error) {
	dataBytes, err := cbor.Marshal(b)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize identity key: %v", err)
	}
	return dataBytes, nil
}

func NewLibP2PNetworkPeerUniqueIdentifierFromDeserialize(data []byte) (*LibP2PNetworkPeerUniqueIdentifier, error) {
	// Variable we will use to return.
	identityKey := &LibP2PNetworkPeerUniqueIdentifier{}

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
