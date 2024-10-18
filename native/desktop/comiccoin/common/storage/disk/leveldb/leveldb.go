package db

import (
	"log"
	"log/slog"

	"github.com/syndtr/goleveldb/leveldb"
	dberr "github.com/syndtr/goleveldb/leveldb/errors"
	"github.com/syndtr/goleveldb/leveldb/filter"
	"github.com/syndtr/goleveldb/leveldb/opt"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/common/storage"
)

// storageImpl implements the db.Database interface.
// It uses a LevelDB database to store key-value pairs.
type storageImpl struct {
	// The LevelDB database instance.
	db *leveldb.DB
}

// NewDiskStorage creates a new instance of the storageImpl.
// It opens the database file at the specified path and returns an error if it fails.
func NewDiskStorage(dbPath string, dbName string, logger *slog.Logger) storage.Storage {
	if dbPath == "" {
		log.Fatal("cannot have empty filepath for the database")
	}

	o := &opt.Options{
		Filter: filter.NewBloomFilter(10),
	}

	filePath := dbPath + "/" + dbName

	db, err := leveldb.OpenFile(filePath, o)
	if err != nil {
		log.Fatalf("failed loading up key value storer adapter at %v", filePath)
	}
	return &storageImpl{
		db: db,
	}
}

// Get retrieves a value from the database by its key.
// It returns an error if the key is not found.
func (impl *storageImpl) Get(k string) ([]byte, error) {
	bin, err := impl.db.Get([]byte(k), nil)
	if err == dberr.ErrNotFound {
		return nil, nil
	}
	return bin, nil
}

// Set sets a value in the database by its key.
// It returns an error if the operation fails.
func (impl *storageImpl) Set(k string, val []byte) error {
	impl.db.Delete([]byte(k), nil)
	err := impl.db.Put([]byte(k), val, nil)
	if err == dberr.ErrNotFound {
		return nil
	}
	return err
}

// Delete deletes a value from the database by its key.
// It returns an error if the operation fails.
func (impl *storageImpl) Delete(k string) error {
	err := impl.db.Delete([]byte(k), nil)
	if err == dberr.ErrNotFound {
		return nil
	}
	return err
}

// Iterate iterates over the key-value pairs in the database, starting from the specified key prefix.
// It calls the provided function for each pair.
// It returns an error if the iteration fails.
func (impl *storageImpl) Iterate(processFunc func(key, value []byte) error) error {
	iter := impl.db.NewIterator(nil, nil)
	for ok := iter.First(); ok; ok = iter.Next() {
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
func (impl *storageImpl) Close() error {
	return impl.db.Close()
}