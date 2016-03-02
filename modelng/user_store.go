package modelng

// U groups methods for accessing users.
var U = &userStore{}

type userStore struct{}

// New creates a new user with the specified username
func (z *userStore) New(username string) *User {
	return &User{
		Username: username,
	}
}
