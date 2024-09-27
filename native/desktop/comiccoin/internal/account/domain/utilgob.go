package domain

import (
	"bytes"
	"encoding/gob"
	"fmt"
)

func (b *Account) Serialize() ([]byte, error) {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)
	err := encoder.Encode(b)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize account: %v", err)
	}
	return result.Bytes(), nil
}

func NewAccountFromDeserialize(data []byte) (*Account, error) {
	// Variable we will use to return.
	account := &Account{}

	// Defensive code: If programmer entered empty bytes then we will
	// return nil deserialization result.
	if data == nil {
		return nil, nil
	}

	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&account)
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize account: %v", err)
	}
	return account, nil
}
