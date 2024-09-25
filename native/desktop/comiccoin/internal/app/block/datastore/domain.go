package datastore

import "time"

// Genesis represents the genesis file.
type Genesis struct {
	Date          time.Time         `json:"date"`
	ChainID       uint16            `json:"chain_id"`        // The chain id represents an unique id for this running instance.
	TransPerBlock uint16            `json:"trans_per_block"` // The maximum number of transactions that can be in a block.
	Difficulty    uint16            `json:"difficulty"`      // How difficult it needs to be to solve the work problem.
	MiningReward  uint64            `json:"mining_reward"`   // Reward for mining a block.
	GasPrice      uint64            `json:"gas_price"`       // Fee paid for each transaction mined into a block.
	Balances      map[string]uint64 `json:"balances"`
}

type Block struct {
	Hash         string         `json:"hash"`
	PreviousHash string         `json:"previous_hash"`
	Timestamp    time.Time      `json:"timestamp"`
	Nonce        uint64         `json:"nonce"`
	Difficulty   int            `json:"difficulty"`
	Transactions []*Transaction `json:"transactions"`
}
