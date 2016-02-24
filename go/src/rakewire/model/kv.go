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
	chSep   = "." // must not be a valid key character
	chMax   = "~" // must come after separator character
	fmtTime = "20060102150405Z0700"
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

func kvSave(entityName string, value Object, tx Transaction) error {

	b := tx.Bucket(bucketData).Bucket(entityName)

	fn := func() (uint64, string, error) {
		return kvNextID(entityName, tx)
	}

	if err := value.setIDIfNecessary(fn); err != nil {
		return err
	}

	oldValues := b.GetRecord(value.getID())
	newValues := value.serialize()
	newIndexes := value.serializeIndexes()

	// save item
	if err := b.PutRecord(value.getID(), newValues); err != nil {
		return err
	}

	// create old index values
	oldIndexes := make(map[string]Record)

	if oldValues != nil {

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
			if err := b.Delete(key); err != nil {
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

	keys := []string{}
	b := tx.Bucket(bucketData).Bucket(name)

	// collect keys
	c := b.Cursor()
	min, max := kvKeyMinMax2(value.getID())
	for k, _ := c.Seek(min); k != nil && bytes.Compare(k, max) <= 0; k, _ = c.Next() {
		keys = append(keys, string(k))
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
			if err := b.Delete(key); err != nil {
				return err
			}
		}
	}

	return nil

}

func kvNextID(entityName string, tx Transaction) (uint64, string, error) {

	idStr := "0"
	key := "sequence"
	b := tx.Bucket(bucketConfig)

	// get previous value
	record := b.GetRecord(key)
	if record == nil {
		record = make(Record)
	}
	if value, ok := record[entityName]; ok {
		idStr = string(value)
	}

	// turn into a uint64
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return 0, "", err
	}

	// increment
	id++

	// format as string
	idStr = kvKeyUintEncode(id)

	// save back to database
	record[entityName] = idStr
	if err := b.PutRecord(key, record); err != nil {
		return 0, "", err
	}

	return id, idStr, nil

}

func kvKeyEncode(values ...string) string {
	return strings.Join(values, chSep)
}

func kvKeyEncode2(values ...string) []byte {
	return []byte(kvKeyEncode(values...))
}

func kvKeyDecode(value []byte) []string {
	return strings.Split(string(value), chSep)
}

func kvKeyMax(values ...string) string {
	return strings.Join(values, chSep) + chMax
}

func kvKeyMax2(values ...string) []byte {
	return []byte(kvKeyMax(values...))
}

func kvKeyMinMax(id string) (string, string) {
	return kvKeyEncode(id), kvKeyMax(id)
}

func kvKeyMinMax2(id string) ([]byte, []byte) {
	return []byte(kvKeyEncode(id)), []byte(kvKeyMax(id))
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
