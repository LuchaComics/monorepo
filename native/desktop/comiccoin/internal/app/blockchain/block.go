package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

type Block struct {
	Data         map[string]interface{} `bson:"Data" json:"Data"`
	Hash         string                 `bson:"Hash" json:"Hash"`
	PreviousHash string                 `bson:"previous_Hash" json:"previous_Hash"`
	Timestamp    time.Time              `bson:"Timestamp" json:"Timestamp"`
	POW          int                    `bson:"POW" json:"POW"`
}

func (b Block) calculateHash() string {
	data, _ := json.Marshal(b.Data)
	blockData := b.PreviousHash + string(data) + b.Timestamp.String() + strconv.Itoa(b.POW)
	blockHash := sha256.Sum256([]byte(blockData))
	return fmt.Sprintf("%x", blockHash)
}

func (b *Block) mine(difficulty int) {
	for !strings.HasPrefix(b.Hash, strings.Repeat("0", difficulty)) {
		b.POW++
		b.Hash = b.calculateHash()
	}
}

func Handle(err error) {
	if err != nil {
		log.Fatalf("failed to serialize or deserialize: %v", err)
	}
}

func (b *Block) Serialize() []byte {
	var res bytes.Buffer
	encoder := gob.NewEncoder(&res)

	err := encoder.Encode(b)

	Handle(err)

	return res.Bytes()
}

func Deserialize(data []byte) *Block {
	var block Block

	decoder := gob.NewDecoder(bytes.NewReader(data))

	err := decoder.Decode(&block)

	Handle(err)

	return &block
}
