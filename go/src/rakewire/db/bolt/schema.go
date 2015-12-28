package bolt

import (
	"github.com/boltdb/bolt"
	"log"
	"rakewire/db"
	m "rakewire/model"
	"strconv"
)

const (
	// SchemaVersion of the database
	SchemaVersion = 2
)

func (z *Service) checkDatabase() error {

	// check that buckets exist
	// z.Lock() - called in z.Open
	err := z.db.Update(func(tx *bolt.Tx) error {

		schemaVersion := getSchemaVersion(tx)
		if schemaVersion > 0 && schemaVersion != SchemaVersion {
			if err := z.rebuildIndexes(tx); err != nil {
				return err
			}
		} else {
			if err := z.checkSchema(tx); err != nil {
				return err
			}
		}

		setSchemaVersion(tx, SchemaVersion)

		return nil

	})
	// z.Unlock() - called in z.Open

	return err

}

func (z *Service) checkSchema(tx *bolt.Tx) error {

	log.Printf("%-7s %-7s checking schema...", logDebug, logName)

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

	b, err = bucketData.CreateBucketIfNotExists([]byte(m.GroupEntity))
	if err != nil {
		return err
	}

	b, err = bucketData.CreateBucketIfNotExists([]byte(m.FeedEntity))
	if err != nil {
		return err
	}

	b, err = bucketData.CreateBucketIfNotExists([]byte(m.FeedLogEntity))
	if err != nil {
		return err
	}

	b, err = bucketData.CreateBucketIfNotExists([]byte(m.EntryEntity))
	if err != nil {
		return err
	}

	b, err = bucketData.CreateBucketIfNotExists([]byte(m.UserEntryEntity))
	if err != nil {
		return err
	}

	b, err = bucketData.CreateBucketIfNotExists([]byte(m.UserFeedEntity))
	if err != nil {
		return err
	}

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

	return nil

}

func (z *Service) rebuildIndexes(tx *bolt.Tx) error {

	log.Printf("%-7s %-7s rebuilding indexes", logDebug, logName)

	err := tx.DeleteBucket([]byte(bucketIndex))
	if err != nil {
		return err
	}

	err = z.checkSchema(tx)
	if err != nil {
		return err
	}

	err = z.rebuildIndexesForEntity(m.UserEntity, &m.User{}, tx)
	if err != nil {
		return err
	}

	err = z.rebuildIndexesForEntity(m.GroupEntity, &m.Group{}, tx)
	if err != nil {
		return err
	}

	err = z.rebuildIndexesForEntity(m.FeedEntity, &m.Feed{}, tx)
	if err != nil {
		return err
	}

	err = z.rebuildIndexesForEntity(m.FeedLogEntity, &m.FeedLog{}, tx)
	if err != nil {
		return err
	}

	err = z.rebuildIndexesForEntity(m.EntryEntity, &m.Entry{}, tx)
	if err != nil {
		return err
	}

	err = z.rebuildIndexesForEntity(m.UserEntryEntity, &m.UserEntry{}, tx)
	if err != nil {
		return err
	}

	err = z.rebuildIndexesForEntity(m.UserFeedEntity, &m.UserFeed{}, tx)
	if err != nil {
		return err
	}

	return nil

}

func (z *Service) rebuildIndexesForEntity(entityName string, dao db.DataObject, tx *bolt.Tx) error {

	bEntity := tx.Bucket([]byte(bucketData)).Bucket([]byte(entityName))
	ids, err := kvGetUniqueIDs(bEntity)
	if err != nil {
		return err
	}

	for _, id := range ids {
		dao.Clear()
		if data, ok := kvGet(id, bEntity); ok {
			if err := dao.Deserialize(data); err != nil {
				return err
			}
			if err := kvSaveIndexes(entityName, id, dao.IndexKeys(), nil, tx); err != nil {
				return err
			}
		}
	}

	return nil
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
