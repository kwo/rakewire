package bolt

import (
	"bytes"
	"fmt"
	"github.com/boltdb/bolt"
	"log"
	"rakewire/db"
	"strconv"
	"strings"
)

const (
	empty = ""
	chMax = "~"
	chMin = " "
	chSep = "|"
)

func kvGet(id uint64, b *bolt.Bucket) (map[string]string, bool) {

	found := false
	result := make(map[string]string)

	c := b.Cursor()
	min := []byte(kvMinKey(id))
	nxt := []byte(kvNxtKey(id))
	for k, v := c.Seek(min); k != nil && bytes.Compare(k, nxt) < 0; k, v = c.Next() {
		// assume proper key format of ID/fieldname
		result[strings.SplitN(string(k), chSep, 2)[1]] = string(v)
		found = true
	} // for loop

	return result, found

}

func kvGetFromIndex(name string, index string, keys []string, tx *bolt.Tx) (map[string]string, bool) {

	bIndex := tx.Bucket([]byte(bucketIndex)).Bucket([]byte(name)).Bucket([]byte(index))

	keyStr := kvKeys(keys)
	idStr := string(bIndex.Get([]byte(keyStr)))

	if idStr != empty {

		id, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			log.Printf("%-7s %-7s error parsing index value: %s", logWarn, logName, idStr)
			return nil, false
		}

		b := tx.Bucket([]byte(bucketData)).Bucket([]byte(name))
		return kvGet(id, b)

	}

	return nil, false

}

func kvPut(id uint64, values map[string]string, b *bolt.Bucket) error {

	// remove old records
	c := b.Cursor()
	min := []byte(kvMinKey(id))
	nxt := []byte(kvNxtKey(id))
	for k, _ := c.Seek(min); k != nil && bytes.Compare(k, nxt) < 0; k, _ = c.Next() {
		c.Delete()
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

func kvSave(value db.DataObject, tx *bolt.Tx) error {

	b := tx.Bucket([]byte(bucketData)).Bucket([]byte(value.GetName()))

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
	if err := kvSaveIndexes(value.GetName(), value.GetID(), newIndexes, oldIndexes, tx); err != nil {
		return err
	}

	return nil

}

func kvSaveIndexes(name string, id uint64, newIndexes map[string][]string, oldIndexes map[string][]string, tx *bolt.Tx) error {

	pkey := strconv.FormatUint(id, 10)
	indexBucket := tx.Bucket([]byte(bucketIndex)).Bucket([]byte(name))

	// delete old indexes
	for indexName := range oldIndexes {
		b := indexBucket.Bucket([]byte(indexName))
		indexElements := oldIndexes[indexName]
		keyStr := kvKeys(indexElements)
		if err := b.Delete([]byte(keyStr)); err != nil {
			return err
		}
	}

	// add new indexes
	for indexName := range newIndexes {
		b := indexBucket.Bucket([]byte(indexName))
		indexElements := newIndexes[indexName]
		keyStr := kvKeys(indexElements)
		if err := b.Put([]byte(keyStr), []byte(pkey)); err != nil {
			return err
		}
	}

	return nil

}

func kvKey(id uint64) string {
	return fmt.Sprintf("%020d", id)
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
