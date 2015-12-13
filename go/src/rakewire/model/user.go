package model

//User defines a system user
type User struct {
	ID           string `kv:"Username:2,FeverHash:2"`
	Username     string `kv:"Username:1"`
	PasswordHash string
	FeverHash    string `kv:"FeverHash:1"`
}

// NewUser creates a new user with the specified username
func NewUser(username string) *User {
	return &User{
		ID:       getUUID(),
		Username: username,
	}
}
