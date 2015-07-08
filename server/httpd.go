package server

import (
	"fmt"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"net/http"
	m "rakewire.com/model"
)

// Serve web service
func Serve(cfg m.HttpdConfiguration) {

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

	logger := NewInternalLogger()

	n := negroni.New()
	n.Use(negroni.NewRecovery())
	n.Use(logger)
	n.UseHandler(router)

	addr := fmt.Sprintf("%s:%d", cfg.Address, cfg.Port)
	logger.Printf("listening on http://%s", addr)
	logger.Fatal(http.ListenAndServe(addr, n))

}
