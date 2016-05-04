package httpd

import (
	"golang.org/x/net/context"
	"net/http"
	"rakewire/web"
	"regexp"
)

type staticAPI struct {
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

func newStaticAPI(debugMode bool) *staticAPI {

	var fs http.FileSystem
	if debugMode {
		fs = http.Dir("web/public") // debug mode assumes the application is being started from within the project root
	} else {
		fs = web.FS(false)
	}

	api := &staticAPI{}
	api.indexHandler = http.RedirectHandler("/", http.StatusMovedPermanently)
	api.html5RouteHandler = http.FileServer(oneFS{name: "/index.html", root: fs})
	api.html5RouteMatcher = regexp.MustCompile("^/[a-z0-9/-]+$")
	api.rootHandler = http.FileServer(fs)

	return api
}

func (z *staticAPI) ServeHTTPC(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/index.html" {
		z.indexHandler.ServeHTTP(w, r) // always redirect /index.html to /
	} else if z.html5RouteMatcher.MatchString(r.URL.Path) {
		z.html5RouteHandler.ServeHTTP(w, r) // HTML5 routes: any path without a dot (thus an extension)
	} else {
		z.rootHandler.ServeHTTP(w, r)
	}
}
