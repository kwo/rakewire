package httpd

import (
	"errors"
	"fmt"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"net"
	"net/http"
	"rakewire.com/db"
	m "rakewire.com/model"
)

// Httpd server
type Httpd struct {
	listener net.Listener
	Database db.Database
}

var (
	logger = NewInternalLogger()
)

// Start web service
func (z *Httpd) Start(cfg m.HttpdConfiguration) error {

	if z.Database == nil {
		logger.Println("Cannot start httpd, no database provided")
		return errors.New("No database")
	}

	router := mux.NewRouter()

	// api router
	apiRouter := APIRouter()
	router.PathPrefix("/api").Handler(negroni.New(
		negroni.Wrap(apiRouter),
	))

	// static web site
	router.PathPrefix("/").Handler(negroni.New(
		negroni.Wrap(http.FileServer(http.Dir(cfg.WebAppDir))),
	))

	n := negroni.New()
	n.Use(negroni.NewRecovery())
	n.Use(logger)
	n.UseHandler(router)

	l, err := net.Listen("tcp", fmt.Sprintf("%s:%d", cfg.Address, cfg.Port))
	if err != nil {
		logger.Printf("Cannot start listener: %s\n", err.Error())
		return err
	}
	z.listener = l
	server := http.Server{
		Handler: n,
	}
	logger.Printf("Started httpd on http://%s", z.listener.Addr())
	err = server.Serve(z.listener)
	if err != nil {
		return err
	}

	return nil

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
