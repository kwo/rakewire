package server

import (
	//"fmt"
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
	router.PathPrefix("/").Handler(http.FileServer(http.Dir(cfg.WebAppDir)))

	n := negroni.New()
	n.Use(negroni.NewLogger())
	n.UseHandler(router)
	n.Run(fmt.Sprintf("%s:%d", cfg.Address, cfg.Port))

}
