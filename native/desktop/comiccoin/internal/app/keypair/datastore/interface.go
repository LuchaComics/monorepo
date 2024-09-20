package datastore

import (
	"context"
	"log/slog"

	"github.com/libp2p/go-libp2p/core/crypto"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/adapter/keyvaluestore"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/config"
)

type KeypairStorer interface {
	GetByName(ctx context.Context, name string) (crypto.PrivKey, crypto.PubKey, error)
	GenerateNewKeyPairAndSetByName(ctx context.Context, name string) error
}

type keypairStorerImpl struct {
	logger   *slog.Logger
	dbClient keyvaluestore.KeyValueStorer
}

func NewDatastore(cfg *config.Config, logger *slog.Logger, kvs keyvaluestore.KeyValueStorer) KeypairStorer {
	return &keypairStorerImpl{
		dbClient: kvs,
		logger:   logger,
	}
}
