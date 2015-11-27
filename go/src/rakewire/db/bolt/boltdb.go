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

	if err := checkSchema(z); err != nil {
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
