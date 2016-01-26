package model

import (
	"bytes"
	"fmt"
	"log"
	"strconv"
	"strings"
)

const (
	empty = ""
	chSep = "|"
)

// Record defines a group of key-value pairs that can create a new Object
type Record map[string]string

// GetID return the primary key of the object.
func (z Record) GetID() uint64 {
	id, _ := strconv.ParseUint(z["ID"], 10, 64)
	return id
}

// SetID sets the primary key of the object.
func (z Record) SetID(id uint64) {
	z["ID"] = strconv.FormatUint(id, 10)
}

// OnRecord defines a function type that fires on a new Record
type OnRecord func(Record) error

// Object defines the functions necessary for objects to be persisted to the database
type Object interface {
	getID() uint64
	setID(id uint64)
	clear()
	serialize(...bool) Record
	deserialize(Record, ...bool) error
	indexKeys() map[string][]string
}

// NewDeserializationError returns a new DeserializationError or nil of all arrays are empty.
func newDeserializationError(entityName string, errors []error, missing []string, unknown []string) error {

	if len(errors) > 0 || len(missing) > 0 || len(unknown) > 0 {
		return &DeserializationError{
			Entity:            entityName,
			Errors:            errors,
			MissingFieldnames: missing,
			UnknownFieldnames: unknown,
		}
	}

	return nil

}

// DeserializationError represents multiple errors encountered during deserialization
type DeserializationError struct {
	Entity            string
	Errors            []error
	MissingFieldnames []string
	UnknownFieldnames []string
}

func (z DeserializationError) Error() string {
	texts := []string{}
	for _, err := range z.Errors {
		texts = append(texts, err.Error())
	}
	for _, field := range z.MissingFieldnames {
		texts = append(texts, fmt.Sprintf("Missing field in %s: %s", z.Entity, field))
	}
	for _, field := range z.UnknownFieldnames {
		texts = append(texts, fmt.Sprintf("Unknown field in %s: %s", z.Entity, field))
	}
	return strings.Join(texts, "\n")
}

func kvGet(id uint64, b Bucket) (map[string]string, bool) {

	found := false
	result := make(map[string]string)

	c := b.Cursor()
	min := []byte(kvMinKey(id))
	nxt := []byte(kvNxtKey(id))
	for k, v := c.Seek(min); k != nil && bytes.Compare(k, nxt) < 0; k, v = c.Next() {
		// assume proper key format of ID/fieldname
		result[kvKeyElement(k, 1)] = string(v)
		found = true
	} // for loop

	return result, found

}

func kvGetFromIndex(name string, index string, keys []string, tx Transaction) (map[string]string, bool) {

	bIndex := tx.Bucket(bucketIndex).Bucket(name).Bucket(index)

	keyStr := kvKeys(keys)
	idStr := string(bIndex.Get([]byte(keyStr)))

	if idStr != empty {

		id, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			log.Printf("%-7s %-7s error parsing index value: %s", logWarn, logName, idStr)
			return nil, false
		}

		b := tx.Bucket(bucketData).Bucket(name)
		return kvGet(id, b)

	}

	return nil, false

}

// kvGetUniqueIDs: deprecated
func kvGetUniqueIDs(b Bucket) ([]uint64, error) {

	var result []uint64

	var lastID uint64
	err := b.ForEach(func(k, v []byte) error {
		id, err := kvKeyElementID(k, 0)
		if err != nil {
			return err
		}
		if id != lastID {
			result = append(result, id)
			lastID = id
		}
		return nil
	})

	return result, err

}

func kvPut(id uint64, values map[string]string, b Bucket) error {

	// remove record keys not in new set
	c := b.Cursor()
	min := []byte(kvMinKey(id))
	nxt := []byte(kvNxtKey(id))
	for k, _ := c.Seek(min); k != nil && bytes.Compare(k, nxt) < 0; k, _ = c.Next() {
		if _, ok := values[kvKeyElement(k, 1)]; !ok {
			if err := c.Delete(); err != nil {
				return err
			}
		}
	}

	// add new records
	for fieldName := range values {
		key := []byte(kvKey(id) + chSep + fieldName)
		value := []byte(values[fieldName])
		if err := b.Put(key, value); err != nil {
			return err
		}
	} // for loop object fields

	return nil

}

