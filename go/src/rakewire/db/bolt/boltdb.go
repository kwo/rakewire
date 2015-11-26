package bolt

import (
	"github.com/boltdb/bolt"
	"rakewire/db"
	"rakewire/logging"
	"sync"
	"time"
)

const (
	bucketData           = "Data"
	bucketFeed           = "Feed"
	bucketFeedLog        = "FeedLog"
	bucketIndex          = "Index"
	bucketIndexFeedByURL = "idxFeedByURL"
	bucketIndexNextFetch = "idxNextFetch"
)

// Database implementation of Database
type Database struct {
	sync.Mutex
	db           *bolt.DB
	databaseFile string
}

var (
	logger = logging.New("db")
)

// Open the database
func (z *Database) Open(cfg *db.Configuration) error {

	db, err := bolt.Open(cfg.Location, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		logger.Errorf("Cannot open database at %s. %s", cfg.Location, err.Error())
		return err
	}
	z.db = db
	z.databaseFile = cfg.Location

	// check that buckets exist
	z.Lock()
	err = z.db.Update(func(tx *bolt.Tx) error {

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

		return nil

	})
	z.Unlock()

	if err != nil {
		logger.Errorf("Cannot initialize database: %s", err.Error())
		return err
	}

	logger.Infof("Using database at %s\n", cfg.Location)

	return nil

}

// Close the database
func (z *Database) Close() error {
	if db := z.db; db != nil {
		z.db = nil
		if err := db.Close(); err != nil {
			logger.Warnf("Error closing database: %s\n", err.Error())
			return err
		}
		logger.Info("Closed database")
	}
	return nil
}

// Repair the database
func (z *Database) Repair() error {

	// TODO: reimplement repair database

	return nil

}
