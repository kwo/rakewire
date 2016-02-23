package model

import (
	"bytes"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	empty   = ""
	chSep   = "|" // must come after any valid key character (no closing braces in keys)
	chMax   = "~" // must come after separator character
	fmtTime = time.RFC3339
	fmtUint = "%010d"
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

func kvGet(id string, b Bucket) (map[string]string, bool) {

	found := false
	result := make(map[string]string)

	c := b.Cursor()
	min, max := kvKeyMinMax(id)
	for k, v := c.Seek(min); k != nil && bytes.Compare(k, max) <= 0; k, v = c.Next() {
		// assume proper key format of ID/fieldname
		result[kvKeyDecode(k)[1]] = string(v)
		found = true
	} // for loop

	return result, found

}

func kvGetFromIndex(name, index string, keys []string, tx Transaction) (map[string]string, bool) {

	bIndex := tx.Bucket(bucketIndex).Bucket(name).Bucket(index)
	id := string(bIndex.Get(kvKeyEncode(keys...)))

	if id != empty {
		b := tx.Bucket(bucketData).Bucket(name)
		return kvGet(id, b)
	}

	return nil, false

}

func kvPut(id string, record Record, b Bucket) error {

	keys := [][]byte{}

	// remove record keys not in new set
	c := b.Cursor()
	min, max := kvKeyMinMax(id)
	for k, _ := c.Seek(min); k != nil && bytes.Compare(k, max) <= 0; k, _ = c.Next() {
		if _, ok := record[kvKeyDecode(k)[1]]; !ok {
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
	for fieldName := range record {

		key := kvKeyEncode(id, fieldName)
		value := []byte(record[fieldName])
		if err := b.Put(key, value); err != nil {
			return err
		}
	} // for loop object fields

	return nil

}

func kvSave(entityName string, value Object, tx Transaction) error {

	b := tx.Bucket(bucketData).Bucket(entityName)

	fn := func() (uint64, string, error) {
		return kvNextID(entityName, tx)
	}

	if err := value.setIDIfNecessary(fn); err != nil {
		return err
	}

	oldValues, update := kvGet(value.getID(), b)
	newValues := value.serialize()
	newIndexes := value.serializeIndexes()

	// save item
	if err := kvPut(value.getID(), newValues, b); err != nil {
		return err
	}

	// create old index values
	oldIndexes := make(map[string]Record)

	if update {

		value.clear()
		if err := value.deserialize(oldValues); err != nil {
			return err
		}
		oldIndexes = value.serializeIndexes()

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

func kvSaveIndexes(name string, id string, newIndexes map[string]Record, oldIndexes map[string]Record, tx Transaction) error {

	indexBucket := tx.Bucket(bucketIndex).Bucket(name)

	// delete old indexes
	for indexName := range oldIndexes {
		oldIndexRecord := oldIndexes[indexName]
		b := indexBucket.Bucket(indexName)
		for key := range oldIndexRecord {
			if err := b.Delete([]byte(key)); err != nil {
				return err
			}
		}
	}

	// add new indexes
	for indexName := range newIndexes {
		newIndexRecord := newIndexes[indexName]
		b := indexBucket.Bucket(indexName)
		for key, value := range newIndexRecord {
			if err := b.Put([]byte(key), []byte(value)); err != nil {
				return err
			}
		}
	}

	return nil

}

func kvDelete(name string, value Object, tx Transaction) error {

	if value.getID() == empty {
		return fmt.Errorf("Cannot delete %s with blank ID", name)
	}

	keys := [][]byte{}
	b := tx.Bucket(bucketData).Bucket(name)

	// collect keys
	c := b.Cursor()
	min, max := kvKeyMinMax(value.getID())
	for k, _ := c.Seek(min); k != nil && bytes.Compare(k, max) <= 0; k, _ = c.Next() {
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
	for indexName, record := range value.serializeIndexes() {
		b := indexBucket.Bucket(indexName)
		for key := range record {
			if err := b.Delete([]byte(key)); err != nil {
				return err
			}
		}
	}

	return nil

}

func kvNextID(entityName string, tx Transaction) (uint64, string, error) {

	entryName := "sequence." + strings.ToLower(entityName)

	b := tx.Bucket(bucketConfig)

	// get previous value
	idBytes := b.Get([]byte(entryName))
	if idBytes == nil {
		return 0, "", fmt.Errorf("No sequence found for %s", entityName)
	}

	// turn into a uint64
	id, err := strconv.ParseUint(string(idBytes), 10, 64)
	if err != nil {
		return 0, "", err
	}

	// increment
	id++

	// format as string
	idStr := kvKeyUintEncode(id)

	// save back to database
	if err := b.Put([]byte(entryName), []byte(idStr)); err != nil {
		return 0, "", err
	}

	return id, idStr, nil

}

func kvKeyEncode(values ...string) []byte {
	return []byte(strings.Join(values, chSep))
}

func kvKeyDecode(value []byte) []string {
	return strings.Split(string(value), chSep)
}

func kvKeyMax(values ...string) []byte {
	return []byte(strings.Join(values, chSep) + chMax)
}

func kvKeyMinMax(id string) ([]byte, []byte) {
	return kvKeyEncode(id), kvKeyMax(id)
}

func kvKeyBoolEncode(value bool) string {
	if value {
		return "1"
	}
	return "0"
}

func kvKeyTimeEncode(value time.Time) string {
	return value.UTC().Format(fmtTime)
}

func kvKeyTimeDecode(value string) (time.Time, error) {
	return time.Parse(fmtTime, value)
}

func kvKeyUintEncode(id uint64) string {
	return fmt.Sprintf(fmtUint, id)
}

func kvKeyUintDecode(id string) (uint64, error) {
	return strconv.ParseUint(id, 10, 64)
}

func isStringInArray(a string, b []string) bool {
	for _, x := range b {
		if a == x {
			return true
		}
	}
	return false
}
