package httpd

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"rakewire/db"
	"rakewire/logging"
)

// Service server
type Service struct {
	listener net.Listener
	Database db.Database
}

// Configuration configuration
type Configuration struct {
	Address string
	Port    int
	TLSCert string
	TLSKey  string
}

const (
	apiPrefix        = "/api"
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
	pathUI           = "../../../ui/webapp"
)

var (
	logger = logging.New("httpd")
)

type singleFileSystem struct {
	name string
	root http.FileSystem
}

func (z singleFileSystem) Open(name string) (http.File, error) {
	// ignore name and use z.name
	return z.root.Open(z.name)
}

// Start web service
func (z *Service) Start(cfg *Configuration, chErrors chan<- error) {

	if z.Database == nil {
		logger.Println("Cannot start httpd, no database provided")
		chErrors <- errors.New("No database")
		return
	}

	router := z.mainRouter(chErrors)
	if router == nil {
		return
	}
	mainHandler := Adapt(router, LogAdapter())

	l, err := net.Listen("tcp", fmt.Sprintf("%s:%d", cfg.Address, cfg.Port))
	if err != nil {
		logger.Printf("Cannot start listener: %s\n", err.Error())
		chErrors <- err
		return
	}
	// BACKLOG TLS wrap listener in tls.NewListener
	z.listener = l
	server := http.Server{
		Handler: mainHandler,
	}
	logger.Printf("Started httpd on http://%s", z.listener.Addr())
	err = server.Serve(z.listener)
	if err != nil && z.Running() {
		logger.Printf("Cannot start httpd: %s\n", err.Error())
		chErrors <- err
		return
	}

}

// Stop stop the server
func (z *Service) Stop() error {
	if l := z.listener; l != nil {
		z.listener = nil
		if err := l.Close(); err != nil {
			logger.Printf("Error stopping httpd: %s\n", err.Error())
			return err
		}
		logger.Println("Stopped httpd")
	}
	return nil
}

// Running indicates if server is running or not
func (z *Service) Running() bool {
	return z.listener != nil
}
