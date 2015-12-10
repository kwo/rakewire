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

// Service implementation of Database
type Service struct {
	sync.Mutex
	db           *bolt.DB
	databaseFile string
	running      bool
}

// NewService creates a new database service.
func NewService(cfg *db.Configuration) *Service {
	return &Service{
		databaseFile: cfg.Location,
	}
}

// Start the database
func (z *Service) Start() error {

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

// Stop the database
func (z *Service) Stop() {

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

// IsRunning indicated if the service is active or not.
func (z *Service) IsRunning() bool {
	z.Lock()
	defer z.Unlock()
	return z.running
}

// Repair the database
func (z *Service) Repair() error {

	// TODO: reimplement repair database

	return nil

}
