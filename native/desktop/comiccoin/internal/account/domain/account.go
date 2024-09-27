package domain

import (
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
