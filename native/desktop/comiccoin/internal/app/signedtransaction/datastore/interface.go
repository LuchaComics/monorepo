package datastore

import (
	"context"
	"log/slog"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/adapter/keyvaluestore"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/config"
)

type SignedTransactionStorer interface {
	Insert(ctx context.Context, b *SignedTransaction) error
	GetByNonce(ctx context.Context, nonce uint64) (*SignedTransaction, error)
	List(ctx context.Context) ([]*SignedTransaction, error)
	DeleteByNonce(ctx context.Context, nonce uint64) error
}

type signedTansactionStorerImpl struct {
	logger   *slog.Logger
	dbClient keyvaluestore.KeyValueStorer
}

func NewDatastore(cfg *config.Config, logger *slog.Logger, kvs keyvaluestore.KeyValueStorer) SignedTransactionStorer {
	return &signedTansactionStorerImpl{
		dbClient: kvs,
		logger:   logger,
	}
}
