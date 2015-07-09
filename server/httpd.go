package server

import (
	"fmt"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"net"
	"net/http"
	m "rakewire.com/model"
)

// Httpd server
type Httpd struct {
	listener net.Listener
}

var (
	logger = NewInternalLogger()
)

// Start web service
func (z *Httpd) Start(cfg m.HttpdConfiguration) {

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
		logger.Fatal(err)
		return
	}
	z.listener = l
	server := http.Server{
		Handler: n,
	}
	logger.Printf("Started httpd on http://%s", z.listener.Addr())
	logger.Fatal(server.Serve(z.listener))

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
