package domain

import (
	"bytes"
	"encoding/gob"
	"fmt"
)

// PendingBlockTransaction represents the transaction before it was recorded to
// block because the transaction needs to be mined by the minging system.
type PendingBlockTransaction struct {
	MempoolTransaction
	TimeStamp uint64 `json:"timestamp"` // Ethereum: The time the transaction was received.
	GasPrice  uint64 `json:"gas_price"` // Ethereum: The price of one unit of gas to be paid for fees.
	GasUnits  uint64 `json:"gas_units"` // Ethereum: The number of units of gas used for this transaction.
}

type PendingBlockTransactionRepository interface {
	Upsert(bd *PendingBlockTransaction) error
	ListAll() ([]*PendingBlockTransaction, error)
	DeleteAll() error
}

func (dto *PendingBlockTransaction) Serialize() ([]byte, error) {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)
	err := encoder.Encode(dto)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize block data: %v", err)
	}
	return result.Bytes(), nil
}

func NewPendingBlockTransactionFromDeserialize(data []byte) (*PendingBlockTransaction, error) {
	// Variable we will use to return.
	dto := &PendingBlockTransaction{}

	// Defensive code: If programmer entered empty bytes then we will
	// return nil deserialization result.
	if data == nil {
		return nil, nil
	}

	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&dto)
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize block data: %v", err)
	}
	return dto, nil
}
