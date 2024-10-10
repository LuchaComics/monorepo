package domain

import (
	"crypto/ecdsa"
	"fmt"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/pkg/blockchain/signature"
	"github.com/ethereum/go-ethereum/crypto"
)

// Validator represents a trusted validator in the network.
type Validator struct {
	ID             string
	PublicKeyBytes []byte
}

func (validator *Validator) SignBlockHeader(privateKey *ecdsa.PrivateKey, blockHeader *BlockHeader) (string, error) {
	v, r, s, err := signature.Sign(blockHeader, privateKey)
	if err != nil {
		return "", err
	}
	blockHeaderSignatureString := signature.SignatureString(v, r, s)
	return blockHeaderSignatureString, nil
}

func (validator *Validator) ValidateBlockHeader(blockHeaderSignature string) bool {
	// Defensive Code.
	if blockHeaderSignature == "" {
		return false
	}

	v, r, s, err := signature.ToVRSFromHexSignature(blockHeaderSignature)
	if err != nil {
		return false
	}

	if err := signature.VerifySignature(v, r, s); err != nil {
		return false
	}
	return true
}

func (validator *Validator) GetPublicKeyECDSA() (*ecdsa.PublicKey, error) {
	publicKeyECDSA, err := crypto.UnmarshalPubkey(validator.PublicKeyBytes)
	if err != nil {
		return nil, fmt.Errorf("failed unmarshalling validator public key: %s", err)
	}
	return publicKeyECDSA, nil
}
