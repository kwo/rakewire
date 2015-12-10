package httpd

import (
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"rakewire/db"
	"sync"
)

const (
	logName  = "[httpd]"
	logDebug = "[DEBUG]"
	logInfo  = "[INFO]"
	logWarn  = "[WARN]"
	logError = "[ERROR]"
)

var (
	// ErrRestart indicates that the service cannot be started because it is already running.
	ErrRestart = errors.New("The service is already started")
)

// Service server
type Service struct {
	sync.Mutex
	cfg      *Configuration
	database db.Database
	listener net.Listener
	running  bool
}

// Configuration configuration
type Configuration struct {
	AccessLog string
	Address   string
	Port      int
	UseLocal  bool
}

const (
	hAcceptEncoding  = "Accept-Encoding"
	hContentEncoding = "Content-Encoding"
	hContentLength   = "Content-Length"
	hContentType     = "Content-Type"
	hVary            = "Vary"
	mGet             = "GET"
	mPost            = "POST"
	mPut             = "PUT"
	mimeHTML         = "text/html; charset=utf-8"
	mimeJSON         = "application/json"
	mimeText         = "text/plain; charset=utf-8"
)

// NewService creates a new httpd service.
func NewService(cfg *Configuration, database db.Database) *Service {
	return &Service{
		cfg:      cfg,
		database: database,
	}
}

// Start web service
func (z *Service) Start() error {

	z.Lock()
	defer z.Unlock()
	if z.running {
		log.Printf("%-7s %-7s service already started, exiting...", logWarn, logName)
		return ErrRestart
	}

	if z.database == nil {
		log.Printf("%-7s %-7s cannot start httpd, no database provided", logError, logName)
		return errors.New("No database")
	}

	router, err := z.mainRouter(z.cfg.UseLocal)
	if err != nil {
		log.Printf("%-7s %-7s cannot load router: %s", logError, logName, err.Error())
		return err
	}
	mainHandler := Adapt(router, LogAdapter(z.cfg.AccessLog))

	z.listener, err = net.Listen("tcp", fmt.Sprintf("%s:%d", z.cfg.Address, z.cfg.Port))
	if err != nil {
		log.Printf("%-7s %-7s cannot start listener: %s", logError, logName, err.Error())
		return err
	}

	server := http.Server{
		Handler: mainHandler,
	}

	go server.Serve(z.listener)

	z.running = true
	log.Printf("%-7s %-7s service started on http://%s:%d", logInfo, logName, z.cfg.Address, z.cfg.Port)
	return nil

}

// Stop stop the server
func (z *Service) Stop() {

	z.Lock()
	defer z.Unlock()
	if !z.running {
		log.Printf("%-7s %-7s service already stopped, exiting...", logWarn, logName)
		return
	}

	if err := z.listener.Close(); err != nil {
		log.Printf("%-7s %-7s error stopping httpd: %s", logError, logName, err.Error())
	}

	z.listener = nil
	z.running = false

	log.Printf("%-7s %-7s service stopped", logInfo, logName)

}

// IsRunning indicates if server is running or not
func (z *Service) IsRunning() bool {
	z.Lock()
	defer z.Unlock()
	return z.running
}
