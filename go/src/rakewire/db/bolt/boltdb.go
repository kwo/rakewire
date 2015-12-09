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
	bucketEntry          = "Entry"
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
	running      bool
}

// NewService creates a new database service.
func NewService(cfg *db.Configuration) *Database {
	return &Database{
		databaseFile: cfg.Location,
	}
}

// Open the database
func (z *Database) Open() error {

	z.Lock()
	defer z.Unlock()
	if z.running {
		log.Printf("%-7s %-7s Database already opened, exiting...", logWarn, logName)
		return nil
	}

	db, err := bolt.Open(z.databaseFile, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		log.Printf("%-7s %-7s Cannot open database at %s. %s", logError, logName, z.databaseFile, err.Error())
		return err
	}
	z.db = db

	if err := checkSchema(z); err != nil {
		log.Printf("%-7s %-7s Cannot initialize database: %s", logError, logName, err.Error())
		return err
	}

	z.running = true
	log.Printf("%-7s %-7s Using database at %s", logInfo, logName, z.databaseFile)
	return nil

}

// Close the database
func (z *Database) Close() {

	z.Lock()
	defer z.Unlock()
	if !z.running {
		log.Printf("%-7s %-7s Database already closed, exiting...", logWarn, logName)
		return
	}

	if err := z.db.Close(); err != nil {
		log.Printf("%-7s %-7s Error closing database: %s", logWarn, logName, err.Error())
		return
	}

	z.db = nil
	z.running = false
	log.Printf("%-7s %-7s Closed database", logInfo, logName)

}

// Repair the database
func (z *Database) Repair() error {

	// TODO: reimplement repair database

	return nil

}
