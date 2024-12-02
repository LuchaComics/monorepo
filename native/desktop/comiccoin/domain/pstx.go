package domain

import (
	"fmt"
	"math/big"

	"github.com/fxamacker/cbor/v2"

	auth_domain "github.com/LuchaComics/monorepo/cloud/comiccoin-authority/domain"
)

// PendingSignedTransaction struct is responsible for storing all the
// transactions that the user, of this current local wallet, submitted to the
// authority and is waiting for this transaction to be either added (or rejected)
// by the authority.
//
// Use this struct to prevent the user, of this current local wallet, to make
// any transactions unless there are no more
type PendingSignedTransaction auth_domain.SignedTransaction

// SignedTransactionToPendingSignedTransaction method converts a
// `SignedTransaction` data type into a `PendingSignedTransaction` data type.
func SignedTransactionToPendingSignedTransaction(pstx *auth_domain.SignedTransaction) *PendingSignedTransaction {
	return (*PendingSignedTransaction)(pstx)
}

func PendingSignedTransactionToSignedTransaction(pstx *PendingSignedTransaction) *auth_domain.SignedTransaction {
	return (*auth_domain.SignedTransaction)(pstx)
}

const (
	PendingSignedTransactionStateNotReady = "not_ready"
	PendingSignedTransactionStateReady    = "ready"
)

type PendingSignedTransactionRepository interface {
	// Upsert inserts or updates an nfts in the repository.
	Upsert(acc *PendingSignedTransaction) error

	// GetByID retrieves an nft by its token id.
	GetByNonce(nonce *big.Int) (*PendingSignedTransaction, error)

	// ListAll retrieves all nfts in the repository.
	ListAll() ([]*PendingSignedTransaction, error)

	// // ListByIDs retrieves nfts that have the token is in the parameter
	// ListWithFilterByIDs(tokIDs []*big.Int) ([]*PendingSignedTransaction, error)

	// DeleteByNonce deletes a pending signed transaction based by a nonce value.
	DeleteByNonce(pstxNonce *big.Int) error

	OpenTransaction() error

	CommitTransaction() error

	DiscardTransaction()
}

// Serialize serializes the wallet into a byte slice.
// This method uses the cbor library to marshal the wallet into a byte slice.
func (b *PendingSignedTransaction) Serialize() ([]byte, error) {
	// Marshal the wallet into a byte slice using the cbor library.
	dataBytes, err := cbor.Marshal(b)
	if err != nil {
		// Return an error if the marshaling fails.
		return nil, fmt.Errorf("failed to serialize wallet: %v", err)
	}
	return dataBytes, nil
}

// NewPendingSignedTransactionFromDeserialize deserializes an wallet from a byte slice.
// This method uses the cbor library to unmarshal the byte slice into an wallet.
func NewPendingSignedTransactionFromDeserialize(data []byte) (*PendingSignedTransaction, error) {
	// Create a new wallet variable to return.
	wallet := &PendingSignedTransaction{}

	// Defensive code: If the input data is empty, return a nil deserialization result.
	if data == nil {
		return nil, nil
	}

	// Unmarshal the byte slice into the wallet variable using the cbor library.
	if err := cbor.Unmarshal(data, &wallet); err != nil {
		// Return an error if the unmarshaling fails.
		return nil, fmt.Errorf("failed to deserialize wallet: %v", err)
	}
	return wallet, nil
}

func ToPendingSignedTransactionNoncesArray(txs []*PendingSignedTransaction) []*big.Int {
	txNonces := make([]*big.Int, len(txs))
	for _, tx := range txs {
		txNonces = append(txNonces, new(big.Int).SetBytes(tx.Transaction.NonceBytes))
	}
	return txNonces
}
