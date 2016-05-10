package httpd

import (
	"crypto/tls"
	"errors"
	"github.com/kwo/rakewire/api"
	"github.com/kwo/rakewire/fever"
	"github.com/kwo/rakewire/logger"
	"github.com/kwo/rakewire/model"
	"github.com/kwo/rakewire/web"
	"golang.org/x/net/context"
	"net"
	"net/http"
	"strings"
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
	DebugMode          bool
	InsecureSkipVerify bool
	ListenHostPort     string
	PublicHostPort     string
	TLSCertFile        string
	TLSKeyFile         string
}

// Service server
type Service struct {
	sync.Mutex
	appstart           int64
	database           model.Database
	debugMode          bool
	insecureSkipVerify bool
	listener           net.Listener
	listenHostPort     string // listening address
	publicHostPort     string
	running            bool
	tlsCertFile        string
	tlsKeyFile         string
	version            string
}

// NewService creates a new httpd service.
func NewService(cfg *Configuration, database model.Database, version string, appStart int64) *Service {
	return &Service{
		database:           database,
		debugMode:          cfg.DebugMode,
		insecureSkipVerify: cfg.InsecureSkipVerify,
		listenHostPort:     cfg.ListenHostPort,
		publicHostPort:     cfg.PublicHostPort,
		tlsCertFile:        cfg.TLSCertFile,
		tlsKeyFile:         cfg.TLSKeyFile,
		version:            version,
		appstart:           appStart,
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

	handler := z.newHandler()

	server := http.Server{
		Addr:      z.listenHostPort,
		Handler:   handler,
		TLSConfig: tlsConfig,
	}

	go server.Serve(z.listener)

	log.Infof("listening on %s, reachable at %s", z.listenHostPort, z.publicHostPort)

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

func (z *Service) newHandler() http.Handler {

	apiPath := "/api/"
	apiHandler := Chain(api.New(z.database, apiPath, z.version, z.appstart), Authorize(""))
	feverPath := "/fever/"
	feverHandler := fever.New(z.database)
	webHandler := web.New(z.debugMode)

	router := HandlerFuncC(func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, apiPath) {
			apiHandler.ServeHTTPC(ctx, w, r)
		} else if strings.HasPrefix(r.URL.Path, feverPath) {
			feverHandler.ServeHTTPC(ctx, w, r)
		} else {
			webHandler.ServeHTTPC(ctx, w, r)
		}
	})

	handler := Chain(
		router,
		NoCache(),
		CanonicalHost(z.publicHostPort, http.StatusMovedPermanently),
		AccessLog(),
		Authenticate(z.database),
		TimeoutHandler(10*time.Second), // TODO: configurable request timeout
		CloseHandler(),
	)

	return Adapt(context.Background(), handler) // TODO: logging

}
