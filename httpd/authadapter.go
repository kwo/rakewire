package httpd

import (
	"github.com/gorilla/context"
	"net/http"
	"rakewire/auth"
	"rakewire/model"
)

// TODO: use x/net/context instead of gorilla context

// Authenticator authenticates requests, placing the user object in the request context
func Authenticator(db model.Database) Adapter {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			if user, err := auth.Authenticate(db, r.Header.Get("Authorization")); err == nil {
				context.Set(r, "user", user)
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
			h.ServeHTTP(w, r)

		})
	}
}
