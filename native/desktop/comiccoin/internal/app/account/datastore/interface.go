package datastore

import (
	"context"
	"log/slog"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/adapter/keyvaluestore"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/config"
)

type Account struct {
	Name           string         `json:"name"`
	WalletFilepath string         `json:"wallet_filepath"`
	WalletAddress  common.Address `json:"wallet_address"`
}

type AccountStorer interface {
	Insert(ctx context.Context, accountName, accountWalletPassword string) (*Account, error)
	GetByName(ctx context.Context, name string) (*Account, error)
	List(ctx context.Context) ([]*Account, error)
	DeleteByName(ctx context.Context, name string) error
	GetKeyByNameAndPassword(ctx context.Context, accountName string, accountWalletPassword string) (*keystore.Key, error)
}

type accountStorerImpl struct {
	config   *config.Config
	logger   *slog.Logger
	dbClient keyvaluestore.KeyValueStorer
}

func NewDatastore(cfg *config.Config, logger *slog.Logger, kvs keyvaluestore.KeyValueStorer) AccountStorer {
	return &accountStorerImpl{
		config:   cfg,
		dbClient: kvs,
		logger:   logger,
	}
}
