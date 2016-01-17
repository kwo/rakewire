package middleware

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/gorilla/context"
	"net/http"
	"rakewire/model"
	"strings"
)

// BasicAuthOptions stores the configuration for HTTP Basic Authentication.
type BasicAuthOptions struct {
	Realm    string
	Database model.Database
}

// BasicAuth provides HTTP middleware for protecting URIs with HTTP Basic Authentication
// as per RFC 2617. The server authenticates a user:password combination provided in the
// "Authorization" HTTP header.
func BasicAuth(z *BasicAuthOptions) Adapter {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			// Check that the provided details match
			if !z.authenticate(r) {
				z.requestAuth(w, r)
				return
			}

			// Call the next handler on success.
			h.ServeHTTP(w, r)

		})
	}
}

// Require authentication, and serve our error handler otherwise.
func (z *BasicAuthOptions) requestAuth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("WWW-Authenticate", fmt.Sprintf("Basic realm=%q", z.Realm))
	http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
}

// authenticate retrieves and then validates the user:password combination provided in
// the request header. Returns 'false' if the user has not successfully authenticated.
func (z *BasicAuthOptions) authenticate(r *http.Request) bool {
	const basicScheme string = "Basic "

	if r == nil {
		return false
	}

	// return true if user has already been set
	if user := context.Get(r, "user"); user != nil {
		return true
	}

	// Confirm the request is sending Basic Authentication credentials.
	auth := r.Header.Get("Authorization")
	if !strings.HasPrefix(auth, basicScheme) {
		return false
	}

	if user := z.authenticateDatabase(auth[len(basicScheme):]); user != nil {
		context.Set(r, "user", user)
		return true
	}

	return false

}

func (z *BasicAuthOptions) authenticateDatabase(token string) *model.User {

	var result *model.User

	username, password, err := decodeUsernamePassword(token)
	if err != nil {
		return nil
	}

	z.Database.Select(func(tx model.Transaction) error {
		if user, err := model.UserGetByUsername(username, tx); err == nil && user != nil {
			if user.MatchPassword(password) {
				result = user
			}
		}
		return nil
	})

	return result

}

func decodeUsernamePassword(token string) (string, string, error) {

	// Get the plain-text username and password from the request.
	// The first six characters are skipped - e.g. "Basic ".
	str, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		return "", "", nil
	}

	// Split on the first ":" character only, with any subsequent colons assumed to be part
	// of the password. Note that the RFC2617 standard does not place any limitations on
	// allowable characters in the password.
	creds := bytes.SplitN(str, []byte(":"), 2)
	if len(creds) != 2 {
		return "", "", nil
	}

	username := string(creds[0])
	password := string(creds[1])

	return username, password, nil

}
