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
func (z *Service) Start(chErrors chan<- error) {

	z.Lock()
	defer z.Unlock()
	if z.running {
		log.Printf("%-7s %-7s Server already started, exiting...", logWarn, logName)
		return
	}

	if z.database == nil {
		log.Printf("%-7s %-7s Cannot start httpd, no database provided", logError, logName)
		chErrors <- errors.New("No database")
		return
	}

	router, err := z.mainRouter(z.cfg.UseLocal)
	if err != nil {
		log.Printf("%-7s %-7s Cannot load router: %s", logError, logName, err.Error())
		chErrors <- err
		return
	}
	mainHandler := Adapt(router, LogAdapter(z.cfg.AccessLog))

	z.listener, err = net.Listen("tcp", fmt.Sprintf("%s:%d", z.cfg.Address, z.cfg.Port))
	if err != nil {
		log.Printf("%-7s %-7s Cannot start listener: %s", logError, logName, err.Error())
		chErrors <- err
		return
	}

	server := http.Server{
		Handler: mainHandler,
	}

	z.running = true
	log.Printf("%-7s %-7s Started httpd on http://%s:%d", logInfo, logName, z.cfg.Address, z.cfg.Port)

	go func() {
		err = server.Serve(z.listener)
		if err != nil && z.IsRunning() {
			log.Printf("%-7s %-7s Cannot start httpd: %s", logError, logName, err.Error())
			chErrors <- err
		}
	}()

}

// Stop stop the server
func (z *Service) Stop() {

	z.Lock()
	defer z.Unlock()
	if !z.running {
		log.Printf("%-7s %-7s Server already stopped, exiting...", logWarn, logName)
		return
	}

	if err := z.listener.Close(); err != nil {
		log.Printf("%-7s %-7s Error stopping httpd: %s", logError, logName, err.Error())
	}

	z.listener = nil
	z.running = false

	log.Printf("%-7s %-7s Stopped httpd", logInfo, logName)

}

// IsRunning indicates if server is running or not
func (z *Service) IsRunning() bool {
	z.Lock()
	defer z.Unlock()
	return z.running
}
