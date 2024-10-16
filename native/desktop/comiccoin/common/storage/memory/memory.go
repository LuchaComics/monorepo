package memory

import (
	"fmt"
	"log/slog"
	"sync"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/common/storage"
)

type cacheValue struct {
	value []byte
}

// keyValueStorerImpl implements the db.Database interface.
// It uses a LevelDB database to store key-value pairs.
type keyValueStorerImpl struct {
	data map[string]cacheValue
	lock sync.Mutex
}

// NewInMemoryStorage creates a new instance of the keyValueStorerImpl.
func NewInMemoryStorage(logger *slog.Logger) storage.Storage {
	return &keyValueStorerImpl{
		data: make(map[string]cacheValue),
	}
}

// Get retrieves a value from the database by its key.
// It returns an error if the key is not found.
func (impl *keyValueStorerImpl) Get(k string) ([]byte, error) {
	impl.lock.Lock()
	defer impl.lock.Unlock()

	cachedValue, ok := impl.data[k]
	if !ok {
		delete(impl.data, k)
		return nil, fmt.Errorf("does not exist for: %v", k)
	}

	return cachedValue.value, nil
}

// Set sets a value in the database by its key.
// It returns an error if the operation fails.
func (impl *keyValueStorerImpl) Set(k string, val []byte) error {
	impl.lock.Lock()
	defer impl.lock.Unlock()

	impl.data[k] = cacheValue{
		value: val,
	}
	return nil
}

// Delete deletes a value from the database by its key.
// It returns an error if the operation fails.
func (impl *keyValueStorerImpl) Delete(k string) error {
	impl.lock.Lock()
	defer impl.lock.Unlock()

	delete(impl.data, k)
	return nil
}

// Iterate iterates over the key-value pairs in the database, starting from the specified key prefix.
// It calls the provided function for each pair.
// It returns an error if the iteration fails.
func (impl *keyValueStorerImpl) Iterate(processFunc func(key, value []byte) error) error {
	impl.lock.Lock()
	defer impl.lock.Unlock()

	// Iterate over the key-value pairs in the database, starting from the starting point
	for k, v := range impl.data {
		// Call the provided function for each pair
		if err := processFunc([]byte(k), v.value); err != nil {
			return err
		}
	}

	return nil
}

// Close closes the database.
// It returns an error if the operation fails.
func (impl *keyValueStorerImpl) Close() error {
	impl.lock.Lock()
	defer impl.lock.Unlock()

	// Clear the data map
	impl.data = make(map[string]cacheValue)

	return nil
}
