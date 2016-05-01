package auth

import (
	"errors"
	"rakewire/model"
	"strings"
)

const (
	schemeBasic = "Basic "
	schemeJWT   = "Bearer "
)

// package level errors
var (
	ErrBadHeader       = errors.New("Cannot parse authorization header")
	ErrUnauthenticated = errors.New("Unauthenticated")
	ErrUnauthorized    = errors.New("Unauthorized")
)

// User contains the username and roles of a user
type User struct {
	Name  string
	Roles []string
}

// Authenticate will authenticate and authorize a user
func Authenticate(db model.Database, authHeader string, roles ...string) (*User, error) {

	if len(authHeader) == 0 {
		return nil, ErrUnauthenticated
	} else if strings.HasPrefix(authHeader, schemeBasic) {
		return authenticateBasic(db, authHeader, roles...)
	} else if strings.HasPrefix(authHeader, schemeJWT) {
		return authenticateJWT(authHeader, roles...)
	}

	return nil, ErrUnauthenticated // unknown authentication scheme

}
