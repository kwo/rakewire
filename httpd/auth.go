package httpd

import (
	"github.com/kwo/rakewire/auth"
	"github.com/kwo/rakewire/model"
	"golang.org/x/net/context"
	"net/http"
)

// Authenticator authenticates requests, placing the user object in the request context
func Authenticator(db model.Database) Middleware {
	return func(next HandlerC) HandlerC {
		return HandlerFuncC(func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
			if user, err := auth.Authenticate(db, r.Header.Get("Authorization")); err == nil {
				ctx = context.WithValue(ctx, "user", user)
			} else if err == auth.ErrUnauthenticated {
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			} else if err == auth.ErrUnauthorized {
				http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
				return
			} else if err == auth.ErrBadHeader {
				http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
				return
			} else {
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}

			// Call the next handler on success.
			next.ServeHTTPC(ctx, w, r)

		})
	}
}
