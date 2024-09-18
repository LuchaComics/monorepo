package keyvaluestore

// KeyValueStorer interface used to implement actions you can take onto a
// key-value database which stores the data in a persistence manner.
type KeyValueStorer interface {
	Get(key []byte) ([]byte, error)
	Set(key []byte, val []byte) error
	Delete(key []byte) error
}
