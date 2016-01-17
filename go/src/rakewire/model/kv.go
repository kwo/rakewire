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

const (
	bucketConfig = "Config"
	bucketData   = "Data"
	bucketIndex  = "Index"
)

const (
	logName  = "[bolt]"
	logTrace = "[TRACE]"
	logDebug = "[DEBUG]"
	logInfo  = "[INFO]"
	logWarn  = "[WARN]"
	logError = "[ERROR]"
)

// DataObject defines the functions necessary for objects to be persisted to the database
type DataObject interface {
	GetID() uint64
	SetID(id uint64)
	Clear()
	Serialize(flags ...bool) map[string]string
	Deserialize(map[string]string) error
	IndexKeys() map[string][]string
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

func kvSave(name string, value DataObject, tx Transaction) error {

	b := tx.Bucket(bucketData).Bucket(name)

	if value.GetID() == 0 {
		id, err := b.NextSequence()
		if err != nil {
			return err
		}
		value.SetID(id)
	}

	oldValues, _ := kvGet(value.GetID(), b)
	newValues := value.Serialize()
	newIndexes := value.IndexKeys()

	// save entry
	if err := kvPut(value.GetID(), newValues, b); err != nil {
		return err
	}

	// create old index values
	value.Clear()
	if err := value.Deserialize(oldValues); err != nil {
		return err
	}
	oldIndexes := value.IndexKeys()

	// put new values back
	value.Clear()
	if err := value.Deserialize(newValues); err != nil {
		return err
	}

	// save indexes
	if err := kvSaveIndexes(name, value.GetID(), newIndexes, oldIndexes, tx); err != nil {
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

func kvDelete(name string, value DataObject, tx Transaction) error {

	if value.GetID() == 0 {
		return fmt.Errorf("Cannot delete %s with ID of 0", name)
	}

	// delete data
	b := tx.Bucket(bucketData).Bucket(name)
	c := b.Cursor()
	min := []byte(kvMinKey(value.GetID()))
	nxt := []byte(kvNxtKey(value.GetID()))
	for k, _ := c.Seek(min); k != nil && bytes.Compare(k, nxt) < 0; k, _ = c.Next() {
		if err := c.Delete(); err != nil {
			return err
		}
	}

	// delete indexes
	indexBucket := tx.Bucket(bucketIndex).Bucket(name)
	for indexName, indexElements := range value.IndexKeys() {
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

func kvKey(id uint64) string {
	return fmt.Sprintf("%05d", id)
}

func kvMinKey(id uint64) string {
	return kvKey(id)
}

func kvNxtKey(id uint64) string {
	return kvKey(id + 1)
}

func kvKeys(elements []string) string {
	return strings.Join(elements, chSep)
}
