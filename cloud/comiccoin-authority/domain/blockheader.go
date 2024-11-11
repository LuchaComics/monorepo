package domain

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/fxamacker/cbor/v2"
)

// BlockHeader represents common information required for each block.
type BlockHeader struct {
	ChainID       uint16         `bson:"chain_id" json:"chain_id"`               // Keep track of which chain this block belongs to.
	Number        uint64         `bson:"number" json:"number"`                   // Ethereum: Block number in the chain.
	PrevBlockHash string         `bson:"prev_block_hash" json:"prev_block_hash"` // Bitcoin: Hash of the previous block in the chain.
	TimeStamp     uint64         `bson:"timestamp" json:"timestamp"`             // Bitcoin: Time the block was mined.
	Beneficiary   common.Address `bson:"beneficiary" json:"beneficiary"`         // Ethereum: The account who is receiving fees and tips.
	Difficulty    uint16         `bson:"difficulty" json:"difficulty"`           // Ethereum: Number of 0's needed to solve the hash solution.
	MiningReward  uint64         `bson:"mining_reward" json:"mining_reward"`     // Ethereum: The reward for mining this block.

	// The StateRoot represents a hash of the in-memory account balance
	// database. This field allows the blockchain to provide a guarantee that
	// the accounting of the transactions and fees for each account on each
	// node is exactly the same.
	StateRoot string `json:"state_root"` // Ethereum: Represents a hash of the accounts and their balances.

	TransRoot string `json:"trans_root"` // Both: Represents the merkle tree root hash for the transactions in this block.
	Nonce     uint64 `json:"nonce"`      // Both: Value identified to solve the hash solution.

	LatestTokenID uint64 `json:"latest_token_id"` // ComicCoin: The latest token that the blockchain points to.
	TokensRoot    string `json:"tokens_root"`     // ComicCoin: Represents the hash of all the tokens and their owners.
}

// Serialize serializes a block header into a byte array.
// It returns the serialized byte array and an error if one occurs.
func (b *BlockHeader) Serialize() ([]byte, error) {
	// Marshal the block data into a byte array using CBOR.
	dataBytes, err := cbor.Marshal(b)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize block header: %v", err)
	}
	return dataBytes, nil
}
