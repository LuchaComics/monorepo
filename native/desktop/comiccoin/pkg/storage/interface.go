package storage

// Storage interface defines the methods that can be used to interact with a key-value database.
type Storage interface {
	// Get returns the value associated with the specified key, or an error if the key is not found.
	Get(prefix, key string) ([]byte, error)

	// Set sets the value associated with the specified key.
	// If the key already exists, its value is updated.
	Set(prefix, key string, val []byte) error

	// Delete removes the value associated with the specified key from the database.
	Delete(prefix, key string) error

	// Getf is a variant of Get that allows the key to be constructed using a format string.
	Getf(format string, a ...any) ([]byte, error)

	// Setf is a variant of Set that allows the key to be constructed using a format string.
	Setf(val []byte, format string, a ...any) error

	// Deletef is a variant of Delete that allows the key to be constructed using a format string.
	Deletef(format string, a ...any) error

	// View iterates over the key-value pairs in the database, starting from the specified key.
	// The processFunc callback is called for each pair, and can return an error to terminate the iteration.
	View(key string, processFunc func(key, value []byte) error) error

	// ViewFromFirst is similar to View, but starts from the first key in the database.
	ViewFromFirst(processFunc func(key, value []byte) error) error

	// Iterate is similar to View, but allows the iteration to start from a specific key prefix.
	// The seekThenIterateKey parameter can be used to specify a key to seek to before starting the iteration.
	Iterate(keyPrefix string, seekThenIterateKey string, processFunc func(key, value []byte) error) error

	// Close closes the database, releasing any system resources it holds.
	Close() error
}
