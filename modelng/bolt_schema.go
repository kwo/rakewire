package modelng

import (
	"github.com/boltdb/bolt"
)

func checkSchema(tx *bolt.Tx) error {

	var b *bolt.Bucket

	// top level
	bData, err := tx.CreateBucketIfNotExists([]byte(bucketData))
	if err != nil {
		return err
	}
	bIndex, err := tx.CreateBucketIfNotExists([]byte(bucketIndex))
	if err != nil {
		return err
	}

	// data & indexes
	for entityName, entityIndexes := range allEntities {
		if _, err = bData.CreateBucketIfNotExists([]byte(entityName)); err != nil {
			return err
		}
		if b, err = bIndex.CreateBucketIfNotExists([]byte(entityName)); err == nil {
			for _, indexName := range entityIndexes {
				if _, err = b.CreateBucketIfNotExists([]byte(indexName)); err != nil {
					return err
				}
			} // entityIndexes
		} else {
			return err
		}
	} // allEntities

	return nil

}
