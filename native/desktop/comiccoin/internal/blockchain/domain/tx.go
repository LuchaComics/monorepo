package domain

import "github.com/ethereum/go-ethereum/common"

// Transaction structure represents a transfer of coins between accounts
// which have not been added to the blockchain yet and are waiting for the miner
// to receive and verify. Once  transactions have been veriried
// they will be deleted from our system as they will live in the blockchain
// afterwords.
type Transaction struct {
	ChainID uint16         `json:"chain_id"` // Ethereum: The chain id that is listed in the genesis file.
	Nonce   uint64         `json:"nonce"`    // Ethereum: Unique id for the transaction supplied by the user.
	From    common.Address `json:"from"`     // Ethereum: Account sending the transaction. Will be checked against signature.
	To      common.Address `json:"to"`       // Ethereum: Account receiving the benefit of the transaction.
	Value   uint64         `json:"value"`    // Ethereum: Monetary value received from this transaction.
	Tip     uint64         `json:"tip"`      // Ethereum: Tip offered by the sender as an incentive to mine this transaction.
	Data    []byte         `json:"data"`     // Ethereum: Extra data related to the transaction.
}
