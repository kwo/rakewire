package bolt

import (
	"bytes"
	"github.com/boltdb/bolt"
	"rakewire/db"
	"strconv"
	"strings"
)

const (
	chMin = " "
	chMax = "~"
	chSep = "/"
	empty = ""
)

func kvSave(value db.DataObject, tx *bolt.Tx) error {

	b := tx.Bucket([]byte(bucketData)).Bucket([]byte(value.GetName()))

	if value.GetID() == 0 {
		id, err := b.NextSequence()
		if err != nil {
			return err
		}
		value.SetID(id)
	}

	oldValues := kvGet(value.GetID(), b)
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
	if err := kvIndex(value.GetName(), value.GetID(), newIndexes, oldIndexes, tx); err != nil {
		return err
	}

	return nil

}

func kvGet(id uint64, b *bolt.Bucket) map[string]string {

	pkey := strconv.FormatUint(id, 10)
	result := make(map[string]string)

	c := b.Cursor()
	min := []byte(pkey + chMin)
	max := []byte(pkey + chMax)
	for k, v := c.Seek(min); k != nil && bytes.Compare(k, max) <= 0; k, v = c.Next() {
		// assume proper key format of ID/fieldname
		result[strings.SplitN(string(k), chSep, 2)[1]] = string(v)
	} // for loop

	return result

}

func kvGetIndex(name string, index string, keys []string, tx *bolt.Tx) (map[string]string, error) {

	bIndex := tx.Bucket([]byte(bucketIndex)).Bucket([]byte(name)).Bucket([]byte(index))

	keyStr := strings.Join(keys, chSep)
	idStr := string(bIndex.Get([]byte(keyStr)))

	if idStr != empty {

		pkey, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			return nil, err
		}

		bData := tx.Bucket([]byte(bucketData)).Bucket([]byte(name))
		return kvGet(pkey, bData), nil

	}

	return nil, nil

}

func kvPut(id uint64, values map[string]string, b *bolt.Bucket) error {

	pkey := strconv.FormatUint(id, 10)

	// remove old records
	c := b.Cursor()
	min := []byte(pkey + chMin)
	max := []byte(pkey + chMax)
	for k, _ := c.Seek(min); k != nil && bytes.Compare(k, max) <= 0; k, _ = c.Next() {
		c.Delete()
	}

	// add new records
	for fieldName := range values {
		key := []byte(pkey + chSep + fieldName)
		value := []byte(values[fieldName])
		if err := b.Put(key, value); err != nil {
			return err
		}
	} // for loop object fields

	return nil

}

func kvIndex(name string, id uint64, newIndexes map[string][]string, oldIndexes map[string][]string, tx *bolt.Tx) error {

	pkey := strconv.FormatUint(id, 10)
	indexBucket := tx.Bucket([]byte(bucketIndex)).Bucket([]byte(name))

	// delete old indexes
	for indexName := range oldIndexes {
		b := indexBucket.Bucket([]byte(indexName))
		indexElements := oldIndexes[indexName]
		keyStr := strings.Join(indexElements, chSep)
		if err := b.Delete([]byte(keyStr)); err != nil {
			return err
		}
	}

	// add new indexes
	for indexName := range newIndexes {
		b := indexBucket.Bucket([]byte(indexName))
		indexElements := newIndexes[indexName]
		keyStr := strings.Join(indexElements, chSep)
		if err := b.Put([]byte(keyStr), []byte(pkey)); err != nil {
			return err
		}
	}

	return nil

}
