package blockchain

import (
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/adapter/keyvaluestore"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/config"
)

type Blockchain struct {
	LastHash   string
	Difficulty int
	Database   keyvaluestore.KeyValueStorer

	// Channel to send new blocks or results to peer nodes
	resultCh chan *Block
}

func NewBlockchainWithCoinbaseKey(
	cfg *config.Config,
	kvs keyvaluestore.KeyValueStorer,
	coinbaseKey *keystore.Key,
) *Blockchain {
	bc := &Blockchain{
		Difficulty: cfg.BlockchainDifficulty,
		Database:   kvs,
		resultCh:   make(chan *Block),
	}

	// Defensive code.
	if cfg.BlockchainDifficulty <= 0 {
		log.Fatal("cannot have blochain difficulty less then or equal to zero")
	}

	// Check if we have a genesis block. If we do not have it then we will need
	// to generate it here in our initialization.
	lastHashBin, err := kvs.Getf("LAST_HASH")
	if err != nil || lastHashBin == nil || string(lastHashBin) == "" {
		// No existing blockchain found, create genesis block
		genesisBlock := NewGenesisBlock(coinbaseKey)
		genesisBlock.Mine(bc.Difficulty)

		// Store genesis block in our database, to begin, we need to
		// serialize it into []bytes array.
		genesisBlockBin, err := json.Marshal(genesisBlock)
		if err != nil {
			log.Fatalf("Failed to marshal genesis block: %v", err)
		}
		if genesisBlockBin == nil { // Defensive code for programmer.
			log.Fatal("did not marshal genesis block")
		}

		// Save our genesis hash.
		err = kvs.Setf(genesisBlockBin, "BLOCK_%v", genesisBlock.Hash)
		if err != nil {
			log.Fatalf("Failed to store genesis block: %v", err)
		}

		// Save our last hash
		err = kvs.Setf([]byte(genesisBlock.Hash), "LAST_HASH")
		if err != nil {
			log.Fatalf("Failed to store last hash: %v", err)
		}

		bc.LastHash = genesisBlock.Hash
	} else {
		bc.LastHash = string(lastHashBin)
	}
	return bc
}

func NewBlockchain(
	cfg *config.Config,
	kvs keyvaluestore.KeyValueStorer,
) *Blockchain {
	bc := &Blockchain{
		Difficulty: cfg.BlockchainDifficulty,
		Database:   kvs,
	}

	// Defensive code.
	if cfg.BlockchainDifficulty <= 0 {
		log.Fatal("cannot have blochain difficulty less then or equal to zero")
	}

	// Check for our latest block
	lastHashBin, err := kvs.Getf("LAST_HASH")
	if err != nil || lastHashBin == nil || string(lastHashBin) == "" {
		if err != nil {
			log.Fatalf("failed initializing blockchain with error: %v", err)
		}
		log.Fatalf("failed initializing blockchain because of missing genesis block")
	} else {
		bc.LastHash = string(lastHashBin)
	}
	return bc
}

func (bc *Blockchain) Close() error {
	return bc.Database.Close()
}

func (bc *Blockchain) AddBlock(transactions []*Transaction) error {
	// Defensive code to protect the programmer.
	lastHash := bc.LastHash
	if lastHash == "" {
		log.Fatal("cannot have empty last hash!")
	}
	log.Println("fetching last block")
	// Fetch the last known block to compare with our newly created block.
	oldBlockBin, err := bc.Database.Getf("BLOCK_%v", lastHash)
	if err != nil {
		return fmt.Errorf("failed to lookup `BLOCK_%v` in database: %v", lastHash, err)
	}
	log.Println("fetched last block:", string(oldBlockBin))
	oldBlock := DeserializeBlock(oldBlockBin)
	log.Println("deserialized last block")

	newBlock := &Block{
		PreviousHash: lastHash,
		Timestamp:    time.Now(),
		Difficulty:   bc.Difficulty,
		Transactions: transactions,
	}

	newBlock.Mine(bc.Difficulty)

	// Defensive code for programmer code sanity check.
	if newBlock.Hash == "" {
		log.Fatal("cannot have empty newBlock.Hash!")
	}

	if isBlockValid(newBlock, oldBlock) {
		// Store new block
		blockData, err := json.Marshal(newBlock)
		if err != nil {
			return fmt.Errorf("failed to marshal new block: %v", err)
		}

		err = bc.Database.Setf(blockData, "BLOCK_%v", newBlock.Hash)
		if err != nil {
			return fmt.Errorf("failed to store new block: %v", err)
		}

		// Update last hash
		err = bc.Database.Setf([]byte(newBlock.Hash), "LAST_HASH")
		if err != nil {
			return fmt.Errorf("failed to update last hash: %v", err)
		}

		bc.LastHash = newBlock.Hash
	}

	return nil
}

func (bc *Blockchain) GetBalance(address common.Address) (*big.Int, error) {
	balance := new(big.Int)
	currentHash := bc.LastHash

	for {
		blockData, err := bc.Database.Getf("BLOCK_%v", currentHash)
		if err != nil {
			return nil, fmt.Errorf("failed to get block data: %v", err)
		}

		var block Block
		err = json.Unmarshal(blockData, &block)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal block data: %v", err)
		}

		for _, tx := range block.Transactions {
			if tx.From == address {
				balance.Sub(balance, tx.Value)
			}
			if tx.To == address {
				balance.Add(balance, tx.Value)
			}
		}

		if block.PreviousHash == "" {
			break // Genesis block reached
		}
		currentHash = block.PreviousHash
	}

	return balance, nil
}

func isBlockValid(newBlock, oldBlock *Block) bool {
	// if oldBlock.Index+1 != newBlock.Index {
	// 	return false
	// }

	if oldBlock.Hash != newBlock.PreviousHash {
		return false
	}

	if newBlock.CalculateHash() != newBlock.Hash {
		return false
	}

	var txVerified bool = true
	for _, tx := range newBlock.Transactions {
		txVerified = txVerified && tx.Verify()
	}

	return txVerified
}

// Subscribe returns the channel so you can listen for new results.
func (bc *Blockchain) Subscribe() <-chan *Block {
	return bc.resultCh
}
