package httpd

import (
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"rakewire/db"
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
	listener net.Listener
	Database db.Database
}

// Configuration configuration
type Configuration struct {
	AccessLog string
	Address   string
	Port      int
	TLSCert   string
	TLSKey    string
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

// Start web service
func (z *Service) Start(cfg *Configuration, chErrors chan<- error) {

	if z.Database == nil {
		log.Printf("%-7s %-7s Cannot start httpd, no database provided", logError, logName)
		chErrors <- errors.New("No database")
		return
	}

	router, err := z.mainRouter()
	if err != nil {
		log.Printf("%-7s %-7s Cannot load router: %s", logError, logName, err.Error())
		chErrors <- err
		return
	}
	mainHandler := Adapt(router, LogAdapter(cfg.AccessLog))

	l, err := net.Listen("tcp", fmt.Sprintf("%s:%d", cfg.Address, cfg.Port))
	if err != nil {
		log.Printf("%-7s %-7s Cannot start listener: %s", logError, logName, err.Error())
		chErrors <- err
		return
	}
	z.listener = l
	server := http.Server{
		Handler: mainHandler,
	}
	log.Printf("%-7s %-7s Started httpd on http://%s", logInfo, logName, z.listener.Addr())
	err = server.Serve(z.listener)
	if err != nil && z.Running() {
		log.Printf("%-7s %-7s Cannot start httpd: %s", logError, logName, err.Error())
		chErrors <- err
		return
	}

}

// Stop stop the server
func (z *Service) Stop() error {
	if l := z.listener; l != nil {
		z.listener = nil
		if err := l.Close(); err != nil {
			log.Printf("%-7s %-7s Error stopping httpd: %s", logError, logName, err.Error())
			return err
		}
		log.Printf("%-7s %-7s Stopped httpd", logInfo, logName)
	}
	return nil
}

// Running indicates if server is running or not
func (z *Service) Running() bool {
	return z.listener != nil
}
