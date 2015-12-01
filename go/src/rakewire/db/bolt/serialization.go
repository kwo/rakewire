package bolt

import (
	"bytes"
	"fmt"
	"github.com/boltdb/bolt"
	"rakewire/kv"
	"strings"
)

const (
	chMin = " "
	chMax = "~"
	chSep = "/"
	empty = ""
)

// TODO: need delete function similar to query
// TODO: add function to register data types on start to create buckets
// TODO: create function to rebuild indexes

// Get retrieves the object with the given ID from the specified bucket.
func Get(object interface{}, tx *bolt.Tx) error {

	meta, data, err := kv.Encode(object)
	if err != nil {
		return err
	}

	values := load(data.Name, data.Key, tx)
	if len(values) == 0 {
		// if the record is not found, unset the id on the object
		values[meta.Key] = empty
	}
	if err := kv.Decode(object, values); err != nil {
		return err
	}

	return nil

}

// Put saves the object with the given ID to the specified bucket.
func Put(object interface{}, tx *bolt.Tx) error {

	meta, data, err := kv.Encode(object)
	if err != nil {
		return err
	}

	valuesOld := load(data.Name, data.Key, tx)
	dataOld := kv.DataFrom(meta, valuesOld)

	if err = save(data.Name, data.Key, data.Values, tx); err != nil {
		return err
	}

	if err = index(data.Name, data.Key, data.Indexes, dataOld.Indexes, tx); err != nil {
		return err
	}

	return nil

}

// Query retrieves objects using the given criteria.
func Query(name string, index string, min []interface{}, max []interface{}, add func() interface{}, tx *bolt.Tx) error {

	marker := func(v []interface{}) ([]byte, error) {
		if v == nil {
			return nil, nil
		}
		fields, err := kv.EncodeFields(v...)
		if err != nil {
			return nil, err
		}
		return []byte(strings.Join(fields, chSep)), nil
	}

	minB, err := marker(min)
	if err != nil {
		return err
	}
	maxB, err := marker(max)
	if err != nil {
		return err
	}

	if index != empty {
		return rangeIndex(name, index, minB, maxB, add, tx)
	}

	return rangeBucket(name, minB, maxB, add, tx)

}

func index(name string, pkey string, indexesNew map[string][]string, indexesOld map[string][]string, tx *bolt.Tx) error {

	indexBucket := tx.Bucket([]byte(bucketIndex)).Bucket([]byte(name))

	// delete old indexes
	for indexName := range indexesOld {
		b := indexBucket.Bucket([]byte(indexName))
		indexElements := indexesOld[indexName]
		keyStr := strings.Join(indexElements, chSep)
		if err := b.Delete([]byte(keyStr)); err != nil {
			return err
		}
	}

	// add new indexes
	for indexName := range indexesNew {
		b := indexBucket.Bucket([]byte(indexName))
		indexElements := indexesNew[indexName]
		keyStr := strings.Join(indexElements, chSep)
		if err := b.Put([]byte(keyStr), []byte(pkey)); err != nil {
			return err
		}
	}

	return nil

}

func load(name string, pkey string, tx *bolt.Tx) map[string]string {

	//logger.Debugf("Loading %s ...", name)
	bucketData := tx.Bucket([]byte(bucketData))
	b := bucketData.Bucket([]byte(name))

	result := make(map[string]string)

	// seek to the min key for the given ID
	// loop through cursor until max ID
	c := b.Cursor()
	min := []byte(pkey + chMin)
	max := []byte(pkey + chMax)
	for k, v := c.Seek(min); k != nil && bytes.Compare(k, max) <= 0; k, v = c.Next() {
		// assume proper key format of ID/fieldname
		result[strings.SplitN(string(k), chSep, 2)[1]] = string(v)
	} // for loop

	return result

}

func rangeBucket(name string, min []byte, max []byte, add func() interface{}, tx *bolt.Tx) error {

	b := tx.Bucket([]byte(bucketData)).Bucket([]byte(name))
	c := b.Cursor()

	initCursor := func(v []byte) ([]byte, []byte) {
		if v == nil || len(v) == 0 {
			return c.First()
		}
		return c.Seek(v)
	}

	lastKey := empty

	for k, _ := initCursor(min); k != nil && bytes.Compare(k, max) <= 0; k, _ = c.Next() {

		// assume proper key format of ID/fieldname
		pkey := strings.SplitN(string(k), chSep, 2)[0]

		if pkey != lastKey {
			lastKey = pkey
			values := load(name, pkey, tx)
			if err := kv.Decode(add(), values); err != nil {
				return err
			}
		}

	} // cursor

	return nil

}

func rangeIndex(name string, index string, min []byte, max []byte, add func() interface{}, tx *bolt.Tx) error {

	b := tx.Bucket([]byte(bucketIndex)).Bucket([]byte(name)).Bucket([]byte(index))
	c := b.Cursor()

	initCursor := func(v []byte) ([]byte, []byte) {
		if v == nil || len(v) == 0 {
			return c.First()
		}
		return c.Seek(v)
	}

	for k, v := initCursor(min); k != nil && bytes.Compare(k, max) <= 0; k, v = c.Next() {

		pkey := string(v)

		values := load(name, pkey, tx)
		if err := kv.Decode(add(), values); err != nil {
			return err
		}

	} // cursor

	return nil

}

func save(name string, pkey string, values map[string]string, tx *bolt.Tx) error {

	keyList := make(map[string]bool)

	//logger.Debugf("Saving %s ...", name)
	bucketData := tx.Bucket([]byte(bucketData))
	b := bucketData.Bucket([]byte(name))

	// loop thru values
	for fieldName := range values {

		keyStr := fmt.Sprintf("%s%s%s", pkey, chSep, fieldName)
		valueStr := values[fieldName]

		key := []byte(keyStr)
		value := []byte(valueStr)
		if len(value) > 0 {
			keyList[keyStr] = true
			if err := b.Put(key, value); err != nil {
				return err
			}
		}

	} // for loop object fields

	// loop with cursor thru record, remove database fields not in field list
	c := b.Cursor()
	min := []byte(pkey + chMin)
	max := []byte(pkey + chMax)
	// logger.Debugf("marshall cursor min/max: %s / %s", min, max)
	for k, _ := c.Seek(min); k != nil && bytes.Compare(k, max) <= 0; k, _ = c.Next() {
		if !keyList[string(k)] {
			// logger.Debugf("Removing key: %s", k)
			c.Delete()
		}
	}

	return nil

}
