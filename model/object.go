package model

import (
	"fmt"
	"strings"
	"time"
)

const (
	bucketData  = "Data"
	bucketIndex = "Index"
	bucketTmp   = "tmp"
	chMax       = "~"
	chSep       = "|"
	empty       = ""
	fmtTime     = "20060102150405Z0700"
	fmtUint     = "%010d"
)

var (
	allEntities = map[string][]string{
		entityEntry:        indexesEntry,
		entityFeed:         indexesFeed,
		entityGroup:        indexesGroup,
		entityItem:         indexesItem,
		entitySubscription: indexesSubscription,
		entityTransmission: indexesTransmission,
		entityUser:         indexesUser,
	}
)

// Object defines the functions necessary for objects to be persisted to a store
type Object interface {
	GetID() string
	decode([]byte) error
	encode() ([]byte, error)
	hasIncrementingID() bool
	indexes() map[string][]string
	setID(Transaction) error
}

func getObject(entityName string) Object {
	switch entityName {
	case entityEntry:
		return &Entry{}
	case entityFeed:
		return &Feed{}
	case entityGroup:
		return &Group{}
	case entityItem:
		return &Item{}
	case entitySubscription:
		return &Subscription{}
	case entityTransmission:
		return &Transmission{}
	case entityUser:
		return &User{}
	}
	return nil
}

func deleteObject(tx Transaction, entityName string, id string) error {

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
				if err := bIndex.Delete([]byte(keyEncode(indexKeys...))); err != nil {
					return err
				}
			}

			// delete object
			if err := bData.Delete([]byte(object.GetID())); err != nil {
				return err
			}

		} // data not nil

	} // valid id

	return nil

}

func keyEncode(values ...string) string {
	return strings.Join(values, chSep)
}

func keyEncodeBool(value bool) string {
	if value {
		return "1"
	}
	return "0"
}

func keyEncodeTime(t time.Time) string {
	return t.UTC().Truncate(time.Second).Format(fmtTime)
}

func keyEncodeUint(id uint64) string {
	return fmt.Sprintf(fmtUint, id)
}

func keyMax(key string) []byte {
	return []byte(key + chMax)
}

func keyMinMax(key string) ([]byte, []byte) {
	return []byte(key), keyMax(key)
}

func saveObject(tx Transaction, entityName string, object Object) error {

	bData := tx.Bucket(bucketData, entityName)
	bIndexes := tx.Bucket(bucketIndex, entityName)

	// delete old indexes, if ID not empty
	if object.GetID() != empty {

		// save new data
		newdata, err := object.encode()
		if err != nil {
			return err
		}

		// retrieve old data
		if olddata := bData.Get([]byte(object.GetID())); olddata != nil {
			if err := object.decode(olddata); err != nil {
				return err
			}
			for indexName, indexKeys := range object.indexes() {
				bIndex := bIndexes.Bucket(indexName)
				if err := bIndex.Delete([]byte(keyEncode(indexKeys...))); err != nil {
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
	if object.GetID() == empty {
		if err := object.setID(tx); err != nil {
			return err
		}
	}

	// save entity
	if data, err := object.encode(); err == nil {
		if err := bData.Put([]byte(object.GetID()), data); err != nil {
			return err
		}
	} else {
		return err
	}

	// save new indexes
	for indexName, indexKeys := range object.indexes() {
		bIndex := bIndexes.Bucket(indexName)
		if err := bIndex.Put([]byte(keyEncode(indexKeys...)), []byte(object.GetID())); err != nil {
			return err
		}
	}

	return nil

}
