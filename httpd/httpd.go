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
	httpdAddress            = "httpd.address"
	httpdHost               = "httpd.host"
	httpdPort               = "httpd.port"
	httpdTLSPort            = "httpd.tls.port"
	httpdTLSPrivate         = "httpd.tls.private"
	httpdTLSPublic          = "httpd.tls.public"
	httpdUseTLS             = "httpd.tls.active"
	httpdAccessLevelDefault = "DEBUG"
	httpdAddressDefault     = ""
	httpdHostDefault        = "localhost"
	httpdPortDefault        = 8888
	httpdTLSPortDefault     = 4444
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
	tlsListener net.Listener
	running     bool
	accessLevel string
	address     string // binding address, empty string means 0.0.0.0
	host        string // discard requests not made to this host
	port        int
	tlsPort     int
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
		tlsPort:     cfg.GetInt(httpdTLSPort, httpdTLSPortDefault),
		tlsPublic:   cfg.Get(httpdTLSPublic, httpdTLSPublicDefault),
		tlsPrivate:  cfg.Get(httpdTLSPrivate, httpdTLSPrivateDefault),
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

	if z.useTLS {
		if err := z.startHTTPS(); err != nil {
			return err
		}
	}

	if err := z.startHTTP(); err != nil {
		return err
	}

	z.running = true

	return nil

}

func (z *Service) startHTTP() error {

	restrictToStatusOnly := z.useTLS
	router, err := z.mainRouter(restrictToStatusOnly)
	if err != nil {
		log.Printf("%-7s %-7s cannot load router: %s", logError, logName, err.Error())
		return err
	}

	mainHandler := middleware.Adapt(router, middleware.NoCache(), gorillaHandlers.CompressHandler, LogAdapter(z.accessLevel))

	z.listener, err = net.Listen("tcp", fmt.Sprintf("%s:%d", z.address, z.port))
	if err != nil {
		log.Printf("%-7s %-7s cannot start listener: %s", logError, logName, err.Error())
		return err
	}

	server := http.Server{
		Handler: mainHandler,
	}
	go server.Serve(z.listener)

	log.Printf("%-7s %-7s service started on http://%s:%d", logInfo, logName, z.address, z.port)

	return nil

}

func (z *Service) startHTTPS() error {

	router, err := z.mainRouter()
	if err != nil {
		log.Printf("%-7s %-7s cannot load router: %s", logError, logName, err.Error())
		return err
	}

	mainHandler := middleware.Adapt(router, middleware.NoCache(), gorillaHandlers.CompressHandler, LogAdapter(z.accessLevel))

	cert, err := tls.X509KeyPair([]byte(z.tlsPublic), []byte(z.tlsPrivate))
	if err != nil {
		log.Printf("%-7s %-7s cannot create tls key pair: %s", logError, logName, err.Error())
		return err
	}
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
	}
	z.tlsListener, err = tls.Listen("tcp", fmt.Sprintf("%s:%d", z.address, z.tlsPort), tlsConfig)
	if err != nil {
		log.Printf("%-7s %-7s cannot start tls listener: %s", logError, logName, err.Error())
		return err
	}

	server := http.Server{
		Handler: mainHandler,
	}

	go server.Serve(z.tlsListener)

	log.Printf("%-7s %-7s service started on https://%s:%d", logInfo, logName, z.host, z.tlsPort)

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
