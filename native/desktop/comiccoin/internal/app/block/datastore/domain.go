package datastore

import "time"

type Block struct {
	Hash         string         `json:"hash"`
	PreviousHash string         `json:"previous_hash"`
	Timestamp    time.Time      `json:"timestamp"`
	Nonce        uint64         `json:"nonce"`
	Difficulty   int            `json:"difficulty"`
	Transactions []*Transaction `json:"transactions"`
}
