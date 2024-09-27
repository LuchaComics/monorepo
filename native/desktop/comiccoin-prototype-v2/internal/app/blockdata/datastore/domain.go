package datastore

import (
	signedtx "github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/app/signedtransaction/datastore"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/utils/blockchain/merkle"
	"github.com/ethereum/go-ethereum/common"
)

// =============================================================================

// BlockData represents what can be serialized to disk and over the network.
type BlockData struct {
	Hash   string      `json:"hash"`
	Header BlockHeader `json:"block"`
	Trans  []BlockTx   `json:"trans"`
}

// =============================================================================

// BlockHeader represents common information required for each block.
type BlockHeader struct {
	Number        uint64         `json:"number"`          // Ethereum: Block number in the chain.
	PrevBlockHash string         `json:"prev_block_hash"` // Bitcoin: Hash of the previous block in the chain.
	TimeStamp     uint64         `json:"timestamp"`       // Bitcoin: Time the block was mined.
	Beneficiary   common.Address `json:"beneficiary"`     // Ethereum: The account who is receiving fees and tips.
	Difficulty    uint16         `json:"difficulty"`      // Ethereum: Number of 0's needed to solve the hash solution.
	MiningReward  uint64         `json:"mining_reward"`   // Ethereum: The reward for mining this block.
	StateRoot     string         `json:"state_root"`      // Ethereum: Represents a hash of the accounts and their balances.
	TransRoot     string         `json:"trans_root"`      // Both: Represents the merkle tree root hash for the transactions in this block.
	Nonce         uint64         `json:"nonce"`           // Both: Value identified to solve the hash solution.
}

// Block represents a group of transactions batched together.
type Block struct {
	Header     BlockHeader
	MerkleTree *merkle.Tree[BlockTx]
}

// POWArgs represents the set of arguments required to run POW.
type POWArgs struct {
	Beneficiary  common.Address
	Difficulty   uint16
	MiningReward uint64
	PrevBlock    Block
	StateRoot    string
	Trans        []BlockTx
	EvHandler    func(v string, args ...any)
}

// =============================================================================

// BlockTx represents the transaction as it's recorded inside a block. This
// includes a timestamp and gas fees.
type BlockTx struct {
	signedtx.SignedTransaction
	TimeStamp uint64 `json:"timestamp"` // Ethereum: The time the transaction was received.
	GasPrice  uint64 `json:"gas_price"` // Ethereum: The price of one unit of gas to be paid for fees.
	GasUnits  uint64 `json:"gas_units"` // Ethereum: The number of units of gas used for this transaction.
}
