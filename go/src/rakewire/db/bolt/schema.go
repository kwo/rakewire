package bolt

import (
	"fmt"
	"github.com/boltdb/bolt"
	"log"
	m "rakewire/model"
	"strconv"
)

const (
	// SchemaVersion of the database
	SchemaVersion = 1
)

func checkSchema(z *Service) error {

	// check that buckets exist
	// z.Lock() - called in z.Open
	err := z.db.Update(func(tx *bolt.Tx) error {

		var err error

		for {

			schemaVersion := getSchemaVersion(tx)
			log.Printf("%-7s %-7s schema version: %d", logDebug, logName, schemaVersion)
			if schemaVersion == SchemaVersion {
				break
			}

			switch schemaVersion {
			case 0:
				err = upgradeTo1(tx)
				if err != nil {
					break
				}
			// case 1:
			// 	err = upgradeTo2(tx)
			// 	if err != nil {
			// 		break
			// 	}
			default:
				err = fmt.Errorf("Unhandled schema version: %d", schemaVersion)
			}

		} // loop schemaVersion

		return nil

	})
	// z.Unlock() - called in z.Open

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

func bumpSequence(b *bolt.Bucket) error {

	for {
		n, err := b.NextSequence()
		if err != nil {
			return err
		}
		if n == 1000 {
			break
		}
	}

	return nil

}

func upgradeTo1(tx *bolt.Tx) error {

	var b *bolt.Bucket

	// top level
	bucketData, err := tx.CreateBucketIfNotExists([]byte(bucketData))
	if err != nil {
		return err
	}
	bucketIndex, err := tx.CreateBucketIfNotExists([]byte(bucketIndex))
	if err != nil {
		return err
	}

	// data
	b, err = bucketData.CreateBucketIfNotExists([]byte(bucketUser))
	if err != nil {
		return err
	}
	bumpSequence(b)

	b, err = bucketData.CreateBucketIfNotExists([]byte(bucketFeed))
	if err != nil {
		return err
	}
	bumpSequence(b)
	b, err = bucketData.CreateBucketIfNotExists([]byte(bucketFeedLog))
	if err != nil {
		return err
	}
	bumpSequence(b)
	b, err = bucketData.CreateBucketIfNotExists([]byte(bucketEntry))
	if err != nil {
		return err
	}
	bumpSequence(b)

	// indexes

	// index user
	bucketIndexUser, err := bucketIndex.CreateBucketIfNotExists([]byte("User"))
	if err != nil {
		return err
	}
	_, err = bucketIndexUser.CreateBucketIfNotExists([]byte("Username"))
	if err != nil {
		return err
	}
	_, err = bucketIndexUser.CreateBucketIfNotExists([]byte("FeverHash"))
	if err != nil {
		return err
	}

	// index feed
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

	// index feedlog
	bucketIndexFeedLog, err := bucketIndex.CreateBucketIfNotExists([]byte("FeedLog"))
	if err != nil {
		return err
	}
	_, err = bucketIndexFeedLog.CreateBucketIfNotExists([]byte("FeedTime"))
	if err != nil {
		return err
	}

	// index entry
	bucketIndexEntry, err := bucketIndex.CreateBucketIfNotExists([]byte(bucketEntry))
	if err != nil {
		return err
	}
	if _, err = bucketIndexEntry.CreateBucketIfNotExists([]byte("Date")); err != nil {
		return err
	}
	if _, err = bucketIndexEntry.CreateBucketIfNotExists([]byte("GUID")); err != nil {
		return err
	}

	u := m.NewUser("testuser@localhost")
	u.SetPassword("abcdefg")

	if err := kvSave(u, tx); err != nil {
		return err
	}

	return setSchemaVersion(tx, 1)

}

// func upgradeTo2(tx *bolt.Tx) error {
//
// 	u := m.NewUser("testuser@localhost")
// 	u.SetPassword("abcdefg")
//
// 	if err := kvSave(u, tx); err != nil {
// 		return err
// 	}
//
// 	return setSchemaVersion(tx, 2)
//
// }
