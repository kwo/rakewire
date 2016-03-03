package modelng

import (
	"fmt"
	"strings"
	"time"
)

// Object defines the functions necessary for objects to be persisted to a store
type Object interface {
	getID() string
	setID(Transaction) error
	encode() ([]byte, error)
	decode([]byte) error
	indexes() map[string][]string
}

func delete(entityName string, id string, tx Transaction) error {

	if id != empty {

		bData := tx.Bucket(bucketData, entityName)
		bIndexes := tx.Bucket(bucketIndex, entityName)

		// retrieve data
		if data := bData.Get([]byte(id)); data != nil {

			object := getObject(entityName)

			// decode object
			if err := object.decode(data); err != nil {
				return err
			}

			// delete indexes
			for indexName, indexKeys := range object.indexes() {
				bIndex := bIndexes.Bucket(indexName)
				if err := bIndex.Delete(keyEncode(indexKeys...)); err != nil {
					return err
				}
			}

			// delete object
			if err := bData.Delete(keyEncode(object.getID())); err != nil {
				return err
			}

		} // data not nil

	} // valid id

	return nil

}

func keyDecode(value []byte) []string {
	return strings.Split(string(value), chSep)
}

func keyEncode(values ...string) []byte {
	return []byte(keyEncodeString(values...))
}

func keyEncodeString(values ...string) string {
	return strings.Join(values, chSep)
}

func keyEncodeBool(value bool) string {
	if value {
		return "1"
	}
	return "0"
}

func keyEncodeTime(t time.Time) string {
	return t.Format(fmtTime)
}

func keyEncodeUint(id uint64) string {
	return fmt.Sprintf(fmtUint, id)
}

func save(entityName string, object Object, tx Transaction) error {

	bData := tx.Bucket(bucketData, entityName)
	bIndexes := tx.Bucket(bucketIndex, entityName)

	// delete old indexes, if ID not empty
	if object.getID() != empty {

		// save new data
		newdata, err := object.encode()
		if err != nil {
			return err
		}

		// retrieve old data
		if olddata := bData.Get([]byte(object.getID())); olddata != nil {
			if err := object.decode(olddata); err != nil {
				return err
			}
			for indexName, indexKeys := range object.indexes() {
				bIndex := bIndexes.Bucket(indexName)
				if err := bIndex.Delete(keyEncode(indexKeys...)); err != nil {
					return err
				}
			}
		} // olddata not nil

		// reinstate new data
		if err := object.decode(newdata); err != nil {
			return err
		}

	} // delete old indexes

	// assign new ID, if necessary
	if object.getID() == empty {
		if err := object.setID(tx); err != nil {
			return err
		}
	}

	// save user entity
	if data, err := object.encode(); err == nil {
		if err := bData.Put(keyEncode(object.getID()), data); err != nil {
			return err
		}
	} else {
		return err
	}

	// save new indexes
	for indexName, indexKeys := range object.indexes() {
		bIndex := bIndexes.Bucket(indexName)
		if err := bIndex.Put(keyEncode(indexKeys...), keyEncode(object.getID())); err != nil {
			return err
		}
	}

	return nil

}
