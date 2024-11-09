package domain

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/crypto"
)

// Validator represents a trusted validator in the network.
type Validator struct {
	ID             string `bson:"id" json:"id"`
	PublicKeyBytes []byte `bson:"public_key_bytes" json:"public_key_bytes"`
}

func (validator *Validator) Sign(privateKey *ecdsa.PrivateKey, value any) ([]byte, error) {
	data, err := json.Marshal(value)
	if err != nil {
		return nil, err
	}

	// Prepare the data for signing.
	hash := sha256.Sum256(data)

	// Sign the data
	hashSignature, err := ecdsa.SignASN1(rand.Reader, privateKey, hash[:])
	if err != nil {
		return nil, err
	}

	// Return our result.
	return hashSignature, nil
}

func (validator *Validator) Verify(sig []byte, data any) bool {
	// Defensive Code.
	if sig == nil || data == nil {
		log.Printf("VALIDATOR: VERIFY FAILED: %v\n", "sig == nil || data == nil")
		log.Printf("VALIDATOR: VERIFY FAILED: sig %v\n", sig)
		log.Printf("VALIDATOR: VERIFY FAILED: data %v\n", data)
		return false
	}

	// Prepare the data for signing.
	dataBytes, err := json.Marshal(data)
	if err != nil {
		log.Printf("VALIDATOR: VERIFY FAILED: json.Marshal(value) err %v\n", err)
		return false
	}

	// Prepare the data for signing.
	hash := sha256.Sum256(dataBytes)

	// Get our validators public key.
	validatorPubKey, err := validator.GetPublicKeyECDSA()
	if err != nil {
		log.Printf("VALIDATOR: VERIFY FAILED: GetPublicKeyECDSA err %v\n", err)
		return false
	}

	// Verify the signature
	return ecdsa.VerifyASN1(validatorPubKey, hash[:], sig)
}

func (validator *Validator) GetPublicKeyECDSA() (*ecdsa.PublicKey, error) {
	if validator == nil {
		return nil, fmt.Errorf("validator error: %v", "d.n.e.")
	}

	publicKeyECDSA, err := crypto.UnmarshalPubkey(validator.PublicKeyBytes)
	if err != nil {
		return nil, fmt.Errorf("failed unmarshalling validator public key: %s", err)
	}
	return publicKeyECDSA, nil
}
