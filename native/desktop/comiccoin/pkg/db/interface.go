package db

// Database interface used to implement actions you can take onto a
// key-value database which stores the data in a persistence manner.
type Database interface {
	Get(prefix, key string) ([]byte, error)
	Set(prefix, key string, val []byte) error
	Delete(prefix, key string) error

	Getf(format string, a ...any) ([]byte, error)
	Setf(val []byte, format string, a ...any) error
	Deletef(format string, a ...any) error

	View(key string, processFunc func(key, value []byte) error) error
	ViewFromFirst(processFunc func(key, value []byte) error) error
	Iterate(keyPrefix string, seekThenIterateKey string, processFunc func(key, value []byte) error) error

	Close() error
}
