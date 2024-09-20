package leveldb

import (
	"fmt"
	"log"
	"log/slog"
	"path/filepath"

	"github.com/syndtr/goleveldb/leveldb"
	dberr "github.com/syndtr/goleveldb/leveldb/errors"

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

func NewKeyValueStorer(cfg *config.Config, logger *slog.Logger) keyvaluestore.KeyValueStorer {
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

func (impl *keyValueStorerImpl) Get(key string) ([]byte, error) {
	bin, err := impl.db.Get([]byte(key), nil)
	if err == dberr.ErrNotFound {
		return nil, nil
	}
	return bin, nil
}

func (impl *keyValueStorerImpl) Getf(format string, a ...any) ([]byte, error) {
	k := fmt.Sprintf(format, a...)
	return impl.Get(k)
}

func (impl *keyValueStorerImpl) Set(key string, val []byte) error {
	impl.db.Delete([]byte(key), nil)
	err := impl.db.Put([]byte(key), val, nil)
	if err == dberr.ErrNotFound {
		return nil
	}
	return err
}

func (impl *keyValueStorerImpl) Setf(val []byte, format string, a ...any) error {
	k := fmt.Sprintf(format, a...)
	return impl.Set(k, val)
}

func (impl *keyValueStorerImpl) Delete(key string) error {
	err := impl.db.Delete([]byte(key), nil)
	if err == dberr.ErrNotFound {
		return nil
	}
	return err
}

func (impl *keyValueStorerImpl) View(key string, processFunc func(key, value []byte) error) error {
	iter := impl.db.NewIterator(nil, nil)
	for ok := iter.Seek([]byte(key)); ok; ok = iter.Next() {
		// Call the passed function for each key-value pair.
		err := processFunc(iter.Key(), iter.Value())
		if err == dberr.ErrNotFound {
			return nil
		}
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
		log.Println("ViewFromFirst: key:", iter.Key(), "val:", iter.Value())
		// Call the passed function for each key-value pair.
		err := processFunc(iter.Key(), iter.Value())
		if err == dberr.ErrNotFound {
			return nil
		}
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
