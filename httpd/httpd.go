package httpd

import (
	"crypto/tls"
	"errors"
	"github.com/rs/xhandler"
	"golang.org/x/net/context"
	"net"
	"net/http"
	"rakewire/api"
	"rakewire/fever"
	"rakewire/logger"
	"rakewire/model"
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
	api                *api.API
	appstart           time.Time
	database           model.Database
	debugMode          bool
	insecureSkipVerify bool
	listener           net.Listener
	listenHostPort     string // listening address
	publicHostPort     string // TODO: discard requests not made to this host
	running            bool
	tlsCertFile        string
	tlsKeyFile         string
	version            string
}

const (
	hContentType = "Content-Type"
)

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

	handler, errHandler := z.newHandler(localHostPort, tlsConfig)
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

	z.api.Shutdown()

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

func (z *Service) newHandler(endpoint string, tlsConfig *tls.Config) (http.Handler, error) {

	rootContext := context.Background()
	// TODO: grpc server.Stop() instead of quitters

	z.api = api.New(z.database, z.version, z.appstart)

	feverHandler := fever.New(z.database)
	opmlHandler := &opmlAPI{db: z.database}
	staticHandler := newStaticAPI(z.debugMode)

	apiHandler, err := z.api.NewHandler(rootContext, endpoint, tlsConfig)
	if err != nil {
		log.Debugf("cannot start API: %s", err.Error())
		return nil, err
	}
	apiServer := z.api.NewServer(tlsConfig)

	coreHandler := xhandler.HandlerFuncC(func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/fever/") {
			feverHandler.ServeHTTPC(ctx, w, r)
		} else if r.URL.Path == "/subscriptions.opml" {
			c := xhandler.Chain{}
			c.UseC(Authenticator(z.database))
			c.HandlerC(opmlHandler).ServeHTTPC(ctx, w, r)
		} else {
			staticHandler.ServeHTTPC(ctx, w, r)
		}
	})

	handler := xhandler.HandlerFuncC(func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		// TODO: logging, access
		if r.ProtoMajor == 2 && strings.Contains(r.Header.Get(hContentType), "application/grpc") {
			apiServer.ServeHTTP(w, r) // no context
		} else if strings.HasPrefix(r.URL.Path, "/api/") {
			apiHandler.ServeHTTP(w, r) // no context as this calls the grpc server
		} else {
			c := xhandler.Chain{}
			c.UseC(NoCache)
			c.HandlerC(coreHandler).ServeHTTPC(ctx, w, r)
		}
	})

	c := xhandler.Chain{}
	c.UseC(xhandler.CloseHandler)

	return xhandler.New(rootContext, c.HandlerC(handler)), nil // TODO: logging

}
