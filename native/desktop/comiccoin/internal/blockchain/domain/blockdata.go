package domain

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

// =============================================================================

// BlockData represents what can be serialized to disk and over the network.
type BlockData struct {
	Hash   string              `json:"hash"`
	Header *BlockHeader        `json:"block_header"`
	Trans  []*BlockTransaction `json:"trans"`
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

// SignedTransaction is a signed version of the transaction. This is how
// clients like a wallet provide transactions for inclusion into the blockchain.
type SignedTransaction struct {
	Transaction
	V *big.Int `json:"v"` // Ethereum: Recovery identifier, either 29 or 30 with ardanID.
	R *big.Int `json:"r"` // Ethereum: First coordinate of the ECDSA signature.
	S *big.Int `json:"s"` // Ethereum: Second coordinate of the ECDSA signature.
}

// BlockTransaction represents the transaction as it's recorded inside a block. This
// includes a timestamp and gas fees.
type BlockTransaction struct {
	SignedTransaction
	TimeStamp uint64 `json:"timestamp"` // Ethereum: The time the transaction was received.
	GasPrice  uint64 `json:"gas_price"` // Ethereum: The price of one unit of gas to be paid for fees.
	GasUnits  uint64 `json:"gas_units"` // Ethereum: The number of units of gas used for this transaction.
}

// =============================================================================

func (b *BlockData) Serialize() ([]byte, error) {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)
	err := encoder.Encode(b)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize block data: %v", err)
	}
	return result.Bytes(), nil
}

func NewBlockDataFromDeserialize(data []byte) (*BlockData, error) {
	// Variable we will use to return.
	account := &BlockData{}

	// Defensive code: If programmer entered empty bytes then we will
	// return nil deserialization result.
	if data == nil {
		return nil, nil
	}

	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&account)
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize block data: %v", err)
	}
	return account, nil
}

// =============================================================================

type BlockDataRepository interface {
	Upsert(bd *BlockData) error
	GetByHash(hash string) (*BlockData, error)
	List() ([]*BlockData, error)
	DeleteByHash(hash string) error
}
