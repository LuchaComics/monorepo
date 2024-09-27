package db

import (
	"fmt"
	"log"
	"log/slog"
	"path/filepath"

	"github.com/syndtr/goleveldb/leveldb"
	dberr "github.com/syndtr/goleveldb/leveldb/errors"
	"github.com/syndtr/goleveldb/leveldb/util"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/pkg/db"
)

type keyValueStorerImpl struct {
	db *leveldb.DB
}

const dbDirName = "blocks"

func GetDBDirPath(dataDir string) string {
	return filepath.Join(dataDir, dbDirName)
}

func NewDatabase(cfg *config.Config, logger *slog.Logger) db.Database {
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

func (impl *keyValueStorerImpl) Get(prefix, key string) ([]byte, error) {
	return impl.Getf("%s-%s", prefix, key)
}

func (impl *keyValueStorerImpl) Set(prefix, key string, val []byte) error {
	return impl.Setf(val, "%s-%s", prefix, key)
}

func (impl *keyValueStorerImpl) Delete(prefix, key string) error {
	return impl.Deletef("%s-%s", prefix, key)
}
func (impl *keyValueStorerImpl) Getf(format string, a ...any) ([]byte, error) {
	k := fmt.Sprintf(format, a...)
	bin, err := impl.db.Get([]byte(k), nil)
	if err == dberr.ErrNotFound {
		return nil, nil
	}
	return bin, nil
}

func (impl *keyValueStorerImpl) Setf(val []byte, format string, a ...any) error {
	k := fmt.Sprintf(format, a...)
	impl.db.Delete([]byte(k), nil)
	err := impl.db.Put([]byte(k), val, nil)
	if err == dberr.ErrNotFound {
		return nil
	}
	return err
}

func (impl *keyValueStorerImpl) Deletef(format string, a ...any) error {
	k := fmt.Sprintf(format, a...)
	err := impl.db.Delete([]byte(k), nil)
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

// Iterate function used to provide a list of key and values for your code
// to iterate through.
func (impl *keyValueStorerImpl) Iterate(keyPrefix string, seekThenIterateKey string, processFunc func(key, value []byte) error) error {
	iter := impl.db.NewIterator(util.BytesPrefix([]byte(keyPrefix)), nil)

	// Apply filter, else do not.
	if seekThenIterateKey == "" {
		if ok := iter.First(); !ok {
			return nil
		}
	} else {
		if ok := iter.Seek([]byte(seekThenIterateKey)); !ok {
			return nil
		}
	}

	for iter.Next() {
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
