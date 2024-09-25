package datastore

import (
	"bytes"
	"encoding/gob"
	"fmt"
)

func (b *SignedPendingTransaction) Serialize() ([]byte, error) {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)
	err := encoder.Encode(b)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize pendingTransaction: %v", err)
	}
	return result.Bytes(), nil
}

func NewSignedPendingTransactionFromDeserialize(data []byte) (*SignedPendingTransaction, error) {
	// Variable we will use to return.
	pendingTransaction := &SignedPendingTransaction{}

	// Defensive code: If programmer entered empty bytes then we will
	// return nil deserialization result.
	if data == nil {
		return nil, nil
	}

	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&pendingTransaction)
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize pendingTransaction: %v", err)
	}
	return pendingTransaction, nil
}
