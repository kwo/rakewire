package auth

import (
	"errors"
	"rakewire/logger"
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
	log                = logger.New("auth")
)

// User contains the username and roles of a user
type User struct {
	ID    string
	Name  string
	Roles []string
}

// Authenticate will authenticate and authorize a user
func Authenticate(db model.Database, authHeader string, roles ...string) (user *User, err error) {

	var scheme string

	if len(authHeader) == 0 {
		err = ErrUnauthenticated
	} else if strings.HasPrefix(authHeader, schemeBasic) {
		scheme = schemeBasic
		user, err = authenticateBasic(db, authHeader, roles...)
	} else if strings.HasPrefix(authHeader, schemeJWT) {
		scheme = schemeJWT
		user, err = authenticateJWT(authHeader, roles...)
	} else {
		err = ErrUnauthenticated // unknown authentication scheme
	}

	if err == nil {
		log.Infof("%s%s", scheme, user.Name)
	}

	return

}
