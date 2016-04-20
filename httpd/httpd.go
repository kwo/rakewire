package httpd

import (
	"crypto/tls"
	"errors"
	"fmt"
	gorillaHandlers "github.com/gorilla/handlers"
	"net"
	"net/http"
	"rakewire/api"
	"rakewire/logger"
	"rakewire/middleware"
	"rakewire/model"
	"sync"
	"time"
)

const (
	httpdAddress           = "httpd.address"
	httpdHost              = "httpd.host"
	httpdPort              = "httpd.port"
	httpdTLSActive         = "httpd.tls.active"
	httpdTLSPort           = "httpd.tls.port"
	httpdTLSPrivate        = "httpd.tls.private"
	httpdTLSPublic         = "httpd.tls.public"
	httpdAddressDefault    = ""
	httpdHostDefault       = "localhost"
	httpdPortDefault       = 8888
	httpdTLSActiveDefault  = false
	httpdTLSPortDefault    = 4444
	httpdTLSPrivateDefault = ""
	httpdTLSPublicDefault  = ""
)

var (
	// ErrRestart indicates that the service cannot be started because it is already running.
	ErrRestart = errors.New("The service is already started")
	log        = logger.New("httpd")
)

// Service server
type Service struct {
	sync.Mutex
	database    model.Database
	listener    net.Listener
	tlsListener net.Listener
	running     bool
	address     string // binding address, empty string means 0.0.0.0
	host        string // TODO: discard requests not made to this host
	port        int
	tlsActive   bool
	tlsPort     int
	tlsPublic   string
	tlsPrivate  string
	version     string
	appstart    time.Time
}

const (
	hContentType = "Content-Type"
	mGet         = "GET"
	mimeText     = "text/plain; charset=utf-8"
)

// NewService creates a new httpd service.
func NewService(cfg *model.Configuration, database model.Database) *Service {
	return &Service{
		database:   database,
		address:    cfg.GetStr(httpdAddress, httpdAddressDefault),
		host:       cfg.GetStr(httpdHost, httpdHostDefault),
		port:       cfg.GetInt(httpdPort, httpdPortDefault),
		tlsActive:  cfg.GetBool(httpdTLSActive, httpdTLSActiveDefault),
		tlsPort:    cfg.GetInt(httpdTLSPort, httpdTLSPortDefault),
		tlsPublic:  cfg.GetStr(httpdTLSPublic, httpdTLSPublicDefault),
		tlsPrivate: cfg.GetStr(httpdTLSPrivate, httpdTLSPrivateDefault),
		version:    cfg.GetStr("app.version", "Rakewire"),
		appstart:   time.Unix(cfg.GetInt64("app.start", time.Now().Unix()), 0).Truncate(time.Second),
	}
}

// Start web service
func (z *Service) Start() error {

	z.Lock()
	defer z.Unlock()
	if z.running {
		log.Debugf("service already started, exiting...")
		return ErrRestart
	}

	if z.database == nil {
		log.Debugf("cannot start httpd, no database provided")
		return errors.New("No database")
	}

	if z.tlsActive {
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

	restrictToStatusOnly := z.tlsActive
	router, err := z.mainRouter(restrictToStatusOnly)
	if err != nil {
		log.Debugf("cannot load router: %s", err.Error())
		return err
	}

	mainHandler := middleware.Adapt(router, middleware.NoCache(), gorillaHandlers.CompressHandler, LogAdapter())

	endpoint := fmt.Sprintf("%s:%d", z.address, z.port)
	z.listener, err = net.Listen("tcp", endpoint)
	if err != nil {
		log.Debugf("cannot start listener: %s", err.Error())
		return err
	}

	server := http.Server{
		Handler: mainHandler,
	}
	go server.Serve(z.listener)

	log.Debugf("service started on http://%s", endpoint)

	return nil

}

func (z *Service) startHTTPS() error {

	cert, err := tls.X509KeyPair([]byte(z.tlsPublic), []byte(z.tlsPrivate))
	if err != nil {
		log.Debugf("cannot create tls key pair: %s", err.Error())
		return err
	}

	endpointListen := fmt.Sprintf("%s:%d", z.address, z.tlsPort)
	endpointConnect := fmt.Sprintf("%s:%d", z.host, z.tlsPort)
	tlsConfigListen := &tls.Config{
		Certificates: []tls.Certificate{cert},
	}
	tlsConfigConnect := &tls.Config{
		ServerName:   z.host,
		Certificates: []tls.Certificate{cert},
	}

	// TODO
	// incorporate routes right here into this function
	// migrate status handler in this package to api package
	// migrate opml
	// shitcan rest package
	// authentication
	// add static pages

	router, err := z.mainRouter()
	if err != nil {
		log.Debugf("cannot load router: %s", err.Error())
		return err
	}

	apiHandler, apiGRPCServer, err := api.NewAPI(z.database, z.version, z.appstart).Router(endpointConnect, tlsConfigConnect)
	if err != nil {
		log.Debugf("cannot start API: %s", err.Error())
		return err
	}

	// TODO: api authentication

	router.PathPrefix("/api").Handler(apiHandler)

	mainHandler := middleware.Adapt(router, middleware.NoCache(), gorillaHandlers.CompressHandler, LogAdapter())

	z.tlsListener, err = tls.Listen("tcp", endpointListen, tlsConfigListen)
	if err != nil {
		log.Debugf("cannot start tls listener: %s", err.Error())
		return err
	}

	server := http.Server{
		Addr:      endpointListen,
		Handler:   grpcHandler(apiGRPCServer, mainHandler),
		TLSConfig: tlsConfigListen,
	}

	go server.Serve(z.tlsListener)

	log.Debugf("service started on https://%s", endpointListen)

	return nil

}

// Stop stop the server
func (z *Service) Stop() {

	z.Lock()
	defer z.Unlock()
	if !z.running {
		log.Debugf("service already stopped, exiting...")
		return
	}

	if err := z.listener.Close(); err != nil {
		log.Debugf("error stopping httpd: %s", err.Error())
	}

	z.listener = nil
	z.running = false

	log.Debugf("service stopped")

}

// IsRunning indicates if server is running or not
func (z *Service) IsRunning() bool {
	z.Lock()
	defer z.Unlock()
	return z.running
}
