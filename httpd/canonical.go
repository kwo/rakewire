package httpd

import (
	"golang.org/x/net/context"
	"net/http"
	"strings"
)

// CanonicalHost is HTTP middleware that re-directs requests to the canonical
// domain. It accepts a domain and a status code (e.g. 301 or 302) and
// re-directs clients to this domain. The existing request path is maintained.
func CanonicalHost(hostport string, code int) Middleware {

	// remove standard port, if present
	host := strings.TrimSuffix(hostport, ":443")

	return func(next HandlerC) HandlerC {
		return HandlerFuncC(func(ctx context.Context, w http.ResponseWriter, r *http.Request) {

			if r.Host != host {
				if len(r.URL.Scheme) == 0 {
					r.URL.Scheme = "https"
				}
				r.URL.Host = host
				http.Redirect(w, r, r.URL.String(), code)
				return
			}

			next.ServeHTTPC(ctx, w, r)

		})
	}
}
