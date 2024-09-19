package blockchain

import (
	"crypto/ecdsa"
	"encoding/hex"
	"encoding/json"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

type Transaction struct {
	From      common.Address `json:"from"`
	To        common.Address `json:"to"`
	Value     *big.Int       `json:"value"`
	Data      []byte         `json:"data"`
	Nonce     uint64         `json:"nonce"`
	Signature []byte         `json:"signature"`
}

type SignedTransaction struct {
	Transaction
	Sig []byte `json:"signature"`
}

func NewSignedTransaction(tx Transaction, sig []byte) SignedTransaction {
	return SignedTransaction{tx, sig}
}

func NewTransaction(from, to common.Address, value *big.Int, data []byte, nonce uint64) *Transaction {
	return &Transaction{
		From:  from,
		To:    to,
		Value: value,
		Data:  data,
		Nonce: nonce,
	}
}

func (tx *Transaction) Hash() common.Hash {
	return crypto.Keccak256Hash(tx.encodeRLP())
}

func (tx Transaction) Encode() ([]byte, error) {
	return json.Marshal(tx)
}

func (tx *Transaction) encodeRLP() []byte {
	encoded, _ := json.Marshal(struct {
		From  common.Address `json:"from"`
		To    common.Address `json:"to"`
		Value *big.Int       `json:"value"`
		Data  []byte         `json:"data"`
		Nonce uint64         `json:"nonce"`
	}{
		From:  tx.From,
		To:    tx.To,
		Value: tx.Value,
		Data:  tx.Data,
		Nonce: tx.Nonce,
	})
	return encoded
}

func (tx *Transaction) Sign(privateKey *ecdsa.PrivateKey) error {
	hash := tx.Hash()
	signature, err := crypto.Sign(hash.Bytes(), privateKey)
	if err != nil {
		return err
	}
	tx.Signature = signature
	return nil
}

func (tx *Transaction) Verify() bool {
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

func (tx *Transaction) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		From      string `json:"from"`
		To        string `json:"to"`
		Value     string `json:"value"`
		Data      string `json:"data"`
		Nonce     uint64 `json:"nonce"`
		Signature string `json:"signature"`
	}{
		From:      tx.From.Hex(),
		To:        tx.To.Hex(),
		Value:     tx.Value.String(),
		Data:      hex.EncodeToString(tx.Data),
		Nonce:     tx.Nonce,
		Signature: hex.EncodeToString(tx.Signature),
	})
}

func (tx *Transaction) UnmarshalJSON(data []byte) error {
	var temp struct {
		From      string `json:"from"`
		To        string `json:"to"`
		Value     string `json:"value"`
		Data      string `json:"data"`
		Nonce     uint64 `json:"nonce"`
		Signature string `json:"signature"`
	}
	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}
	tx.From = common.HexToAddress(temp.From)
	tx.To = common.HexToAddress(temp.To)
	tx.Value, _ = new(big.Int).SetString(temp.Value, 10)
	tx.Data, _ = hex.DecodeString(temp.Data)
	tx.Nonce = temp.Nonce
	tx.Signature, _ = hex.DecodeString(temp.Signature)
	return nil
}
