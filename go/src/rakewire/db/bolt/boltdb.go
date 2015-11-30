package bolt

import (
	"github.com/boltdb/bolt"
	"log"
	"rakewire/db"
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

const (
	logName  = "[bolt]"
	logDebug = "[DEBUG]"
	logInfo  = "[INFO]"
	logWarn  = "[WARN]"
	logError = "[ERROR]"
)

// Database implementation of Database
type Database struct {
	sync.Mutex
	db           *bolt.DB
	databaseFile string
}

// Open the database
func (z *Database) Open(cfg *db.Configuration) error {

	db, err := bolt.Open(cfg.Location, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		log.Printf("%-7s %-7s Cannot open database at %s. %s", logError, logName, cfg.Location, err.Error())
		return err
	}
	z.db = db
	z.databaseFile = cfg.Location

	if err := checkSchema(z); err != nil {
		log.Printf("%-7s %-7s Cannot initialize database: %s", logError, logName, err.Error())
		return err
	}

	log.Printf("%-7s %-7s Using database at %s", logInfo, logName, cfg.Location)

	return nil

}

// Close the database
func (z *Database) Close() error {
	if db := z.db; db != nil {
		z.db = nil
		if err := db.Close(); err != nil {
			log.Printf("%-7s %-7s Error closing database: %s", logWarn, logName, err.Error())
			return err
		}
		log.Printf("%-7s %-7s Closed database", logInfo, logName)
	}
	return nil
}

// Repair the database
func (z *Database) Repair() error {

	// TODO: reimplement repair database

	return nil

}
