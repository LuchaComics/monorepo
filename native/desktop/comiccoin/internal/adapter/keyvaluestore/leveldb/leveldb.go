package leveldb

import (
	"log"

	"github.com/syndtr/goleveldb/leveldb"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/adapter/keyvaluestore"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/config"
)

type keyValueStorerImpl struct {
	db *leveldb.DB
}

func NewKeyValueStorer(cfg *config.Config) keyvaluestore.KeyValueStorer {
	db, err := leveldb.OpenFile(cfg.DB.DataDir, nil)
	if err != nil {
		log.Fatal("failed loading up key value storer adapter")
	}
	return &keyValueStorerImpl{
		db: db,
	}
}

func (impl *keyValueStorerImpl) Get(key []byte) ([]byte, error) {
	return nil, nil
}

func (impl *keyValueStorerImpl) Set(key []byte, val []byte) error {
	return nil
}

func (impl *keyValueStorerImpl) Delete(key []byte) error {
	return nil
}
