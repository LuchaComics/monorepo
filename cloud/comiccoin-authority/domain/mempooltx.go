package domain

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"math/big"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/common/blockchain/signature"
	"github.com/fxamacker/cbor/v2"
)

// MempoolTransaction represents a transaction that is stored in the mempool.
// It contains the transaction data, as well as the ECDSA signature and recovery identifier.
type MempoolTransaction struct {
	ID primitive.ObjectID `bson:"_id" json:"id"`

	// The transaction data, including sender, recipient, amount, and other metadata.
	Transaction

	// The recovery identifier, either 29 or 30, depending on the Ethereum network.
	V *big.Int `bson:"v" json:"v"`

	// The first coordinate of the ECDSA signature.
	R *big.Int `bson:"r" json:"r"`

	// The second coordinate of the ECDSA signature.
	S *big.Int `bson:"s" json:"s"`
}

// Validate checks if the transaction is valid.
// It verifies the signature, makes sure the account addresses are correct,
// and checks if the 'from' and 'to' accounts are not the same.
func (tx MempoolTransaction) Validate(chainID uint16, isPoA bool) error {
	if tx.V == nil || tx.R == nil || tx.S == nil {
		return errors.New("V, R, or S is nil")
	}

	if tx.V.Uint64() == 0 || tx.R.Uint64() == 0 || tx.S.Uint64() == 0 {
		return errors.New("V, R, or S is zero")
	}

	if tx.V.Uint64() < 29 || tx.V.Uint64() > 30 {
		return errors.New("V is out of range")
	}

	// Check if the transaction's chain ID matches the expected one.
	if tx.ChainID != chainID {
		return fmt.Errorf("invalid chain id, got[%d] exp[%d]", tx.ChainID, chainID)
	}

	// Ensure the 'from' and 'to' accounts are not the same.
	if tx.From == tx.To {
		// ... unless you are the proof of authority.
		if !isPoA {
			return fmt.Errorf("transaction invalid, sending money to yourself, from %s, to %s", tx.From, tx.To)
		}
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

type MempoolTransactionInsertionDetector interface {
	Monitor(ctx context.Context, onChangeDetectedFunc func(data *MempoolTransaction) error) error
	Close(ctx context.Context)
}

// MempoolTransactionRepository interface defines the methods for interacting with
// the mempool transaction repository.
// This interface provides a way to manage mempool transactions, including upserting, listing, and deleting.
type MempoolTransactionRepository interface {
	GetByID(ctx context.Context, id primitive.ObjectID) (*MempoolTransaction, error)

	// Upsert inserts or updates a mempool transaction in the repository.
	Upsert(ctx context.Context, mempoolTx *MempoolTransaction) error

	// ListAll retrieves all mempool transactions in the repository.
	ListByChainID(ctx context.Context, chainID uint16) ([]*MempoolTransaction, error)

	// DeleteAll deletes all mempool transactions in the repository.
	DeleteByChainID(ctx context.Context, chainID uint16) error

	GetInsertionChangeStreamChannel(ctx context.Context) (<-chan MempoolTransaction, chan struct{}, error)
}

// Serialize serializes the mempool transaction into a byte slice.
// This method uses the cbor library to marshal the transaction into a byte slice.
func (stx *MempoolTransaction) Serialize() ([]byte, error) {
	// Marshal the transaction into a byte slice using the cbor library.
	dataBytes, err := cbor.Marshal(stx)
	if err != nil {
		// Return an error if the marshaling fails.
		return nil, fmt.Errorf("failed to serialize mempool transaction: %v", err)
	}
	return dataBytes, nil
}

// NewMempoolTransactionFromDeserialize deserializes a mempool transaction from a byte slice.
// This method uses the cbor library to unmarshal the byte slice into a transaction.
func NewMempoolTransactionFromDeserialize(data []byte) (*MempoolTransaction, error) {
	// Create a new transaction variable to return.
	stx := &MempoolTransaction{}

	// Defensive code: If the input data is empty, return a nil deserialization result.
	if data == nil {
		return nil, nil
	}

	// Unmarshal the byte slice into the transaction variable using the cbor library.
	if err := cbor.Unmarshal(data, &stx); err != nil {
		// Return an error if the unmarshaling fails.
		return nil, fmt.Errorf("failed to deserialize mempool transaction: %v", err)
	}
	return stx, nil
}

// FromAddress extracts the account address from the signed transaction by
// recovering the public key from the signature.
func (tx MempoolTransaction) FromAddress() (string, error) {
	return signature.FromAddress(tx.Transaction, tx.V, tx.R, tx.S)
}

func (tx MempoolTransaction) FromPublicKey() (*ecdsa.PublicKey, error) {
	return signature.GetPublicKeyFromSignature(tx.Transaction, tx.V, tx.R, tx.S)
}

// ToSignedTransaction converts the mempool transaction to a signed transaction.
func (tx MempoolTransaction) ToSignedTransaction() *SignedTransaction {
	return &SignedTransaction{
		Transaction: tx.Transaction,
		V:           tx.V,
		R:           tx.R,
		S:           tx.S,
	}
}

// MarshalBSON overrides the default serializer to handle a bug with mongodb.
func (tx *MempoolTransaction) MarshalBSON() ([]byte, error) {
	// Developers note:
	// The reason *big.Int fields (like V, R, and S in MempoolTransaction) aren't
	// being saved in MongoDB is because MongoDB's bson package does not natively
	// support encoding or decoding *big.Int values. By default, MongoDB doesn't
	// know how to handle big.Int types, so they end up being ignored.
	//
	// To fix this, you need to manually convert *big.Int values to a format
	// MongoDB can store (such as a string or integer) and then convert them back
	// on retrieval. Here’s one way to achieve this by adding custom serialization
	// for these fields

	type Alias MempoolTransaction // Alias to avoid recursion
	return bson.Marshal(&struct {
		V string `bson:"v"`
		R string `bson:"r"`
		S string `bson:"s"`
		*Alias
	}{
		V:     tx.V.String(),
		R:     tx.R.String(),
		S:     tx.S.String(),
		Alias: (*Alias)(tx),
	})
}

// UnmarshalBSON overrides the default deserializer to handle a bug with mongodb.
func (tx *MempoolTransaction) UnmarshalBSON(data []byte) error {
	// Developers note:
	// The reason *big.Int fields (like V, R, and S in MempoolTransaction) aren't
	// being saved in MongoDB is because MongoDB's bson package does not natively
	// support encoding or decoding *big.Int values. By default, MongoDB doesn't
	// know how to handle big.Int types, so they end up being ignored.
	//
	// To fix this, you need to manually convert *big.Int values to a format
	// MongoDB can store (such as a string or integer) and then convert them back
	// on retrieval. Here’s one way to achieve this by adding custom serialization
	// for these fields

	type Alias MempoolTransaction // Alias to avoid recursion
	aux := &struct {
		V string `bson:"v"`
		R string `bson:"r"`
		S string `bson:"s"`
		*Alias
	}{
		Alias: (*Alias)(tx),
	}

	if err := bson.Unmarshal(data, aux); err != nil {
		return err
	}

	tx.V, _ = new(big.Int).SetString(aux.V, 10)
	tx.R, _ = new(big.Int).SetString(aux.R, 10)
	tx.S, _ = new(big.Int).SetString(aux.S, 10)

	return nil
}
