package httpd

import (
	"crypto/tls"
	"errors"
	"fmt"
	gorillaHandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"net"
	"net/http"
	"rakewire/api"
	"rakewire/fever"
	"rakewire/logger"
	"rakewire/model"
	"sync"
	"time"
)

var (
	// ErrRestart indicates that the service cannot be started because it is already running.
	ErrRestart = errors.New("The service is already started")
	log        = logger.New("httpd")
)

// Configuration contains all parameters for the httpd service
type Configuration struct {
	Address     string
	Host        string
	Port        int
	TLSCertFile string
	TLSKeyFile  string
}

// Service server
type Service struct {
	sync.Mutex
	database    model.Database
	listener    net.Listener
	running     bool
	address     string // binding address, empty string means 0.0.0.0
	host        string // TODO: discard requests not made to this host
	port        int
	tlsCertFile string
	tlsKeyFile  string
	version     string
	appstart    time.Time
}

const (
	hContentType = "Content-Type"
	mGet         = "GET"
	mPut         = "PUT"
)

// NewService creates a new httpd service.
func NewService(cfg *Configuration, database model.Database, version string, appStart int64) *Service {
	return &Service{
		database:    database,
		address:     cfg.Address,
		host:        cfg.Host,
		port:        cfg.Port,
		tlsCertFile: cfg.TLSCertFile,
		tlsKeyFile:  cfg.TLSKeyFile,
		version:     version,
		appstart:    time.Unix(appStart, 0).Truncate(time.Second),
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

	log.Infof("starting...")
	log.Infof("address:  %s", z.address)
	log.Infof("host:     %s", z.host)
	log.Infof("port:     %d", z.port)
	log.Infof("tls cert: %s", z.tlsCertFile)
	log.Infof("tls key:  %s", z.tlsKeyFile)

	cert, err := tls.LoadX509KeyPair(z.tlsCertFile, z.tlsKeyFile)
	if err != nil {
		log.Debugf("cannot create tls key pair: %s", err.Error())
		return err
	}

	endpointListen := fmt.Sprintf("%s:%d", z.address, z.port)
	endpointConnect := fmt.Sprintf("%s:%d", z.host, z.port)
	tlsConfig := &tls.Config{
		ServerName:   z.host,
		Certificates: []tls.Certificate{cert},
	}

	z.listener, err = tls.Listen("tcp", endpointListen, tlsConfig)
	if err != nil {
		log.Debugf("cannot start tls listener: %s", err.Error())
		return err
	}

	handler, errHandler := z.router(endpointConnect, tlsConfig)
	if errHandler != nil {
		return errHandler
	}

	server := http.Server{
		Addr:      endpointListen,
		Handler:   handler,
		TLSConfig: tlsConfig,
	}

	go server.Serve(z.listener)

	log.Infof("listening on %s, reachable at %s", endpointListen, endpointConnect)

	z.running = true

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

	log.Infof("stopped")

}

// IsRunning indicates if server is running or not
func (z *Service) IsRunning() bool {
	z.Lock()
	defer z.Unlock()
	return z.running
}

func (z *Service) router(endpoint string, tlsConfig *tls.Config) (http.Handler, error) {

	router := mux.NewRouter()

	// fever api router
	// no authentication necessary as it uses apiKey and feverhash
	feverPrefix := "/fever/"
	feverAPI := fever.NewAPI(feverPrefix, z.database)
	router.PathPrefix(feverPrefix).Handler(
		feverAPI.Router(),
	)

	apiHandler, apiGRPCServer, err := api.NewAPI(z.database, z.version, z.appstart).Router(endpoint, tlsConfig)
	if err != nil {
		log.Debugf("cannot start API: %s", err.Error())
		return nil, err
	}
	// Note: the api prefix must match the path in api/pb/api.proto
	// No authentication necessary because each GRPC method authenticates
	router.PathPrefix("/api").Handler(apiHandler)

	// oddballs router
	oddballsAPI := &oddballs{db: z.database}
	router.PathPrefix("/").Handler(
		Adapt(oddballsAPI.router(), Authenticator(z.database)),
	)

	mainHandler := Adapt(router, NoCache(), gorillaHandlers.CompressHandler, LogAdapter())

	return grpcHandler(apiGRPCServer, mainHandler), nil

}
