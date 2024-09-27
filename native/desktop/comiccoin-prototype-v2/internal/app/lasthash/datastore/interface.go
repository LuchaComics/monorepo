package datastore

import (
	"context"
	"log/slog"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/adapter/keyvaluestore"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/config"
)

type LastHashStorer interface {
	Get(ctx context.Context) (string, error)
	Set(ctx context.Context, hash string) error
}

type lastHashStorerImpl struct {
	logger   *slog.Logger
	dbClient keyvaluestore.KeyValueStorer
}

func NewDatastore(cfg *config.Config, logger *slog.Logger, kvs keyvaluestore.KeyValueStorer) LastHashStorer {
	return &lastHashStorerImpl{
		dbClient: kvs,
		logger:   logger,
	}
}
