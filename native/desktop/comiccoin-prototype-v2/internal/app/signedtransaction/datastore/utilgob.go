package datastore

import (
	"bytes"
	"encoding/gob"
	"fmt"
)

func (b *SignedTransaction) Serialize() ([]byte, error) {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)
	err := encoder.Encode(b)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize Transaction: %v", err)
	}
	return result.Bytes(), nil
}

func NewSignedTransactionFromDeserialize(data []byte) (*SignedTransaction, error) {
	// Variable we will use to return.
	Transaction := &SignedTransaction{}

	// Defensive code: If programmer entered empty bytes then we will
	// return nil deserialization result.
	if data == nil {
		return nil, nil
	}

	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&Transaction)
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize Transaction: %v", err)
	}
	return Transaction, nil
}
