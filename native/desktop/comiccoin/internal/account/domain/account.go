package domain

import (
	"bytes"
	"encoding/gob"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
)

// Account struct represents a users wallets that contain keys (private/public).
// Note: The keys remain encrypted at rest and requires user password input to
// decrypt so the account can be used for transfering coins.
// Note: Private key - Used to sign transactions
// Note: Public key - Acts as the user’s address.
type Account struct {
	ID             string         `json:"id"`
	WalletFilepath string         `json:"wallet_filepath"`
	WalletAddress  common.Address `json:"wallet_address"`
}

type AccountRepository interface {
	Upsert(acc *Account) error
	GetByID(id string) (*Account, error)
	List() ([]*Account, error)
	DeleteByID(id string) error
}

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
