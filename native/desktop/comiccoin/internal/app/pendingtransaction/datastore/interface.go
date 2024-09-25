package datastore

import (
	"context"
	"log/slog"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/adapter/keyvaluestore"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/config"
)

type PendingTransactionStorer interface {
	Insert(ctx context.Context, b *SignedPendingTransaction) error
	GetByNonce(ctx context.Context, nonce uint64) (*SignedPendingTransaction, error)
	List(ctx context.Context) ([]*SignedPendingTransaction, error)
	DeleteByNonce(ctx context.Context, nonce uint64) error
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
