package domain

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/fxamacker/cbor/v2"
)

// Known Issue:
// MongoDB does not natively support storing `*big.Int` types directly,
// which may cause issues if the `V`, `R`, and `S` fields of the
// `SignedTransaction` are of type `*big.Int`.
// These fields need to be converted to an alternate type (e.g., `string`
// or `[]byte`) before saving to MongoDB, as `big.Int` values are not
// natively serialized by BSON.
// To resolve this, we convert `*big.Int` fields to `[]byte` in the model.

// BlockHeader represents common information required for each block.
type BlockHeader struct {
	ChainID        uint16         `bson:"chain_id" json:"chain_id"`               // Keep track of which chain this block belongs to.
	NumberBytes    []byte         `bson:"number_bytes" json:"number_bytes"`       // Ethereum: Block number in the chain.
	PrevBlockHash  string         `bson:"prev_block_hash" json:"prev_block_hash"` // Bitcoin: Hash of the previous block in the chain.
	TimeStamp      uint64         `bson:"timestamp" json:"timestamp"`             // Bitcoin: Time the block was mined.
	Difficulty     uint16         `bson:"difficulty" json:"difficulty"`           // Ethereum: Number of 0's needed to solve the hash solution.
	Beneficiary    common.Address `bson:"beneficiary" json:"beneficiary"`         // Ethereum: The account who is receiving fees .
	TransactionFee uint64         `bson:"transaction_fee" json:"transaction_fee"` // ComicCoin: Fee that must be paid for every transaction. This value is provided by the authority.

	// The StateRoot represents a hash of the in-memory account balance
	// database. This field allows the blockchain to provide a guarantee that
	// the accounting of the transactions and fees for each account on each
	// node is exactly the same.
	StateRoot string `bson:"state_root" json:"state_root"` // Ethereum: Represents a hash of the accounts and their balances.

	TransRoot  string `bson:"trans_root" json:"trans_root"`   // Both: Represents the merkle tree root hash for the transactions in this block.
	NonceBytes []byte `bson:"nonce_bytes" json:"nonce_bytes"` // Both: Value identified to solve the hash solution.

	LatestTokenIDBytes []byte `bson:"latest_token_id_bytes" json:"latest_token_id_bytes"` // ComicCoin: The latest token that the blockchain points to.
	TokensRoot         string `bson:"tokens_root" json:"tokens_root"`                     // ComicCoin: Represents the hash of all the tokens and their owners.
}

func (bh *BlockHeader) GetNumber() *big.Int {
	return new(big.Int).SetBytes(bh.NumberBytes)
}

func (bh *BlockHeader) IsNumberZero() bool {
	// Special thanks to "Is there another way of testing if a big.Int is 0?" via https://stackoverflow.com/a/64257532
	return len(bh.GetNumber().Bits()) == 0
}

func (bh *BlockHeader) SetNumber(n *big.Int) {
	bh.NumberBytes = n.Bytes()
}

func (bh *BlockHeader) GetNonce() *big.Int {
	return new(big.Int).SetBytes(bh.NonceBytes)
}

func (bh *BlockHeader) IsNonceZero() bool {
	// Special thanks to "Is there another way of testing if a big.Int is 0?" via https://stackoverflow.com/a/64257532
	return len(bh.GetNonce().Bits()) == 0
}

func (bh *BlockHeader) SetNonce(n *big.Int) {
	bh.NonceBytes = n.Bytes()
}

func (bh *BlockHeader) GetLatestTokenID() *big.Int {
	return new(big.Int).SetBytes(bh.LatestTokenIDBytes)
}

func (bh *BlockHeader) IsLatestTokenIDZero() bool {
	// Special thanks to "Is there another way of testing if a big.Int is 0?" via https://stackoverflow.com/a/64257532
	return len(bh.GetLatestTokenID().Bits()) == 0
}

func (bh *BlockHeader) SeLatestTokenID(n *big.Int) {
	bh.LatestTokenIDBytes = n.Bytes()
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
