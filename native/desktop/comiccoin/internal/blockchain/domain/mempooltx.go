package domain

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/fxamacker/cbor/v2"
)

// MempoolTransaction is a mempool version of the transaction. This is how
// clients like a wallet provide transactions for inclusion into the blockchain.
type MempoolTransaction struct {
	Transaction
	V *big.Int `json:"v"` // Ethereum: Recovery identifier, either 29 or 30 with ardanID.
	R *big.Int `json:"r"` // Ethereum: First coordinate of the ECDSA signature.
	S *big.Int `json:"s"` // Ethereum: Second coordinate of the ECDSA signature.
}

// Validate checks if the transaction is valid. It verifies the signature,
// makes sure the account addresses are correct, and checks if the 'from'
// and 'to' accounts are not the same.
func (tx MempoolTransaction) Validate(chainID uint16) error {
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

type MempoolTransactionRepository interface {
	Upsert(bd *MempoolTransaction) error
	ListAll() ([]*MempoolTransaction, error)
	DeleteAll() error
}

func (stx *MempoolTransaction) Serialize() ([]byte, error) {
	dataBytes, err := cbor.Marshal(stx)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize mempool transaction: %v", err)
	}
	return dataBytes, nil
}

func NewMempoolTransactionFromDeserialize(data []byte) (*MempoolTransaction, error) {
	// Variable we will use to return.
	stx := &MempoolTransaction{}

	// Defensive code: If programmer entered empty bytes then we will
	// return nil deserialization result.
	if data == nil {
		return nil, nil
	}

	if err := cbor.Unmarshal(data, &stx); err != nil {
		return nil, fmt.Errorf("failed to deserialize mempool transaction: %v", err)
	}
	return stx, nil
}

// FromAddress extracts the account address from the signed transaction by
// recovering the public key from the signature.
func (tx MempoolTransaction) FromAddress() (string, error) {

	// Create the hash of the transaction to prepare it for extracting the public key.
	tran, err := tx.HashWithComicCoinStamp()
	if err != nil {
		return "", err
	}

	// Combine R, S, and V into the original signature format.
	sig := toSignatureBytes(tx.V, tx.R, tx.S)

	// Use the signature to get the public key of the account that signed it.
	publicKey, err := crypto.SigToPub(tran, sig)
	if err != nil {
		return "", err
	}

	// Convert the public key to an account address.
	return crypto.PubkeyToAddress(*publicKey).String(), nil
}

func (tx MempoolTransaction) ToSignedTransaction() *SignedTransaction {
	return &SignedTransaction{}
}
