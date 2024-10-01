package domain

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"math/big"
)

// SignedTransaction is a signed version of the transaction. This is how
// clients like a wallet provide transactions for inclusion into the blockchain.
type SignedTransaction struct {
	Transaction
	V *big.Int `json:"v"` // Ethereum: Recovery identifier, either 29 or 30 with ardanID.
	R *big.Int `json:"r"` // Ethereum: First coordinate of the ECDSA signature.
	S *big.Int `json:"s"` // Ethereum: Second coordinate of the ECDSA signature.
}

// Validate checks if the transaction is valid. It verifies the signature,
// makes sure the account addresses are correct, and checks if the 'from'
// and 'to' accounts are not the same.
func (tx SignedTransaction) Validate(chainID uint16) error {
	// Check if the transaction's chain ID matches the expected one.
	if tx.ChainID != chainID {
		return fmt.Errorf("invalid chain id, got[%d] exp[%d]", tx.ChainID, chainID)
	}

	// Ensure the 'from' and 'to' accounts are not the same.
	if tx.From == tx.To {
		return fmt.Errorf("transaction invalid, sending money to yourself, from %s, to %s", tx.From, tx.To)
	}

	// Validate the signature parts (R, S, and V).
	if err := VerifySignature(tx.V, tx.R, tx.S); err != nil {
		return err
	}

	// Verify that the 'from' address matches the one from the signature.
	address, err := tx.FromAddress()
	if err != nil {
		return err
	}

	if address != string(tx.From.Hex()) {
		return errors.New("signature address doesn't match from address")
	}

	return nil
}

type SignedTransactionRepository interface {
	Upsert(bd *SignedTransaction) error
	ListAll() ([]*SignedTransaction, error)
	DeleteAll() error
}

func (stx *SignedTransaction) Serialize() ([]byte, error) {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)
	err := encoder.Encode(stx)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize block data: %v", err)
	}
	return result.Bytes(), nil
}

func NewSignedTransactionFromDeserialize(data []byte) (*SignedTransaction, error) {
	// Variable we will use to return.
	stx := &SignedTransaction{}

	// Defensive code: If programmer entered empty bytes then we will
	// return nil deserialization result.
	if data == nil {
		return nil, nil
	}

	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&stx)
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize block data: %v", err)
	}
	return stx, nil
}