func kvSave(name string, value Object, tx Transaction) error {

	b := tx.Bucket(bucketData).Bucket(name)

	if value.getID() == 0 {
		id, err := b.NextSequence()
		if err != nil {
			return err
		}
		value.setID(id)
	}

	oldValues, update := kvGet(value.getID(), b)
	newValues := value.serialize()
	newIndexes := value.indexKeys()

	// save item
	if err := kvPut(value.getID(), newValues, b); err != nil {
		return err
	}

	// create old index values
	oldIndexes := make(map[string][]string)

	if update {

		value.clear()
		if err := value.deserialize(oldValues); err != nil {
			return err
		}
		oldIndexes = value.indexKeys()

		// put new values back
		value.clear()
		if err := value.deserialize(newValues); err != nil {
			return err
		}

	}

	// save indexes
	if err := kvSaveIndexes(name, value.getID(), newIndexes, oldIndexes, tx); err != nil {
		return err
	}

	return nil

}

func kvSaveIndexes(name string, id uint64, newIndexes map[string][]string, oldIndexes map[string][]string, tx Transaction) error {

	pkey := strconv.FormatUint(id, 10)
	indexBucket := tx.Bucket(bucketIndex).Bucket(name)

	// delete outdated indexes
	for indexName := range oldIndexes {
		oldIndexElements := oldIndexes[indexName]
		newIndexElements := newIndexes[indexName]
		oldKeyStr := kvKeys(oldIndexElements)
		newKeyStr := kvKeys(newIndexElements)
		if oldKeyStr != newKeyStr {
			b := indexBucket.Bucket(indexName)
			if err := b.Delete([]byte(oldKeyStr)); err != nil {
				return err
			}
		}
	}

	// add updated indexes
	for indexName := range newIndexes {
		newIndexElements := newIndexes[indexName]
		oldIndexElements := oldIndexes[indexName]
		newKeyStr := kvKeys(newIndexElements)
		oldKeyStr := kvKeys(oldIndexElements)
		if newKeyStr != oldKeyStr {
			b := indexBucket.Bucket(indexName)
			if err := b.Put([]byte(newKeyStr), []byte(pkey)); err != nil {
				return err
			}
		}
	}

	return nil

}

func kvDelete(name string, value Object, tx Transaction) error {

	if value.getID() == 0 {
		return fmt.Errorf("Cannot delete %s with ID of 0", name)
	}

	// delete data
	b := tx.Bucket(bucketData).Bucket(name)
	c := b.Cursor()
	min := []byte(kvMinKey(value.getID()))
	nxt := []byte(kvNxtKey(value.getID()))
	for k, _ := c.Seek(min); k != nil && bytes.Compare(k, nxt) < 0; k, _ = c.Next() {
		if err := c.Delete(); err != nil {
			return err
		}
	}

	// delete indexes
	indexBucket := tx.Bucket(bucketIndex).Bucket(name)
	for indexName, indexElements := range value.indexKeys() {
		b := indexBucket.Bucket(indexName)
		keyStr := kvKeys(indexElements)
		if err := b.Delete([]byte(keyStr)); err != nil {
			return err
		}
	}

	return nil

}

func kvKeyElement(k []byte, index int) string {
	return strings.Split(string(k), chSep)[index]
}

func kvKeyElementID(k []byte, index int) (uint64, error) {
	return strconv.ParseUint(kvKeyElement(k, index), 10, 64)
}

func kvMinKey(id uint64) string {
	return kvKey(id)
}

func kvNxtKey(id uint64) string {
	return kvKey(id + 1)
}

func kvKey(id uint64) string {
	return fmt.Sprintf("%05d", id)
}

func kvKeys(elements []string) string {
	return strings.Join(elements, chSep)
}

func kvBucketKeyEncode(id uint64, fieldname string) []byte {
	return []byte(kvKeys([]string{kvKey(id), fieldname}))
}

func kvBucketKeyDecode(key []byte) (uint64, string, error) {
	fields := strings.Split(string(key), chSep)
	id, err := strconv.ParseUint(fields[0], 10, 64)
	if len(fields) != 2 {
		err = fmt.Errorf("Invalid key, must have two fields: %s", key)
	}
	return id, fields[1], err
}

func kvIndexKeyEncode(fields ...string) []byte {
	return []byte(kvKeys(fields))
}

func kvIndexKeyDecode(key []byte) []string {
	return strings.Split(string(key), chSep)
}

func isStringInArray(a string, b []string) bool {
	for _, x := range b {
		if a == x {
			return true
		}
	}
	return false
}
