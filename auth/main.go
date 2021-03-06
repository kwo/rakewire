package auth

import (
	"errors"
	"strings"

	"github.com/kwo/rakewire/model"
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
	ID    string
	Name  string
	Roles []string
}

// HasRole tests if the user has been assigned the given role
func (z *User) HasRole(role string) bool {
	result := false
	if len(z.Roles) > 0 {
		for _, value := range z.Roles {
			if value == role {
				return true
			}
		}
	}
	return result
}

// HasAllRoles tests if the user has been assigned all the given roles
func (z *User) HasAllRoles(roles ...string) bool {
	result := true
	for _, role := range roles {
		result = result && z.HasRole(role)
	}
	return result
}

// Authenticate will authenticate and authorize a user
func Authenticate(db model.Database, authHeader string, roles ...string) (user *User, err error) {

	if len(authHeader) == 0 {
		err = ErrUnauthenticated
	} else if strings.HasPrefix(authHeader, schemeBasic) {
		user, err = authenticateBasic(db, authHeader, roles...)
	} else if strings.HasPrefix(authHeader, schemeJWT) {
		user, err = authenticateJWT(authHeader, roles...)
	} else {
		err = ErrUnauthenticated // unknown authentication scheme
	}

	return

}
