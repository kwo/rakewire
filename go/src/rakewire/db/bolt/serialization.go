package bolt

import (
	"bytes"
	"fmt"
	"github.com/boltdb/bolt"
	"rakewire/serial"
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

// TODO: need a control table to keep track of schema version,
// need functions to convert from one schema version to the next,
// and that is how to rename fields.

// Get retrieves the object with the given ID from the specified bucket.
func Get(object interface{}, tx *bolt.Tx) error {

	_, data, err := serial.Encode(object)
	if err != nil {
		return err
	}

	values, err := load(data.Name, data.Key, tx)
	if err != nil {
		return err
	}

	if err = serial.Decode(object, values); err != nil {
		return err
	}

	return nil

}

// Put saves the object with the given ID to the specified bucket.
func Put(object interface{}, tx *bolt.Tx) error {

	meta, data, err := serial.Encode(object)
	if err != nil {
		return err
	}

	valuesOld, err := load(data.Name, data.Key, tx)
	if err != nil {
		return err
	}
	dataOld := serial.DataFrom(meta, valuesOld)

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

	minS, err := serial.EncodeFields(min...)
	if err != nil {
		return err
	}
	maxS, err := serial.EncodeFields(max...)
	if err != nil {
		return err
	}

	minB := []byte(strings.Join(minS, chSep) + chSep + chMin)
	maxB := []byte(strings.Join(maxS, chSep) + chSep + chMax)

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

func load(name string, pkey string, tx *bolt.Tx) (map[string]string, error) {

	logger.Debugf("Loading %s ...", name)
	bucketData := tx.Bucket([]byte(bucketData))
	b := bucketData.Bucket([]byte(name))

	result := make(map[string]string)

	// seek to the min key for the given ID
	// loop through cursor until max ID
	c := b.Cursor()
	min := []byte(pkey + chMin)
	max := []byte(pkey + chMax)
	for k, v := c.Seek(min); k != nil && bytes.Compare(k, max) <= 0; k, v = c.Next() {

		// extract field name from key
		// key format: ID/fieldname
		key := string(k)
		keyParts := strings.SplitN(key, chSep, 2)
		if len(keyParts) != 2 {
			return nil, fmt.Errorf("Malformatted key: %s.", key)
		}
		fieldName := keyParts[1]

		result[fieldName] = string(v)

	} // for loop

	return result, nil

}

func rangeBucket(name string, min []byte, max []byte, add func() interface{}, tx *bolt.Tx) error {

	b := tx.Bucket([]byte(bucketData)).Bucket([]byte(name))
	c := b.Cursor()

	var pkey string
	values := make(map[string]string)

	for k, v := c.Seek(min); k != nil && bytes.Compare(k, max) <= 0; k, v = c.Next() {

		// extract field name from key
		// key format: ID/fieldname
		key := string(k)
		keyParts := strings.SplitN(key, chSep, 2)
		if len(keyParts) != 2 {
			return fmt.Errorf("Malformatted key: %s.", key)
		}
		fieldName := keyParts[1]

		values[fieldName] = string(v)

		akey := keyParts[0]
		if akey != pkey {
			pkey = akey
			err := serial.Decode(add(), values)
			if err != nil {
				return err
			}
			values = make(map[string]string)
		} // new key

	} // cursor

	return nil

}

func rangeIndex(name string, index string, min []byte, max []byte, add func() interface{}, tx *bolt.Tx) error {

	b := tx.Bucket([]byte(bucketIndex)).Bucket([]byte(name)).Bucket([]byte(index))
	c := b.Cursor()

	for k, v := c.Seek(min); k != nil && bytes.Compare(k, max) <= 0; k, v = c.Next() {

		pkey := string(v)

		values, err := load(name, pkey, tx)
		if err != nil {
			return err
		}
		err = serial.Decode(add(), values)
		if err != nil {
			return err
		}

	} // cursor

	return nil

}

func save(name string, pkey string, values map[string]string, tx *bolt.Tx) error {

	keyList := make(map[string]bool)

	logger.Debugf("Saving %s ...", name)
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
