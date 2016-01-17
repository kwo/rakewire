package middleware

import (
	"github.com/boltdb/bolt"
	"io/ioutil"
	"os"
	"rakewire/model"
	"time"
)

const (
	bucketConfig = "Config"
	bucketData   = "Data"
	bucketIndex  = "Index"
)

func openDatabase() (model.Database, error) {

	f, err := ioutil.TempFile("", "bolt-")
	if err != nil {
		return nil, err
	}
	f.Close()

	db, err := bolt.Open(f.Name(), 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return nil, err
	}

	if err = db.Update(func(tx *bolt.Tx) error {
		return checkSchema(tx)
	}); err != nil {
		return nil, err
	}

	database := &model.BoltDatabase{DB: db}

	err = database.Update(func(tx model.Transaction) error {
		return populateDatabase(tx)
	})

	return database, err

}

func closeDatabase(database model.Database) error {

	boltDB := database.(*model.BoltDatabase)

	if err := boltDB.DB.Close(); err != nil {
		return err
	}

	if err := os.Remove(boltDB.DB.Path()); err != nil {
		return err
	}

	return nil

}

func checkSchema(tx *bolt.Tx) error {

	var b *bolt.Bucket

	// top level
	_, err := tx.CreateBucketIfNotExists([]byte(bucketConfig))
	if err != nil {
		return err
	}
	bucketData, err := tx.CreateBucketIfNotExists([]byte(bucketData))
	if err != nil {
		return err
	}
	bucketIndex, err := tx.CreateBucketIfNotExists([]byte(bucketIndex))
	if err != nil {
		return err
	}

	// data
	b, err = bucketData.CreateBucketIfNotExists([]byte(model.UserEntity))
	if err != nil {
		return err
	}

	b, err = bucketData.CreateBucketIfNotExists([]byte(model.GroupEntity))
	if err != nil {
		return err
	}

	b, err = bucketData.CreateBucketIfNotExists([]byte(model.FeedEntity))
	if err != nil {
		return err
	}

	b, err = bucketData.CreateBucketIfNotExists([]byte(model.FeedLogEntity))
	if err != nil {
		return err
	}

	b, err = bucketData.CreateBucketIfNotExists([]byte(model.EntryEntity))
	if err != nil {
		return err
	}

	b, err = bucketData.CreateBucketIfNotExists([]byte(model.UserEntryEntity))
	if err != nil {
		return err
	}

	b, err = bucketData.CreateBucketIfNotExists([]byte(model.UserFeedEntity))
	if err != nil {
		return err
	}

	// indexes

	user := model.NewUser("")
	b, err = bucketIndex.CreateBucketIfNotExists([]byte(model.UserEntity))
	if err != nil {
		return err
	}
	for k := range user.IndexKeys() {
		_, err = b.CreateBucketIfNotExists([]byte(k))
		if err != nil {
			return err
		}
	}

	group := model.NewGroup(0, "")
	b, err = bucketIndex.CreateBucketIfNotExists([]byte(model.GroupEntity))
	if err != nil {
		return err
	}
	for k := range group.IndexKeys() {
		_, err = b.CreateBucketIfNotExists([]byte(k))
		if err != nil {
			return err
		}
	}

	feed := model.NewFeed("")
	b, err = bucketIndex.CreateBucketIfNotExists([]byte(model.FeedEntity))
	if err != nil {
		return err
	}
	for k := range feed.IndexKeys() {
		_, err = b.CreateBucketIfNotExists([]byte(k))
		if err != nil {
			return err
		}
	}

	feedlog := model.NewFeedLog(feed.ID)
	b, err = bucketIndex.CreateBucketIfNotExists([]byte(model.FeedLogEntity))
	if err != nil {
		return err
	}
	for k := range feedlog.IndexKeys() {
		_, err = b.CreateBucketIfNotExists([]byte(k))
		if err != nil {
			return err
		}
	}

	entry := model.NewEntry(feed.ID, "")
	b, err = bucketIndex.CreateBucketIfNotExists([]byte(model.EntryEntity))
	if err != nil {
		return err
	}
	for k := range entry.IndexKeys() {
		_, err = b.CreateBucketIfNotExists([]byte(k))
		if err != nil {
			return err
		}
	}

	ue := model.UserEntry{}
	b, err = bucketIndex.CreateBucketIfNotExists([]byte(model.UserEntryEntity))
	if err != nil {
		return err
	}
	for k := range ue.IndexKeys() {
		_, err = b.CreateBucketIfNotExists([]byte(k))
		if err != nil {
			return err
		}
	}

	uf := model.NewUserFeed(user.ID, feed.ID)
	b, err = bucketIndex.CreateBucketIfNotExists([]byte(model.UserFeedEntity))
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

func populateDatabase(tx model.Transaction) error {

	// add test user
	user := model.NewUser("karl")
	user.SetPassword("abcdefg")
	return user.Save(tx)

}
