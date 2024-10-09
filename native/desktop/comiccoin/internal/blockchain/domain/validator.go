package domain

import (
	"crypto"
	"crypto/ecdsa"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/pkg/blockchain/signature"
)

// Validator represents a trusted validator in the network.
type Validator struct {
	ID        string
	PublicKey *crypto.PublicKey
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
	v, r, s, err := signature.ToVRSFromHexSignature(blockHeaderSignature)
	if err != nil {
		return false
	}

	if err := signature.VerifySignature(v, r, s); err != nil {
		return false
	}
	return true
}
