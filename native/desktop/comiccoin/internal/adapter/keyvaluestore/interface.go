package keyvaluestore

// KeyValueStorer interface used to implement actions you can take onto a
// key-value database which stores the data in a persistence manner.
type KeyValueStorer interface {
	Get(key string) ([]byte, error)
	Getf(format string, a ...any) ([]byte, error)
	Set(key string, val []byte) error
	Setf(val []byte, format string, a ...any) error
	Delete(key string) error
	View(key string, processFunc func(key, value []byte) error) error
	ViewFromFirst(processFunc func(key, value []byte) error) error
	Close() error
}
