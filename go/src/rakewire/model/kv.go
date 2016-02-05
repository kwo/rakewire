package model

import (
	"bytes"
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"
)

const (
	empty = ""
	chSep = "|"
)

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
	sort.Strings(z.MissingFieldnames)
	sort.Strings(z.UnknownFieldnames)
	errors := []string{}
	for _, err := range z.Errors {
		errors = append(errors, err.Error())
	}
	message := fmt.Sprintf("Invalid %s: missing %s; unknown: %s; errors: %s", z.Entity, strings.Join(z.MissingFieldnames, ", "), strings.Join(z.UnknownFieldnames, ", "), strings.Join(errors, ", "))
	return message
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

func kvPut(id uint64, values map[string]string, b Bucket) error {

	keys := [][]byte{}

	// remove record keys not in new set
	c := b.Cursor()
	min := []byte(kvMinKey(id))
	nxt := []byte(kvNxtKey(id))
	for k, _ := c.Seek(min); k != nil && bytes.Compare(k, nxt) < 0; k, _ = c.Next() {
		if _, ok := values[kvKeyElement(k, 1)]; !ok {
			keys = append(keys, k)
		}
	}

	// delete keys
	for _, k := range keys {
		if err := b.Delete(k); err != nil {
			return err
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

func kvSave(entityName string, value Object, tx Transaction) error {

	b := tx.Bucket(bucketData).Bucket(entityName)

	if value.getID() == 0 {
		id, err := kvNextID(entityName, tx)
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
	if err := kvSaveIndexes(entityName, value.getID(), newIndexes, oldIndexes, tx); err != nil {
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

	keys := [][]byte{}

	// delete data
	b := tx.Bucket(bucketData).Bucket(name)
	c := b.Cursor()
	min := []byte(kvMinKey(value.getID()))
	nxt := []byte(kvNxtKey(value.getID()))
	for k, _ := c.Seek(min); k != nil && bytes.Compare(k, nxt) < 0; k, _ = c.Next() {
		keys = append(keys, k)
	}

	// delete keys
	for _, k := range keys {
		if err := b.Delete(k); err != nil {
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

func kvNextID(entityName string, tx Transaction) (uint64, error) {

	entryName := "sequence." + strings.ToLower(entityName)

	b := tx.Bucket(bucketConfig)

	// get previous value
	idBytes := b.Get([]byte(entryName))
	if idBytes == nil {
		return 0, fmt.Errorf("No sequence found for %s", entityName)
	}

	// turn into a uint64
	id, err := strconv.ParseUint(string(idBytes), 10, 64)
	if err != nil {
		return 0, err
	}

	// increment
	id++

	// format as string
	idStr := strconv.FormatUint(id, 10)

	// save back to database
	if err := b.Put([]byte(entryName), []byte(idStr)); err != nil {
		return 0, err
	}

	return id, nil

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
