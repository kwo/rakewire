package httpd

//go:generate esc -o static.go -pkg httpd -prefix $PROJECT_HOME/web $PROJECT_HOME/web/public

import (
	"github.com/gorilla/mux"
	"net/http"
	"rakewire/fever"
	"rakewire/middleware"
	"rakewire/rest"
)

func (z *Service) mainRouter(useLocal bool) (*mux.Router, error) {

	router := mux.NewRouter()

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

func sendOK(w http.ResponseWriter, req *http.Request) {
	sendError(w, http.StatusOK)
}

func badMediaType(w http.ResponseWriter, req *http.Request) {
	sendError(w, http.StatusUnsupportedMediaType)
}

func badRequest(w http.ResponseWriter, req *http.Request) {
	sendError(w, http.StatusBadRequest)
}

func notFound(w http.ResponseWriter, req *http.Request) {
	sendError(w, http.StatusNotFound)
}

func notSupported(w http.ResponseWriter, req *http.Request) {
	sendError(w, http.StatusMethodNotAllowed)
}

func sendError(w http.ResponseWriter, code int) {
	http.Error(w, http.StatusText(code), code)
}

type oneFS struct {
	name string
	root http.FileSystem
}

func (z oneFS) Open(name string) (http.File, error) {
	// ignore name and use z.name
	return z.root.Open(z.name)
}
