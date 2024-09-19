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
}

func NewBlockchainWithCoinbaseKey(
	cfg *config.Config,
	kvs keyvaluestore.KeyValueStorer,
	coinbaseKey *keystore.Key,
) *Blockchain {
	bc := &Blockchain{
		Difficulty: cfg.BlockchainDifficulty,
		Database:   kvs,
	}

	// Defensive code.
	if cfg.BlockchainDifficulty <= 0 {
		log.Fatal("cannot have blochain difficulty less then or equal to zero")
	}

	// Check if we have a genesis block. If we do not have it then we will need
	// to generate it here in our initialization.
	lastHashBin, err := kvs.Get([]byte("LAST_HASH"))
	if err != nil || lastHashBin == nil || string(lastHashBin) == "" {
		// No existing blockchain found, create genesis block
		genesisBlock := NewGenesisBlock(coinbaseKey)
		genesisBlock.Mine(bc.Difficulty)

		// Store genesis block in our database, to begin, we need to
		// serialize it into []bytes array.
		blockData, err := json.Marshal(genesisBlock)
		if err != nil {
			log.Fatalf("Failed to marshal genesis block: %v", err)
		}

		// Save our genesis hash.
		blockKey := fmt.Sprintf("BLOCK_%v", genesisBlock.Hash)
		err = kvs.Set([]byte(blockKey), blockData)
		if err != nil {
			log.Fatalf("Failed to store genesis block: %v", err)
		}

		// Save our last hash
		err = kvs.Set([]byte("LAST_HASH"), []byte(genesisBlock.Hash))
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
	lastHashBin, err := kvs.Get([]byte("LAST_HASH"))
	if err != nil || lastHashBin == nil || string(lastHashBin) == "" {
		if err != nil {
			log.Fatalf("failed initializing blockchain with error: %v", err)
		}
		log.Fatalf("failed initializing blockchain because of missing genesis block")
	} else {
		bc.LastHash = string(lastHashBin)
		log.Println("fetched latest hash:", bc.LastHash)
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

	newBlock := &Block{
		PreviousHash: lastHash,
		Timestamp:    time.Now(),
		Difficulty:   bc.Difficulty,
		Transactions: transactions,
	}

	newBlock.Mine(bc.Difficulty)

	// Defensive code for code sanity check.
	if newBlock.Hash == "" {
		log.Fatal("cannot have empty newBlock.Hash!")
	}

	// Store new block
	blockData, err := json.Marshal(newBlock)
	if err != nil {
		return fmt.Errorf("failed to marshal new block: %v", err)
	}

	blockKey := fmt.Sprintf("BLOCK_%v", newBlock.Hash)
	err = bc.Database.Set([]byte(blockKey), blockData)
	if err != nil {
		return fmt.Errorf("failed to store new block: %v", err)
	}

	// Update last hash
	err = bc.Database.Set([]byte("LAST_HASH"), []byte(newBlock.Hash))
	if err != nil {
		return fmt.Errorf("failed to update last hash: %v", err)
	}

	bc.LastHash = newBlock.Hash
	return nil
}

func (bc *Blockchain) GetBalance(address common.Address) (*big.Int, error) {
	balance := new(big.Int)
	currentHash := bc.LastHash

	for {
		blockKey := fmt.Sprintf("BLOCK_%v", currentHash)
		blockData, err := bc.Database.Get([]byte(blockKey))
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
