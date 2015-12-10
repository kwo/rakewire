package bolt

import (
	"errors"
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

var (
	// ErrRestart indicates that the service cannot be started because it is already running.
	ErrRestart = errors.New("The service is already started")
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
		log.Printf("%-7s %-7s service already started, exiting...", logWarn, logName)
		return ErrRestart
	}

	db, err := bolt.Open(z.databaseFile, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		log.Printf("%-7s %-7s cannot open database at %s. %s", logError, logName, z.databaseFile, err.Error())
		return err
	}
	z.db = db

	if err := checkSchema(z); err != nil {
		log.Printf("%-7s %-7s cannot initialize database: %s", logError, logName, err.Error())
		return err
	}

	z.running = true
	log.Printf("%-7s %-7s service started using %s", logInfo, logName, z.databaseFile)
	return nil

}

// Stop the database
func (z *Service) Stop() {

	z.Lock()
	defer z.Unlock()
	if !z.running {
		log.Printf("%-7s %-7s service already stopped, exiting...", logWarn, logName)
		return
	}

	if err := z.db.Close(); err != nil {
		log.Printf("%-7s %-7s error closing database: %s", logWarn, logName, err.Error())
		return
	}

	z.db = nil
	z.running = false
	log.Printf("%-7s %-7s service stopped", logInfo, logName)

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
