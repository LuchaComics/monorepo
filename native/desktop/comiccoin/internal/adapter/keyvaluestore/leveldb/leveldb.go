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

const dbDirName = "blocks"

func GetDBDirPath(dataDir string) string {
	return filepath.Join(dataDir, dbDirName)
}

func NewKeyValueStorer(cfg *config.Config) keyvaluestore.KeyValueStorer {
	if cfg.DB.DataDir == "" {
		log.Fatal("cannot have empty dir")
	}
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
	impl.db.Delete(key, nil)
	return impl.db.Put(key, val, nil)
}

func (impl *keyValueStorerImpl) Delete(key []byte) error {
	return impl.db.Delete(key, nil)
}

func (impl *keyValueStorerImpl) View(key []byte, processFunc func(key, value []byte) error) error {
	iter := impl.db.NewIterator(nil, nil)
	for ok := iter.Seek(key); ok; ok = iter.Next() {
		// Call the passed function for each key-value pair.
		err := processFunc(iter.Key(), iter.Value())
		if err != nil {
			return err // Exit early if the processing function returns an error.
		}
	}
	iter.Release()
	return iter.Error()
}

func (impl *keyValueStorerImpl) ViewFromFirst(processFunc func(key, value []byte) error) error {
	iter := impl.db.NewIterator(nil, nil)
	for ok := iter.First(); ok; ok = iter.Next() {
		// Call the passed function for each key-value pair.
		err := processFunc(iter.Key(), iter.Value())
		if err != nil {
			return err // Exit early if the processing function returns an error.
		}
	}
	iter.Release()
	return iter.Error()
}

func (impl *keyValueStorerImpl) Close() error {
	return impl.db.Close()
}
