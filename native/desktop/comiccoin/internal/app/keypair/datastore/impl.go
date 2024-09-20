package datastore

import (
	"context"
	"crypto/rand"
	"log/slog"

	"github.com/libp2p/go-libp2p/core/crypto"
)

func (impl *keypairStorerImpl) GetByName(ctx context.Context, name string) (crypto.PrivKey, crypto.PubKey, error) {
	// Get the private key bytes from the database
	privBytes, err := impl.dbClient.Getf("KEYNAME_%v_PRIVATE", name)
	if err != nil {
		impl.logger.Error("failed getting private key", slog.Any("error", err))
		return nil, nil, err
	}

	// Unmarshal the private key from protobuf format
	priv, err := crypto.UnmarshalPrivateKey(privBytes)
	if err != nil {
		impl.logger.Error("failed unmarshalling private key", slog.Any("error", err))
		return nil, nil, err
	}

	// Get the public key bytes from the database
	pubBytes, err := impl.dbClient.Getf("KEYNAME_%v_PUBLIC", name)
	if err != nil {
		impl.logger.Error("failed getting public key", slog.Any("error", err))
		return nil, nil, err
	}

	// Unmarshal the public key from protobuf format
	pub, err := crypto.UnmarshalPublicKey(pubBytes)
	if err != nil {
		impl.logger.Error("failed unmarshalling public key", slog.Any("error", err))
		return nil, nil, err
	}

	return priv, pub, nil
}

func (impl *keypairStorerImpl) GenerateNewKeyPairAndSetByName(ctx context.Context, name string) error {
	r := rand.Reader

	// Generate a key pair for this host
	priv, pub, err := crypto.GenerateKeyPairWithReader(crypto.RSA, 2048, r)
	if err != nil {
		return err
	}

	// Marshal the private key in protobuf format
	privBytes, err := crypto.MarshalPrivateKey(priv)
	if err != nil {
		return err
	}

	// Store the private key in the database
	if err := impl.dbClient.Setf(privBytes, "KEYNAME_%v_PRIVATE", name); err != nil {
		return err
	}

	// Marshal the public key in protobuf format
	pubBytes, err := crypto.MarshalPublicKey(pub)
	if err != nil {
		return err
	}

	// Store the public key in the database
	if err := impl.dbClient.Setf(pubBytes, "KEYNAME_%v_PUBLIC", name); err != nil {
		return err
	}

	return nil
}
