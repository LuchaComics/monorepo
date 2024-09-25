package datastore

import (
	"crypto/ecdsa"
	"encoding/json"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

func (tx *PendingTransaction) encodeRLP() []byte {
	encoded, _ := json.Marshal(struct {
		From   common.Address `json:"from"`
		To     common.Address `json:"to"`
		Amount *big.Int       `json:"value"`
		Data   []byte         `json:"data"`
	}{
		From:   tx.From,
		To:     tx.To,
		Amount: tx.Amount,
		Data:   tx.Data,
	})
	return encoded
}

func (tx *PendingTransaction) Hash() common.Hash {
	return crypto.Keccak256Hash(tx.encodeRLP())
}

func (tx *PendingTransaction) Sign(privateKey *ecdsa.PrivateKey) error {
	hash := tx.Hash()
	signature, err := crypto.Sign(hash.Bytes(), privateKey)
	if err != nil {
		return err
	}
	tx.Signature = signature
	return nil
}

func (tx *PendingTransaction) Verify() bool {
	hash := tx.Hash()
	sigPublicKey, err := crypto.Ecrecover(hash.Bytes(), tx.Signature)
	if err != nil {
		return false
	}
	sigPublicKeyECDSA, err := crypto.UnmarshalPubkey(sigPublicKey)
	if err != nil {
		return false
	}
	recoveredAddr := crypto.PubkeyToAddress(*sigPublicKeyECDSA)
	return recoveredAddr == tx.From
}
