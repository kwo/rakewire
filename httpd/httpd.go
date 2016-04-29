package httpd

import (
	"crypto/tls"
	"errors"
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
	ListenHostPort     string
	PublicHostPort     string
	InsecureSkipVerify bool
	TLSCertFile        string
	TLSKeyFile         string
}

// Service server
type Service struct {
	sync.Mutex
	database           model.Database
	listener           net.Listener
	api                *api.API
	running            bool
	listenHostPort     string // listening address
	publicHostPort     string // TODO: discard requests not made to this host
	insecureSkipVerify bool
	tlsCertFile        string
	tlsKeyFile         string
	version            string
	appstart           time.Time
}

const (
	hContentType = "Content-Type"
	mGet         = "GET"
	mPut         = "PUT"
)

// NewService creates a new httpd service.
func NewService(cfg *Configuration, database model.Database, version string, appStart int64) *Service {
	return &Service{
		database:           database,
		listenHostPort:     cfg.ListenHostPort,
		publicHostPort:     cfg.PublicHostPort,
		insecureSkipVerify: cfg.InsecureSkipVerify,
		tlsCertFile:        cfg.TLSCertFile,
		tlsKeyFile:         cfg.TLSKeyFile,
		version:            version,
		appstart:           time.Unix(appStart, 0).Truncate(time.Second),
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
	log.Infof("listen:   %s", z.listenHostPort)
	log.Infof("public:   %s", z.publicHostPort)
	log.Infof("insecure: %t", z.insecureSkipVerify)
	log.Infof("tls cert: %s", z.tlsCertFile)
	log.Infof("tls key:  %s", z.tlsKeyFile)

	// extract the hostname from public hostport - assign to tlsConfig servername
	publicFQDN, _, errSplitPublic := net.SplitHostPort(z.publicHostPort)
	if errSplitPublic != nil {
		log.Debugf("cannot split public hostport: %s", errSplitPublic.Error())
		return errSplitPublic
	}
	cert, err := tls.LoadX509KeyPair(z.tlsCertFile, z.tlsKeyFile)
	if err != nil {
		log.Debugf("cannot create tls key pair: %s", err.Error())
		return err
	}
	tlsConfig := &tls.Config{
		ServerName:         publicFQDN,
		InsecureSkipVerify: z.insecureSkipVerify,
		Certificates:       []tls.Certificate{cert},
	}

	z.listener, err = tls.Listen("tcp", z.listenHostPort, tlsConfig)
	if err != nil {
		log.Debugf("cannot start tls listener: %s", err.Error())
		return err
	}

	localHostname, localPort, errSplitLocal := net.SplitHostPort(z.listenHostPort)
	if errSplitLocal != nil {
		log.Debugf("cannot split listen hostport: %s", errSplitLocal.Error())
		return errSplitLocal
	}
	if localHostname == net.IPv4zero.String() {
		localHostname = "localhost"
	}
	localHostPort := net.JoinHostPort(localHostname, localPort)

	handler, errHandler := z.router(localHostPort, tlsConfig)
	if errHandler != nil {
		return errHandler
	}

	server := http.Server{
		Addr:      z.listenHostPort,
		Handler:   handler,
		TLSConfig: tlsConfig,
	}

	go server.Serve(z.listener)

	log.Infof("listening on %s, reachable at %s, local %s", z.listenHostPort, z.publicHostPort, localHostPort)

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

	z.api.Stop()

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

	z.api = api.NewAPI(z.database, z.version, z.appstart)

	apiHandler, apiGRPCServer, err := z.api.Router(endpoint, tlsConfig)
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
