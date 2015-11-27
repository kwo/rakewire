package bolt

import (
	"fmt"
	"github.com/boltdb/bolt"
	"strconv"
	"strings"
)

const (
	// SchemaVersion of the database
	SchemaVersion = 2
)

func checkSchema(z *Database) error {

	// check that buckets exist
	z.Lock()
	err := z.db.Update(func(tx *bolt.Tx) error {

		var err error

		for {

			schemaVersion := getSchemaVersion(tx)
			logger.Debugf("Schema Version: %d", schemaVersion)
			if schemaVersion == SchemaVersion {
				break
			}

			switch schemaVersion {
			case 0:
				err = upgradeTo1(tx)
				if err != nil {
					break
				}
			case 1:
				err = upgradeTo2(tx)
				if err != nil {
					break
				}
			default:
				err = fmt.Errorf("Unhandled schema version: %d", schemaVersion)
			}

		} // loop schemaVersion

		return nil

	})
	z.Unlock()

	return err

}

func getSchemaVersion(tx *bolt.Tx) int {

	bucketInfo := tx.Bucket([]byte("Info"))
	if bucketInfo != nil {
		data := bucketInfo.Get([]byte("schema-version"))
		if len(data) > 0 {
			value, err := strconv.ParseInt(string(data), 10, 64)
			if err == nil && value > 0 {
				return int(value)
			}
		}
	}

	return 0

}

func setSchemaVersion(tx *bolt.Tx, version int) error {
	bucketInfo, err := tx.CreateBucketIfNotExists([]byte("Info"))
	if err != nil {
		return err
	}
	return bucketInfo.Put([]byte("schema-version"), []byte(strconv.FormatInt(int64(version), 10)))
}

func upgradeTo1(tx *bolt.Tx) error {

	bucketData, err := tx.CreateBucketIfNotExists([]byte(bucketData))
	if err != nil {
		return err
	}
	bucketIndex, err := tx.CreateBucketIfNotExists([]byte(bucketIndex))
	if err != nil {
		return err
	}

	_, err = bucketData.CreateBucketIfNotExists([]byte(bucketFeed))
	if err != nil {
		return err
	}
	_, err = bucketData.CreateBucketIfNotExists([]byte(bucketFeedLog))
	if err != nil {
		return err
	}

	bucketIndexFeed, err := bucketIndex.CreateBucketIfNotExists([]byte("Feed"))
	if err != nil {
		return err
	}
	if _, err = bucketIndexFeed.CreateBucketIfNotExists([]byte("URL")); err != nil {
		return err
	}
	if _, err = bucketIndexFeed.CreateBucketIfNotExists([]byte("NextFetch")); err != nil {
		return err
	}

	bucketIndexFeedLog, err := bucketIndex.CreateBucketIfNotExists([]byte("FeedLog"))
	if err != nil {
		return err
	}
	_, err = bucketIndexFeedLog.CreateBucketIfNotExists([]byte("FeedTime"))
	if err != nil {
		return err
	}

	return setSchemaVersion(tx, 1)

}

func upgradeTo2(tx *bolt.Tx) error {

	bucketFeed := tx.Bucket([]byte("Data")).Bucket([]byte("Feed"))
	c := bucketFeed.Cursor()

	for k, _ := c.First(); k != nil; k, _ = c.Next() {
		fieldName := strings.SplitN(string(k), chSep, 2)[1]
		if fieldName == "Last" || fieldName == "Last200" {
			c.Delete()
		}
	}

	return setSchemaVersion(tx, 2)

}
