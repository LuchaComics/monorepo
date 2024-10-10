package domain

import (
	"crypto/ecdsa"
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/crypto"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/pkg/blockchain/signature"
)

// Validator represents a trusted validator in the network.
type Validator struct {
	ID             string
	PublicKeyBytes []byte
}

func (validator *Validator) Sign(privateKey *ecdsa.PrivateKey, data any) ([]byte, error) {
	// Prepare the data for signing.
	hash, err := signature.HashWithComicCoinStamp(data)
	if err != nil {
		return nil, err
	}

	// Sign the hash with the private key to produce a signature.
	hashSignature, err := crypto.Sign(hash, privateKey)
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
		return false
	}

	// Prepare the data for signing.
	hash, err := signature.HashWithComicCoinStamp(data)
	if err != nil {
		log.Printf("VALIDATOR: VERIFY FAILED: HashWithComicCoinStamp err %v\n", err)
		return false
	}

	// Get our validators public key.
	validatorPubKey, err := validator.GetPublicKeyECDSA()
	if err != nil {
		log.Printf("VALIDATOR: VERIFY FAILED: GetPublicKeyECDSA err %v\n", err)
		return false
	}

	// Get the public key from the signature.
	sigPubKey, err := crypto.SigToPub(hash, sig)
	if err != nil {
		log.Printf("VALIDATOR: VERIFY FAILED: crypto.SigToPub err %v\n", err)
		return false
	}

	// Verify signature public key and validator public key match.
	if validatorPubKey != sigPubKey { //TODO: CONFIRM THIS WORKS
		log.Printf("VALIDATOR: VERIFY FAILED: %v\n", "validatorPubKey != sigPubKey")
		return false
	}

	// Perform our verification.
	sigPubKeyBytes, err := crypto.Ecrecover(hash, sig)
	if err != nil {
		log.Printf("VALIDATOR: VERIFY FAILED: crypto.Ecrecover err %v\n", err)
		return false
	}
	return crypto.VerifySignature(sigPubKeyBytes, hash, sig)
}

func (validator *Validator) GetPublicKeyECDSA() (*ecdsa.PublicKey, error) {
	publicKeyECDSA, err := crypto.UnmarshalPubkey(validator.PublicKeyBytes)
	if err != nil {
		return nil, fmt.Errorf("failed unmarshalling validator public key: %s", err)
	}
	return publicKeyECDSA, nil
}
