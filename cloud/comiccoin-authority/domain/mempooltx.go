package domain

import (
	"context"
	"crypto/ecdsa"
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/common/blockchain/signature"
	"github.com/fxamacker/cbor/v2"
)

// MempoolTransaction represents a transaction that is stored in the mempool.
// It contains the transaction data, as well as the ECDSA signature and recovery identifier.
type MempoolTransaction struct {
	ID primitive.ObjectID `bson:"_id" json:"id"`

	// The signed transaction data, including sender, recipient, amount, and other metadata like V, R, S, etc.
	SignedTransaction
}

// Validate checks if the transaction is valid.
// It verifies the signature, makes sure the account addresses are correct,
// and checks if the 'from' and 'to' accounts are not the same.
func (mtx MempoolTransaction) Validate(chainID uint16, isPoA bool) error {
	return mtx.SignedTransaction.Validate(chainID, isPoA)
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

	// DeleteByChainID deletes all mempool transactions in the repository for the particular chainID.
	DeleteByChainID(ctx context.Context, chainID uint16) error

	DeleteByID(ctx context.Context, id primitive.ObjectID) error

	GetInsertionChangeStreamChannel(ctx context.Context) (<-chan MempoolTransaction, chan struct{}, error)
}

// Serialize serializes the mempool transaction into a byte slice.
// This method uses the cbor library to marshal the transaction into a byte slice.
func (mtx *MempoolTransaction) Serialize() ([]byte, error) {
	// Marshal the transaction into a byte slice using the cbor library.
	dataBytes, err := cbor.Marshal(mtx)
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
	mtx := &MempoolTransaction{}

	// Defensive code: If the input data is empty, return a nil deserialization result.
	if data == nil {
		return nil, nil
	}

	// Unmarshal the byte slice into the transaction variable using the cbor library.
	if err := cbor.Unmarshal(data, &mtx); err != nil {
		// Return an error if the unmarshaling fails.
		return nil, fmt.Errorf("failed to deserialize mempool transaction: %v", err)
	}
	return mtx, nil
}

// FromAddress extracts the account address from the signed transaction by
// recovering the public key from the signature.
func (mtx MempoolTransaction) FromAddress() (string, error) {
	// Note: MongoDB doesn't support `*big.Int` so we are forced to do this.
	v, r, s := mtx.SignedTransaction.GetBigIntFields()

	return signature.FromAddress(mtx.SignedTransaction.Transaction, v, r, s)
}

func (mtx MempoolTransaction) FromPublicKey() (*ecdsa.PublicKey, error) {
	// Note: MongoDB doesn't support `*big.Int` so we are forced to do this.
	v, r, s := mtx.SignedTransaction.GetBigIntFields()

	return signature.GetPublicKeyFromSignature(mtx.SignedTransaction.Transaction, v, r, s)
}
