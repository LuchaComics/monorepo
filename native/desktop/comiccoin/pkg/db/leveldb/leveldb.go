package db

import (
	"fmt"
	"log"
	"log/slog"
	"path/filepath"

	"github.com/syndtr/goleveldb/leveldb"
	dberr "github.com/syndtr/goleveldb/leveldb/errors"
	"github.com/syndtr/goleveldb/leveldb/util"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/pkg/db"
)

// keyValueStorerImpl implements the db.Database interface.
// It uses a LevelDB database to store key-value pairs.
type keyValueStorerImpl struct {
	// The LevelDB database instance.
	db *leveldb.DB
}

// dbDirName is the name of the directory where the database is stored.
const dbDirName = "db"

// GetDBDirPath returns the path to the database directory.
func GetDBDirPath(dataDir string) string {
	return filepath.Join(dataDir, dbDirName)
}

// NewDatabase creates a new instance of the keyValueStorerImpl.
// It opens the database file at the specified path and returns an error if it fails.
func NewDatabase(dataDir string, logger *slog.Logger) db.Database {
	if dataDir == "" {
		log.Fatal("cannot have empty dir")
	}
	db, err := leveldb.OpenFile(GetDBDirPath(dataDir), nil)
	if err != nil {
		log.Fatalf("failed loading up key value storer adapter at %v", dataDir)
	}
	return &keyValueStorerImpl{
		db: db,
	}
}

// Get retrieves a value from the database by its key.
// It returns an error if the key is not found.
func (impl *keyValueStorerImpl) Get(prefix, key string) ([]byte, error) {
	return impl.Getf("%s-%s", prefix, key)
}

// Set sets a value in the database by its key.
// It returns an error if the operation fails.
func (impl *keyValueStorerImpl) Set(prefix, key string, val []byte) error {
	return impl.Setf(val, "%s-%s", prefix, key)
}

// Delete deletes a value from the database by its key.
// It returns an error if the operation fails.
func (impl *keyValueStorerImpl) Delete(prefix, key string) error {
	return impl.Deletef("%s-%s", prefix, key)
}

// Getf retrieves a value from the database by its key.
// It returns an error if the key is not found.
func (impl *keyValueStorerImpl) Getf(format string, a ...any) ([]byte, error) {
	k := fmt.Sprintf(format, a...)
	bin, err := impl.db.Get([]byte(k), nil)
	if err == dberr.ErrNotFound {
		return nil, nil
	}
	return bin, nil
}

// Setf sets a value in the database by its key.
// It returns an error if the operation fails.
func (impl *keyValueStorerImpl) Setf(val []byte, format string, a ...any) error {
	k := fmt.Sprintf(format, a...)
	impl.db.Delete([]byte(k), nil)
	err := impl.db.Put([]byte(k), val, nil)
	if err == dberr.ErrNotFound {
		return nil
	}
	return err
}

// Deletef deletes a value from the database by its key.
// It returns an error if the operation fails.
func (impl *keyValueStorerImpl) Deletef(format string, a ...any) error {
	k := fmt.Sprintf(format, a...)
	err := impl.db.Delete([]byte(k), nil)
	if err == dberr.ErrNotFound {
		return nil
	}
	return err
}

// View iterates over the key-value pairs in the database and calls the provided function for each pair.
// It returns an error if the iteration fails.
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

// ViewFromFirst iterates over the key-value pairs in the database, starting from the first pair.
// It calls the provided function for each pair.
// It returns an error if the iteration fails.
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

// Iterate iterates over the key-value pairs in the database, starting from the specified key prefix.
// It calls the provided function for each pair.
// It returns an error if the iteration fails.
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

// Close closes the database.
// It returns an error if the operation fails.
func (impl *keyValueStorerImpl) Close() error {
	return impl.db.Close()
}
