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
	LastHash    string
	Difficulty  int
	Database    keyvaluestore.KeyValueStorer
	CoinbaseKey *keystore.Key
}

func NewBlockchain(
	cfg *config.Config,
	kvs keyvaluestore.KeyValueStorer,
	coinbaseKey *keystore.Key,
) *Blockchain {
	bc := &Blockchain{
		Difficulty:  cfg.BlockchainDifficulty,
		Database:    kvs,
		CoinbaseKey: coinbaseKey,
	}

	// Defensive code.
	if cfg.BlockchainDifficulty <= 0 {
		log.Fatal("cannot have blochain difficulty less then or equal to zero")
	}

	// Check if we have a genesis block
	lastHashBytes, err := kvs.Get([]byte("LAST_HASH"))
	if err != nil || lastHashBytes == nil || string(lastHashBytes) == "" {
		log.Println("err:", err)
		// No existing blockchain found, create genesis block
		genesisBlock := NewGenesisBlock(coinbaseKey)
		genesisBlock.Hash = genesisBlock.CalculateHash()

		// Store genesis block
		blockData, err := json.Marshal(genesisBlock)
		if err != nil {
			log.Fatalf("Failed to marshal genesis block: %v", err)
		}
		err = kvs.Set([]byte(genesisBlock.Hash), blockData)
		if err != nil {
			log.Fatalf("Failed to store genesis block: %v", err)
		}

		// Update last hash
		err = kvs.Set([]byte("LAST_HASH"), []byte(genesisBlock.Hash))
		if err != nil {
			log.Fatalf("Failed to store last hash: %v", err)
		}

		bc.LastHash = genesisBlock.Hash
		log.Println("generated first hash:", bc.LastHash)
	} else {
		bc.LastHash = string(lastHashBytes)
		log.Println("fetched latest hash:", bc.LastHash)
	}
	return bc
}

func (bc *Blockchain) Close() error {
	return bc.Database.Close()
}

func (bc *Blockchain) AddBlock(transactions []*Transaction) error {
	lastHash := bc.LastHash
	if lastHash == "" {
		log.Fatal("cannot have empty last hash!")
	} else {
		log.Println("AddBlock: lastHash:", lastHash)
	}

	newBlock := &Block{
		PreviousHash: lastHash,
		Timestamp:    time.Now(),
		Difficulty:   bc.Difficulty,
		Transactions: transactions,
	}

	newBlock.Mine(bc.Difficulty)

	if newBlock.Hash == "" {
		log.Fatal("cannot have empty newBlock.Hash!")
	} else {
		log.Println("AddBlock: newBlock.Hash:", newBlock.Hash)
	}

	// Store new block
	blockData, err := json.Marshal(newBlock)
	if err != nil {
		return fmt.Errorf("failed to marshal new block: %v", err)
	}
	err = bc.Database.Set([]byte(newBlock.Hash), blockData)
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
		blockData, err := bc.Database.Get([]byte(currentHash))
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
