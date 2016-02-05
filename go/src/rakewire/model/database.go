package model

import (
	"strconv"
)

var (
	allEntities = map[string][]string{
		entryEntity:        entryAllIndexes,
		feedEntity:         feedAllIndexes,
		groupEntity:        groupAllIndexes,
		itemEntity:         itemAllIndexes,
		subscriptionEntity: subscriptionAllIndexes,
		transmissionEntity: transmissionAllIndexes,
		userEntity:         userAllIndexes,
	}
)

// Record defines a group of key-value pairs that can create a new Object
type Record map[string]string

// GetID return the primary key of the object.
func (z Record) GetID() uint64 {
	id, _ := strconv.ParseUint(z["ID"], 10, 64)
	return id
}

// OnRecord defines a function type that fires on a new Record
type OnRecord func(Record) error

// OnRecord defines a function type that fires on a new Record
type fnNextID func(string, Transaction) (uint64, string, error)

// Object defines the functions necessary for objects to be persisted to the database
type Object interface {
	getID() uint64
	setIDIfNecessary(fnNextID, Transaction) error
	clear()
	serialize(...bool) Record
	deserialize(Record, ...bool) error
	indexKeys() map[string][]string
}

// ContainerSeparator specified the separator character for container names
const ContainerSeparator = "/"

// Database defines the interface to a key-value store
type Database interface {
	Location() string
	Select(fn func(tx Transaction) error) error
	Update(fn func(tx Transaction) error) error
}

// Transaction represents an atomic operation to the database
type Transaction interface {
	Bucket(name string) Bucket
	Container(paths ...string) (Container, error)
}

// Bucket holds key-values
type Bucket interface {
	Bucket(name string) Bucket
	Cursor() Cursor
	Delete(key []byte) error
	Get(key []byte) []byte
	Put(key []byte, value []byte) error
}

// Container operate on records
type Container interface {
	Container(paths ...string) (Container, error)
	Delete(id uint64) error
	Get(id uint64) (Record, error)
	Iterate(onRecord OnRecord, flags ...bool) error
	Put(record Record) error
}

// Cursor loops through values in a bucket
type Cursor interface {
	First() (key []byte, value []byte)
	Last() (key []byte, value []byte)
	Next() (key []byte, value []byte)
	Prev() (key []byte, value []byte)
	Seek(seek []byte) (key []byte, value []byte)
}
