package httpd

//go:generate esc -o static.go -pkg httpd -prefix $PROJECT_HOME/web $PROJECT_HOME/web/public

import (
	"github.com/gorilla/mux"
	"net/http"
	"rakewire/fever"
	"rakewire/middleware"
	"rakewire/rest"
)

func (z *Service) mainRouter(flags ...bool) (*mux.Router, error) {

	flagStatusOnly := len(flags) > 0 && flags[0]

	router := mux.NewRouter()

	if !flagStatusOnly {

		// rest api router
		restPrefix := "/api"
		restAPI := rest.NewAPI(restPrefix, z.database)
		router.PathPrefix(restPrefix).Handler(
			middleware.Adapt(restAPI.Router(), middleware.BasicAuth(&middleware.BasicAuthOptions{
				Database: z.database, Realm: "Rakewire",
			})),
		)

		// fever api router
		feverPrefix := "/fever/"
		feverAPI := fever.NewAPI(feverPrefix, z.database)
		router.PathPrefix(feverPrefix).Handler(
			feverAPI.Router(),
		)

	}

	// status api router
	router.Path("/status").Methods(mGet).HandlerFunc(statusHandler)

	// fs := Dir(useLocal, "/public")
	// ofs := oneFS{name: "/index.html", root: fs}
	//
	// // HTML5 routes: any path without a dot (thus an extension)
	// router.Path("/{route:[a-z0-9/-]+}").Handler(
	// 	http.FileServer(ofs),
	// )
	//
	// // always redirect /index.html to /
	// router.Path("/index.html").Handler(
	// 	http.RedirectHandler("/", http.StatusMovedPermanently),
	// )
	//
	// // static web site
	// router.PathPrefix("/").Handler(
	// 	http.FileServer(fs),
	// )

	return router, nil

}

type oneFS struct {
	name string
	root http.FileSystem
}

func (z oneFS) Open(name string) (http.File, error) {
	// ignore name and use z.name
	return z.root.Open(z.name)
}
