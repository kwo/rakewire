package model

// DatabaseConfiguration configuration for service.
type DatabaseConfiguration struct {
	Location string
}

// Database defines the interface to a key-value store
type Database interface {
	Location() string
	Select(fn func(tx Transaction) error) error
	Update(fn func(tx Transaction) error) error
	Repair() error
}

// Transaction represents an atomic operation to the database
type Transaction interface {
	Bucket(name string) Bucket
}

// Bucket holds key-values
type Bucket interface {
	Bucket(name string) Bucket
	Cursor() Cursor
	Delete(key []byte) error
	ForEach(fn func(key, value []byte) error) error
	Get(key []byte) []byte
	NextSequence() (uint64, error)
	Put(key []byte, value []byte) error
}

// Cursor loops through values in a bucket
type Cursor interface {
	Delete() error
	First() (key []byte, value []byte)
	Last() (key []byte, value []byte)
	Next() (key []byte, value []byte)
	Prev() (key []byte, value []byte)
	Seek(seek []byte) (key []byte, value []byte)
}
