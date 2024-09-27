package datastore

import (
	"context"
	"log/slog"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/adapter/keyvaluestore"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/config"
)

type BlockDataStorer interface {
	Insert(ctx context.Context, b *BlockData) error
	GetByHash(ctx context.Context, hash string) (*BlockData, error)
}

type blockdataStorerImpl struct {
	logger   *slog.Logger
	dbClient keyvaluestore.KeyValueStorer
}

func NewDatastore(cfg *config.Config, logger *slog.Logger, kvs keyvaluestore.KeyValueStorer) BlockDataStorer {
	return &blockdataStorerImpl{
		dbClient: kvs,
		logger:   logger,
	}
}
