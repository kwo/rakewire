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
}

// Transaction represents an atomic operation to the database
type Transaction interface {
	Bucket(name string) Bucket
	Container(name string) Container
}

// Bucket holds key-values
type Bucket interface {
	Bucket(name string) Bucket
	Cursor() Cursor
	Delete(key []byte) error
	ForEach(fn func(key, value []byte) error) error // deprecated
	Get(key []byte) []byte
	NextSequence() (uint64, error)
	Put(key []byte, value []byte) error
}

// Container operate on records
type Container interface {
	Container(name string) Container
	//Delete(id uint64) error
	//Get(id uint64) (Record, error)
	Iterate(onRecord OnRecord, flags ...bool) error
	//NextID() (uint64, error)
	Put(id uint64, record Record) error
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
