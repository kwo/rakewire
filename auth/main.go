package auth

import (
	"errors"
	"rakewire/model"
	"strings"
)

const (
	schemeBasic = "Basic "
)

// package level errors
var (
	ErrBadHeader       = errors.New("Cannot parse authorization header")
	ErrUnauthenticated = errors.New("Unauthenticated")
	ErrUnauthorized    = errors.New("Unauthorized")
)

// Authenticate will authenticate and authorize a user
func Authenticate(db model.Database, authHeader string, roles ...string) (*model.User, error) {

	if authHeader == "" {
		return nil, ErrUnauthenticated
	} else if strings.HasPrefix(authHeader, schemeBasic) {
		return authenticateBasic(db, authHeader, roles...)
	}

	return nil, ErrUnauthenticated // unknown authentication scheme

}
