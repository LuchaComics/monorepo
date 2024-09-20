package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"strconv"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
)

type Block struct {
	Hash         string         `json:"hash"`
	PreviousHash string         `json:"previous_hash"`
	Timestamp    time.Time      `json:"timestamp"`
	Nonce        uint64         `json:"nonce"`
	Difficulty   int            `json:"difficulty"`
	Transactions []*Transaction `json:"transactions"`
}

func NewGenesisBlock(coinbaseKey *keystore.Key) *Block {
	initialSupply, _ := new(big.Int).SetString("50000000000000000000", 10) // 50 coins with 18 decimal places
	genesisTransaction := NewTransaction(
		common.Address{},    // From: zero address for genesis block
		coinbaseKey.Address, // To: coinbase address (usually the miner's address)
		initialSupply,       // 50 coins (with 18 decimal places)
		[]byte("Genesis Block"),
		0, // Nonce
	)

	if err := genesisTransaction.Sign(coinbaseKey.PrivateKey); err != nil {
		log.Fatalf("Failed to sign transaction: %v", err)
	}

	return &Block{
		Hash:         "",
		PreviousHash: "",
		Timestamp:    time.Now(),
		Nonce:        0,
		Difficulty:   1,
		Transactions: []*Transaction{genesisTransaction},
	}
}

func (b *Block) CalculateHash() string {
	data, _ := json.Marshal(b.Transactions)
	blockData := b.PreviousHash + string(data) + b.Timestamp.String() + strconv.FormatUint(b.Nonce, 10) + strconv.Itoa(b.Difficulty)
	blockHash := sha256.Sum256([]byte(blockData))
	return fmt.Sprintf("%x", blockHash)
}

func (b *Block) Mine(difficulty int) {
	// log.Println("Mine: b.Hash (premine):", b.Hash)
	for !strings.HasPrefix(b.Hash, strings.Repeat("0", difficulty)) {
		b.Nonce++
		b.Hash = b.CalculateHash()
	}
	// log.Println("Mine: b.Hash (postmine):", b.Hash)
}

func (b *Block) Serialize() []byte {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)
	err := encoder.Encode(b)
	if err != nil {
		log.Fatalf("failed to serialize block: %v", err)
	}
	return result.Bytes()
}

func DeserializeBlock(data []byte) (*Block, error) {
	block := &Block{}
	err := json.Unmarshal(data, block)
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize block: %v", err)
	}
	return block, nil

	// DEVELOPERS NOTE: Depercated code.
	// var block Block
	// decoder := gob.NewDecoder(bytes.NewReader(data))
	// err := decoder.Decode(&block)
	// if err != nil {
	// 	log.Fatalf("failed to deserialize block: %v", err)
	// }
	// return &block
}
