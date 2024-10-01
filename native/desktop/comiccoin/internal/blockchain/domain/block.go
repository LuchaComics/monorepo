package domain

import (
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/pkg/blockchain/merkle"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/pkg/blockchain/signature"
)

// Block represents a group of transactions batched together.
type Block struct {
	Header     *BlockHeader
	MerkleTree *merkle.Tree[BlockTransaction]
}

// NewBlockData constructs block data from a block.
func NewBlockData(block Block) *BlockData {
	blockData := BlockData{
		Hash:   block.Hash(),
		Header: block.Header,
		Trans:  block.MerkleTree.Values(),
	}

	return &blockData
}

// Hash returns the unique hash for the Block.
func (b Block) Hash() string {
	if b.Header.Number == 0 {
		return signature.ZeroHash
	}

	// CORE NOTE: Hashing the block header and not the whole block so the blockchain
	// can be cryptographically checked by only needing block headers and not full
	// blocks with the transaction data. This will support the ability to have pruned
	// nodes and light clients in the future.
	// - A pruned node stores all the block headers, but only a small number of full
	//   blocks (maybe the last 1000 blocks). This allows for full cryptographic
	//   validation of blocks and transactions without all the extra storage.
	// - A light client keeps block headers and just enough sufficient information
	//   to follow the latest set of blocks being produced. The do not validate
	//   blocks, but can prove a transaction is in a block.

	return signature.Hash(b.Header)
}

// isHashSolved checks the hash to make sure it complies with
// the POW rules. We need to match a difficulty number of 0's.
func isHashSolved(difficulty uint16, hash string) bool {
	const match = "0x00000000000000000"

	if len(hash) != 66 {
		return false
	}

	difficulty += 2
	return hash[:difficulty] == match[:difficulty]
}

// // ValidateBlock takes a block and validates it to be included into the blockchain.
// func (b Block) ValidateBlock(previousBlock Block, stateRoot string, evHandler func(v string, args ...any)) error {
// 	evHandler("database: ValidateBlock: validate: blk[%d]: check: chain is not forked", b.Header.Number)
//
// 	// The node who sent this block has a chain that is two or more blocks ahead
// 	// of ours. This means there has been a fork and we are on the wrong side.
// 	nextNumber := previousBlock.Header.Number + 1
// 	if b.Header.Number >= (nextNumber + 2) {
// 		return ErrChainForked
// 	}
//
// 	evHandler("database: ValidateBlock: validate: blk[%d]: check: block difficulty is the same or greater than parent block difficulty", b.Header.Number)
//
// 	if b.Header.Difficulty < previousBlock.Header.Difficulty {
// 		return fmt.Errorf("block difficulty is less than previous block difficulty, parent %d, block %d", previousBlock.Header.Difficulty, b.Header.Difficulty)
// 	}
//
// 	evHandler("database: ValidateBlock: validate: blk[%d]: check: block hash has been solved", b.Header.Number)
//
// 	hash := b.Hash()
// 	if !isHashSolved(b.Header.Difficulty, hash) {
// 		return fmt.Errorf("%s invalid block hash", hash)
// 	}
//
// 	evHandler("database: ValidateBlock: validate: blk[%d]: check: block number is the next number", b.Header.Number)
//
// 	if b.Header.Number != nextNumber {
// 		return fmt.Errorf("this block is not the next number, got %d, exp %d", b.Header.Number, nextNumber)
// 	}
//
// 	evHandler("database: ValidateBlock: validate: blk[%d]: check: parent hash does match parent block", b.Header.Number)
//
// 	if b.Header.PrevBlockHash != previousBlock.Hash() {
// 		return fmt.Errorf("parent block hash doesn't match our known parent, got %s, exp %s", b.Header.PrevBlockHash, previousBlock.Hash())
// 	}
//
// 	if previousBlock.Header.TimeStamp > 0 {
// 		evHandler("database: ValidateBlock: validate: blk[%d]: check: block's timestamp is greater than parent block's timestamp", b.Header.Number)
//
// 		parentTime := time.Unix(int64(previousBlock.Header.TimeStamp), 0)
// 		blockTime := time.Unix(int64(b.Header.TimeStamp), 0)
// 		if blockTime.Before(parentTime) {
// 			return fmt.Errorf("block timestamp is before parent block, parent %s, block %s", parentTime, blockTime)
// 		}
//
// 		// This is a check that Ethereum does but we can't because we don't run all the time.
//
// 		// evHandler("database: ValidateBlock: validate: blk[%d]: check: block is less than 15 minutes apart from parent block", b.Header.Number)
//
// 		// dur := blockTime.Sub(parentTime)
// 		// if dur.Seconds() > time.Duration(15*time.Second).Seconds() {
// 		// 	return fmt.Errorf("block is older than 15 minutes, duration %v", dur)
// 		// }
// 	}
//
// 	evHandler("database: ValidateBlock: validate: blk[%d]: check: state root hash does match current database", b.Header.Number)
//
// 	if b.Header.StateRoot != stateRoot {
// 		return fmt.Errorf("state of the accounts are wrong, current %s, expected %s", stateRoot, b.Header.StateRoot)
// 	}
//
// 	evHandler("database: ValidateBlock: validate: blk[%d]: check: merkle root does match transactions", b.Header.Number)
//
// 	if b.Header.TransRoot != b.MerkleTree.RootHex() {
// 		return fmt.Errorf("merkle root does not match transactions, got %s, exp %s", b.MerkleTree.RootHex(), b.Header.TransRoot)
// 	}
//
// 	return nil
// }
