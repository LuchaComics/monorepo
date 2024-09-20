package datastore

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

func (b *Block) CalculateHash() string {
	data, _ := json.Marshal(b.Transactions)
	blockData := b.PreviousHash + string(data) + b.Timestamp.String() + strconv.FormatUint(b.Nonce, 10) + strconv.Itoa(b.Difficulty)
	blockHash := sha256.Sum256([]byte(blockData))
	return fmt.Sprintf("%x", blockHash)
}

func (b *Block) Mine(difficulty int) (uint64, string) {
	// log.Println("Mine: b.Hash (premine):", b.Hash)
	for !strings.HasPrefix(b.Hash, strings.Repeat("0", difficulty)) {
		b.Nonce++
		b.Hash = b.CalculateHash()
	}
	return b.Nonce, b.Hash
}
