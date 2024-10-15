package domain

import (
	"errors"
	"fmt"
	"log"
	"math/big"

	"github.com/fxamacker/cbor/v2"
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
// and 'to' accounts are not the same (unless you are the proof of authority!)
func (stx SignedTransaction) Validate(chainID uint16, isPoA bool) error {
	// Check if the transaction's chain ID matches the expected one.
	if stx.ChainID != chainID {
		return fmt.Errorf("invalid chain id, got[%d] exp[%d]", stx.ChainID, chainID)
	}

	// Ensure the 'from' and 'to' accounts are not the same.
	if stx.From == stx.To {
		// ... unless you are the proof of authority.
		if !isPoA {
			return fmt.Errorf("transaction invalid, sending money to yourself, from %s, to %s", stx.From, stx.To)
		}
	}

	// Validate the signature parts (R, S, and V).
	if err := VerifySignature(stx.V, stx.R, stx.S); err != nil {
		return err
	}

	// Verify that the 'from' address matches the one from the signature.
	address, err := stx.FromAddress()
	if err != nil {
		return err
	}

	if address != string(stx.From.Hex()) {
		log.Printf("SignedTransaction: Validate: signature address %v doesn't match from address %v\n", address, stx.From.Hex())
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
	dataBytes, err := cbor.Marshal(stx)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize signed transaction: %v", err)
	}
	return dataBytes, nil
}

func NewSignedTransactionFromDeserialize(data []byte) (*SignedTransaction, error) {
	// Variable we will use to return.
	stx := &SignedTransaction{}

	// Defensive code: If programmer entered empty bytes then we will
	// return nil deserialization result.
	if data == nil {
		return nil, nil
	}

	if err := cbor.Unmarshal(data, &stx); err != nil {
		return nil, fmt.Errorf("failed to deserialize signed transaction: %v", err)
	}
	return stx, nil
}
