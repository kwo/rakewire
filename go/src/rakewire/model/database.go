package model

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

// OnRecord defines a function type that fires on a new Record
type OnRecord func(Record) error

// OnRecord defines a function type that fires on a new Record
type fnUniqueID func() (string, error)

// Object defines the functions necessary for objects to be persisted to the database
type Object interface {
	getID() string
	setID(fnUniqueID) error
	serialize(...bool) Record
	deserialize(Record, ...bool) error
	serializeIndexes() map[string]Record
}

// Database defines the interface to a key-value store
type Database interface {
	Location() string
	Select(fn func(tx Transaction) error) error
	Update(fn func(tx Transaction) error) error
}

// Transaction represents an atomic operation to the database
type Transaction interface {
	Bucket(name ...string) Bucket
}

// Bucket holds key-values
type Bucket interface {
	Bucket(name ...string) Bucket
	Cursor() Cursor

	Delete(key string) error
	Get(key string) string
	Put(key, value string) error

	DeleteRecord(id string) error
	GetRecord(id string) Record
	PutRecord(id string, record Record) error
	Iterate(onRecord OnRecord) error

	GetIndex(b Bucket, id string) Record
	IterateIndex(b Bucket, minID, maxID string, onRecord OnRecord) error
}

// Cursor loops through values in a bucket
type Cursor interface {
	First() (key []byte, value []byte)
	Last() (key []byte, value []byte)
	Next() (key []byte, value []byte)
	Prev() (key []byte, value []byte)
	Seek(seek []byte) (key []byte, value []byte)
}
