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
		if n >= 10000 {
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
	b, err = bucketData.CreateBucketIfNotExists([]byte(m.UserEntity))
	if err != nil {
		return err
	}
	bumpSequence(b)

	b, err = bucketData.CreateBucketIfNotExists([]byte(m.GroupEntity))
	if err != nil {
		return err
	}
	bumpSequence(b)

	b, err = bucketData.CreateBucketIfNotExists([]byte(m.FeedEntity))
	if err != nil {
		return err
	}
	bumpSequence(b)

	b, err = bucketData.CreateBucketIfNotExists([]byte(m.FeedLogEntity))
	if err != nil {
		return err
	}
	bumpSequence(b)

	b, err = bucketData.CreateBucketIfNotExists([]byte(m.EntryEntity))
	if err != nil {
		return err
	}
	bumpSequence(b)

	b, err = bucketData.CreateBucketIfNotExists([]byte(m.UserEntryEntity))
	if err != nil {
		return err
	}
	bumpSequence(b)

	b, err = bucketData.CreateBucketIfNotExists([]byte(m.UserFeedEntity))
	if err != nil {
		return err
	}
	bumpSequence(b)

	// indexes

	user := m.NewUser("")
	b, err = bucketIndex.CreateBucketIfNotExists([]byte(m.UserEntity))
	if err != nil {
		return err
	}
	for k := range user.IndexKeys() {
		_, err = b.CreateBucketIfNotExists([]byte(k))
		if err != nil {
			return err
		}
	}

	group := m.NewGroup(0, "")
	b, err = bucketIndex.CreateBucketIfNotExists([]byte(m.GroupEntity))
	if err != nil {
		return err
	}
	for k := range group.IndexKeys() {
		_, err = b.CreateBucketIfNotExists([]byte(k))
		if err != nil {
			return err
		}
	}

	feed := m.NewFeed("")
	b, err = bucketIndex.CreateBucketIfNotExists([]byte(m.FeedEntity))
	if err != nil {
		return err
	}
	for k := range feed.IndexKeys() {
		_, err = b.CreateBucketIfNotExists([]byte(k))
		if err != nil {
			return err
		}
	}

	feedlog := m.NewFeedLog(feed.ID)
	b, err = bucketIndex.CreateBucketIfNotExists([]byte(m.FeedLogEntity))
	if err != nil {
		return err
	}
	for k := range feedlog.IndexKeys() {
		_, err = b.CreateBucketIfNotExists([]byte(k))
		if err != nil {
			return err
		}
	}

	entry := m.NewEntry(feed.ID, "")
	b, err = bucketIndex.CreateBucketIfNotExists([]byte(m.EntryEntity))
	if err != nil {
		return err
	}
	for k := range entry.IndexKeys() {
		_, err = b.CreateBucketIfNotExists([]byte(k))
		if err != nil {
			return err
		}
	}

	ue := m.UserEntry{}
	b, err = bucketIndex.CreateBucketIfNotExists([]byte(m.UserEntryEntity))
	if err != nil {
		return err
	}
	for k := range ue.IndexKeys() {
		_, err = b.CreateBucketIfNotExists([]byte(k))
		if err != nil {
			return err
		}
	}

	uf := m.NewUserFeed(user.ID, feed.ID)
	b, err = bucketIndex.CreateBucketIfNotExists([]byte(m.UserFeedEntity))
	if err != nil {
		return err
	}
	for k := range uf.IndexKeys() {
		_, err = b.CreateBucketIfNotExists([]byte(k))
		if err != nil {
			return err
		}
	}

	return setSchemaVersion(tx, 1)

}

// func upgradeTo2(tx *bolt.Tx) error {
//
// 	bucketData, err := tx.CreateBucketIfNotExists([]byte(bucketData))
// 	if err != nil {
// 		return err
// 	}
// 	bucketIndex, err := tx.CreateBucketIfNotExists([]byte(bucketIndex))
// 	if err != nil {
// 		return err
// 	}
//
// 	b, err := bucketData.CreateBucketIfNotExists([]byte(m.GroupEntity))
// 	if err != nil {
// 		return err
// 	}
// 	bumpSequence(b)
//
// 	group := m.NewGroup(0, "")
// 	b, err = bucketIndex.CreateBucketIfNotExists([]byte(m.GroupEntity))
// 	if err != nil {
// 		return err
// 	}
// 	for k := range group.IndexKeys() {
// 		_, err = b.CreateBucketIfNotExists([]byte(k))
// 		if err != nil {
// 			return err
// 		}
// 	}
//
// 	return setSchemaVersion(tx, 2)
//
// }
