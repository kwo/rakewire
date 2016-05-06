package web

//go:generate esc -private -o public.go -pkg web -prefix public public

import (
	"golang.org/x/net/context"
	"net/http"
	"regexp"
)

// Handler defines the http handler for the web package
type Handler struct {
	indexHandler      http.Handler
	html5RouteHandler http.Handler
	html5RouteMatcher *regexp.Regexp
	rootHandler       http.Handler
}

type oneFS struct {
	name string
	root http.FileSystem
}

func (z oneFS) Open(name string) (http.File, error) {
	// ignore name and use z.name
	return z.root.Open(z.name)
}

// New returns a new Handler. If debugMode is true files will be used directly from the local filesystem.
func New(debugMode bool) *Handler {

	var fs http.FileSystem
	if debugMode {
		fs = http.Dir("web/public") // debug mode assumes the application is being started from within the project root
	} else {
		fs = _escFS(false)
	}

	api := &Handler{}
	api.indexHandler = http.RedirectHandler("/", http.StatusMovedPermanently)
	api.html5RouteHandler = http.FileServer(oneFS{name: "/index.html", root: fs})
	api.html5RouteMatcher = regexp.MustCompile("^/[a-z0-9/-]+$")
	api.rootHandler = http.FileServer(fs)

	return api
}

// ServeHTTPC implements a context-aware http handler
func (z *Handler) ServeHTTPC(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/index.html" {
		z.indexHandler.ServeHTTP(w, r) // always redirect /index.html to /
	} else if z.html5RouteMatcher.MatchString(r.URL.Path) {
		z.html5RouteHandler.ServeHTTP(w, r) // HTML5 routes: any path without a dot (thus an extension)
	} else {
		z.rootHandler.ServeHTTP(w, r)
	}
}
