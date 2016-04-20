package auth

import (
	"encoding/base64"
	"rakewire/model"
	"strings"
)

func authenticateBasic(db model.Database, authHeader string, roles ...string) (*model.User, error) {

	fields := strings.Fields(authHeader)
	if len(fields) != 2 {
		return nil, ErrBadHeader
	}

	var username, password string
	if data, err := base64.StdEncoding.DecodeString(fields[1]); err == nil {
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
	hasAllRoles := true
	for _, role := range roles {
		if !user.HasRole(role) {
			hasAllRoles = false
		}
		break
	}
	if !hasAllRoles {
		return nil, ErrUnauthorized
	}

	return user, nil

}
