package httpd

import (
	"errors"
	"fmt"
	"github.com/GeertJohan/go.rice"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/phyber/negroni-gzip/gzip"
	"net"
	"net/http"
	"rakewire/db"
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
	logger = newInternalLogger()
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
func (z *Service) Start(cfg *Configuration, chErrors chan error) {

	if z.Database == nil {
		logger.Println("Cannot start httpd, no database provided")
		chErrors <- errors.New("No database")
		return
	}

	router := mux.NewRouter()

	// api router
	router.PathPrefix(apiPrefix).Handler(negroni.New(
		NoCache(),
		negroni.Wrap(z.apiRouter(apiPrefix)),
	))

	// static web site
	box, err := rice.FindBox(pathUI)
	if err != nil {
		logger.Printf("Cannot find box: %s\n", err.Error())
		chErrors <- err
		return
	}
	bfs := box.HTTPBox()

	// HTML5 routes: any path without an entension
	sfs := singleFileSystem{name: "/index.html", root: bfs}
	router.Path("/{route:[a-z0-9/-]+}").Handler(negroni.New(
		gzip.Gzip(gzip.BestCompression),
		NoCache(),
		negroni.Wrap(http.FileServer(sfs)),
	))

	router.PathPrefix("/").Handler(negroni.New(
		gzip.Gzip(gzip.BestCompression),
		NoCache(),
		negroni.Wrap(http.FileServer(bfs)),
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
	// BACKLOG TLS wrap listener in tls.NewListener
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
