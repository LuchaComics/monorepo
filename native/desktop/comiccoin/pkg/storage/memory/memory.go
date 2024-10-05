package memory

import (
	"fmt"
	"log/slog"
	"strings"
	"sync"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/pkg/storage"
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
	impl.lock.Lock()
	defer impl.lock.Unlock()
	key := fmt.Sprintf(format, a...)

	cachedValue, ok := impl.data[key]
	if !ok {
		delete(impl.data, key)
		return nil, fmt.Errorf("does not exist for: %v", key)
	}

	return cachedValue.value, nil
}

// Setf sets a value in the database by its key.
// It returns an error if the operation fails.
func (impl *keyValueStorerImpl) Setf(val []byte, format string, a ...any) error {
	impl.lock.Lock()
	defer impl.lock.Unlock()
	key := fmt.Sprintf(format, a...)

	impl.data[key] = cacheValue{
		value: val,
	}
	return nil
}

// Deletef deletes a value from the database by its key.
// It returns an error if the operation fails.
func (impl *keyValueStorerImpl) Deletef(format string, a ...any) error {
	impl.lock.Lock()
	defer impl.lock.Unlock()
	key := fmt.Sprintf(format, a...)

	delete(impl.data, key)
	return nil
}

// View iterates over the key-value pairs in the database and calls the provided function for each pair.
// It returns an error if the iteration fails.
func (impl *keyValueStorerImpl) View(key string, processFunc func(key, value []byte) error) error {
	impl.lock.Lock()
	defer impl.lock.Unlock()

	// Iterate over the key-value pairs in the database
	for k, v := range impl.data {
		// Check if the key matches the provided key
		if strings.HasPrefix(k, key) {
			// Call the provided function for each pair
			if err := processFunc([]byte(k), v.value); err != nil {
				return err
			}
		}
	}

	return nil
}

// ViewFromFirst iterates over the key-value pairs in the database, starting from the first pair.
// It calls the provided function for each pair.
// It returns an error if the iteration fails.
func (impl *keyValueStorerImpl) ViewFromFirst(processFunc func(key, value []byte) error) error {
	impl.lock.Lock()
	defer impl.lock.Unlock()

	// Iterate over the key-value pairs in the database, starting from the first pair
	for k, v := range impl.data {
		// Call the provided function for each pair
		if err := processFunc([]byte(k), v.value); err != nil {
			return err
		}
	}

	return nil
}

// Iterate iterates over the key-value pairs in the database, starting from the specified key prefix.
// It calls the provided function for each pair.
// It returns an error if the iteration fails.
func (impl *keyValueStorerImpl) Iterate(keyPrefix string, seekThenIterateKey string, processFunc func(key, value []byte) error) error {
	impl.lock.Lock()
	defer impl.lock.Unlock()

	// Find the starting point for the iteration
	var found bool
	var startingKey string
	for k := range impl.data {
		if strings.HasPrefix(k, keyPrefix) && k >= seekThenIterateKey {
			found = true
			startingKey = k
			break
		}
	}

	// If the starting point is not found, return an error
	if !found {
		return fmt.Errorf("starting point not found for key prefix %s and seek key %s", keyPrefix, seekThenIterateKey)
	}

	// Iterate over the key-value pairs in the database, starting from the starting point
	for k, v := range impl.data {
		if strings.HasPrefix(k, keyPrefix) && k >= startingKey {
			// Call the provided function for each pair
			if err := processFunc([]byte(k), v.value); err != nil {
				return err
			}
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
