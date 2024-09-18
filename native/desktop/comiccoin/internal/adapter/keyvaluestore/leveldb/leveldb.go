package leveldb

import (
	"log"
	"path/filepath"

	"github.com/syndtr/goleveldb/leveldb"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/adapter/keyvaluestore"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/config"
)

type keyValueStorerImpl struct {
	db *leveldb.DB
}

const dbDirName = "blockchain"

func GetDBDirPath(dataDir string) string {
	return filepath.Join(dataDir, dbDirName)
}

func NewKeyValueStorer(cfg *config.Config) keyvaluestore.KeyValueStorer {
	db, err := leveldb.OpenFile(GetDBDirPath(cfg.DB.DataDir), nil)
	if err != nil {
		log.Fatal("failed loading up key value storer adapter")
	}
	return &keyValueStorerImpl{
		db: db,
	}
}

func (impl *keyValueStorerImpl) Get(key []byte) ([]byte, error) {
	return impl.db.Get(key, nil)
}

func (impl *keyValueStorerImpl) Set(key []byte, val []byte) error {
	return impl.db.Put(key, val, nil)
}

func (impl *keyValueStorerImpl) Delete(key []byte) error {
	return impl.db.Delete(key, nil)
}
