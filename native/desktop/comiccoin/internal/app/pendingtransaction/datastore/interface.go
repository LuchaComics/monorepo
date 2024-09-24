package datastore

import (
	"context"
	"log/slog"
	"math/big"

	"github.com/ethereum/go-ethereum/common"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/adapter/keyvaluestore"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/config"
)

// PendingTransaction structure represents a transfer of coins between accounts
// which have not been added to the blockchain yet and are waiting for the miner
// to receive and verify. Once pending transactions have been veriried
// they will be deleted from our system as they will live in the blockchain
// afterwords.
type PendingTransaction struct {
	ID        string         `json:"id"`
	From      common.Address `json:"from"`
	To        common.Address `json:"to"`
	Value     *big.Int       `json:"value"`
	Data      []byte         `json:"data"`
	Nonce     uint64         `json:"nonce"`
	Signature []byte         `json:"signature"`
}

type PendingTransactionStorer interface {
	Insert(ctx context.Context, b *PendingTransaction) error
	GetByID(ctx context.Context, id string) (*PendingTransaction, error)
	List(ctx context.Context) ([]*PendingTransaction, error)
	DeleteByID(ctx context.Context, id string) error
}

type pendingTransactionStorerImpl struct {
	logger   *slog.Logger
	dbClient keyvaluestore.KeyValueStorer
}

func NewDatastore(cfg *config.Config, logger *slog.Logger, kvs keyvaluestore.KeyValueStorer) PendingTransactionStorer {
	return &pendingTransactionStorerImpl{
		dbClient: kvs,
		logger:   logger,
	}
}
