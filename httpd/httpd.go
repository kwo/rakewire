package httpd

import (
	"errors"
	"fmt"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/phyber/negroni-gzip/gzip"
	"net"
	"net/http"
	m "rakewire.com/model"
)

// Httpd server
type Httpd struct {
	listener net.Listener
	Database m.Database
}

const (
	apiPrefix        = "/api"
	hAcceptEncoding  = "Accept-Encoding"
	hContentEncoding = "Content-Encoding"
	hContentLength   = "Content-Length"
	hContentType     = "Content-Type"
	mGet             = "GET"
	mPost            = "POST"
	mPut             = "PUT"
	mimeHTML         = "text/html; charset=utf-8"
	mimeJSON         = "application/json"
	mimeText         = "text/plain; charset=utf-8"
)

var (
	logger = newInternalLogger()
)

// Start web service
func (z *Httpd) Start(cfg *m.HttpdConfiguration, chErrors chan error) {

	if z.Database == nil {
		logger.Println("Cannot start httpd, no database provided")
		chErrors <- errors.New("No database")
		return
	}

	router := mux.NewRouter()

	// api router
	router.PathPrefix(apiPrefix).Handler(negroni.New(
		negroni.Wrap(z.apiRouter(apiPrefix)),
	))

	// static web site
	router.PathPrefix("/").Handler(negroni.New(
		gzip.Gzip(gzip.BestCompression),
		negroni.Wrap(http.FileServer(http.Dir(cfg.WebAppDir))),
	))

	n := negroni.New()
	n.Use(negroni.NewRecovery())
	n.Use(logger)
	n.UseHandler(router)

	l, err := net.Listen("tcp", fmt.Sprintf("%s:%d", cfg.Address, cfg.Port))
	if err != nil {
		logger.Printf("Cannot start listener: %s\n", err.Error())
		chErrors <- err
		return
	}
	z.listener = l
	server := http.Server{
		Handler: n,
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
func (z *Httpd) Stop() error {
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
func (z *Httpd) Running() bool {
	return z.listener != nil
}
