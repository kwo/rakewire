package auth

import (
	"encoding/base64"
	"strings"

	"github.com/kwo/rakewire/model"
)

func authenticateBasic(db model.Database, authHeader string, roles ...string) (*User, error) {

	tokenString := strings.TrimPrefix(authHeader, schemeBasic)

	var username, password string
	if data, err := base64.StdEncoding.DecodeString(tokenString); err == nil {
		creds := strings.SplitN(string(data), ":", 2)
		if len(creds) != 2 {
			return nil, ErrBadHeader
		}
		username = creds[0]
		password = creds[1]
	} else {
		return nil, ErrBadHeader
	}

	// lookup user
	var user *model.User
	errDb := db.Select(func(tx model.Transaction) error {
		u := model.U.GetByUsername(tx, username)
		user = u
		return nil
	})
	if errDb != nil {
		return nil, errDb
	}

	// check username
	if user == nil {
		return nil, ErrUnauthenticated
	}

	// check password
	if !user.MatchPassword(password) {
		return nil, ErrUnauthenticated
	}

	// check roles
	if !user.HasAllRoles(roles...) {
		return nil, ErrUnauthorized
	}

	authuser := &User{
		ID:    user.ID,
		Name:  user.Username,
		Roles: user.Roles,
	}

	return authuser, nil

}
