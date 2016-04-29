package model

// Instanz performs high level function upon databases
type Instanz interface {
	Open(location string) (Database, error)
	Close(db Database) error
	// CheckSchema
	// CheckIntegrity
}

// Database defines the interface to a key-value store
type Database interface {
	Location() string
	Select(func(tx Transaction) error) error
	Update(func(tx Transaction) error) error
}

// Transaction represents an atomic operation to the database
type Transaction interface {
	Bucket(name ...string) Bucket
	NextID(name string) (uint64, error)
}

// Bucket holds key-values
type Bucket interface {
	Bucket(name ...string) Bucket
	Cursor() Cursor
	Delete(key []byte) error
	Get(key []byte) []byte
	Put(key, value []byte) error
}

// Cursor loops through values in a bucket
type Cursor interface {
	First() (key, value []byte)
	Last() (key, value []byte)
	Next() (key, value []byte)
	Prev() (key, value []byte)
	Seek(seek []byte) (key, value []byte)
}
