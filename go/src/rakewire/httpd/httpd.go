package httpd

import (
	"crypto/tls"
	"errors"
	"fmt"
	gorillaHandlers "github.com/gorilla/handlers"
	"log"
	"net"
	"net/http"
	"rakewire/middleware"
	"rakewire/model"
	"sync"
)

const (
	logName  = "[httpd]"
	logDebug = "[DEBUG]"
	logInfo  = "[INFO]"
	logWarn  = "[WARN]"
	logError = "[ERROR]"
)

const (
	httpdAccessLevel        = "httpd.accesslevel"
	httpdHost               = "httpd.host"
	httpdPort               = "httpd.port"
	httpdStaticLocal        = "httpd.staticlocal"
	httpdTLSPrivate         = "httpd.tls.private"
	httpdTLSPublic          = "httpd.tls.public"
	httpdUseTLS             = "httpd.tls.active"
	httpdAccessLevelDefault = "DEBUG"
	httpdHostDefault        = "localhost"
	httpdPortDefault        = 4444
	httpdStaticLocalDefault = false
	httpdTLSPrivateDefault  = ""
	httpdTLSPublicDefault   = ""
	httpdUseTLSDefault      = false
)

var (
	// ErrRestart indicates that the service cannot be started because it is already running.
	ErrRestart = errors.New("The service is already started")
)

// Service server
type Service struct {
	sync.Mutex
	database    model.Database
	listener    net.Listener
	running     bool
	accessLevel string
	host        string
	port        int
	staticLocal bool
	tlsPublic   string
	tlsPrivate  string
	useTLS      bool
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
func NewService(cfg *model.Configuration, database model.Database) *Service {
	return &Service{
		database:    database,
		accessLevel: cfg.Get(httpdAccessLevel, httpdAccessLevelDefault),
		host:        cfg.Get(httpdHost, httpdHostDefault),
		port:        cfg.GetInt(httpdPort, httpdPortDefault),
		tlsPublic:   cfg.Get(httpdTLSPublic, httpdTLSPublicDefault),
		tlsPrivate:  cfg.Get(httpdTLSPrivate, httpdTLSPrivateDefault),
		staticLocal: cfg.GetBool(httpdStaticLocal, httpdStaticLocalDefault),
		useTLS:      cfg.GetBool(httpdUseTLS, httpdUseTLSDefault),
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

	router, err := z.mainRouter(z.staticLocal)
	if err != nil {
		log.Printf("%-7s %-7s cannot load router: %s", logError, logName, err.Error())
		return err
	}
	mainHandler := middleware.Adapt(router, middleware.NoCache(), gorillaHandlers.CompressHandler, LogAdapter(z.accessLevel))

	// start http config

	if z.useTLS {
		cert, err := tls.X509KeyPair([]byte(z.tlsPublic), []byte(z.tlsPrivate))
		if err != nil {
			log.Printf("%-7s %-7s cannot start tls listener: %s", logError, logName, err.Error())
			return err
		}
		tlsConfig := &tls.Config{
			Certificates: []tls.Certificate{cert},
		}
		z.listener, err = tls.Listen("tcp", fmt.Sprintf("%s:%d", z.host, z.port), tlsConfig)
		if err != nil {
			log.Printf("%-7s %-7s cannot start tls listener: %s", logError, logName, err.Error())
			return err
		}
	} else {
		z.listener, err = net.Listen("tcp", fmt.Sprintf("%s:%d", z.host, z.port))
		if err != nil {
			log.Printf("%-7s %-7s cannot start listener: %s", logError, logName, err.Error())
			return err
		}
	}

	server := http.Server{
		Handler: mainHandler,
	}

	go server.Serve(z.listener)

	z.running = true
	if z.useTLS {
		log.Printf("%-7s %-7s service started on https://%s:%d", logInfo, logName, z.host, z.port)
	} else {
		log.Printf("%-7s %-7s service started on http://%s:%d", logInfo, logName, z.host, z.port)
	}
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
